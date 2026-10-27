package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"vigilant/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	if err := configureGinMode(); err != nil {
		return err
	}

	handled, err := runMigrationCommand()
	if err != nil {
		return err
	}
	if handled {
		return nil
	}

	databaseURL, err := getDatabaseURL()
	if err != nil {
		return err
	}
	return runApp(databaseURL)
}

func getDatabaseURL() (string, error) {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		return "", fmt.Errorf("DATABASE_URL is not set")
	}
	return databaseURL, nil
}

func configureGinMode() error {
	mode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	if mode == "" {
		mode = gin.DebugMode
	}

	switch mode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		gin.SetMode(mode)
		return nil
	default:
		return fmt.Errorf("invalid GIN_MODE %q; expected debug, release, or test", mode)
	}
}

func runApp(databaseURL string) error {
	db, err := database.Open(databaseURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Close(db); err != nil {
			log.Printf("close database: %v", err)
		}
	}()

	addr := strings.TrimSpace(os.Getenv("HTTP_ADDR"))
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           newRouter(db),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return runServer(server)
}

func newRouter(db *gorm.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/health", healthCheck(db))
	return router
}

func healthCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func runServer(server *http.Server) error {
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("application listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return gracefulShutdown(server)
	}
}

func gracefulShutdown(server *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shut down HTTP server: %w", err)
	}
	log.Println("HTTP server shut down gracefully")
	return nil
}
