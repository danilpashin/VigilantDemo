package handler

import (
	"github.com/gin-gonic/gin"
)

func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}
