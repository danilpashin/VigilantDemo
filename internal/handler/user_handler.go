package handler

import (
	"errors"
	"net/http"
	"strconv"

	"vigilant/internal/domain"
	"vigilant/internal/models"
	"vigilant/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct{ users service.UserService }

func NewUserHandler(users service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request domain.UserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(request.Password) > 72 {
		respondWithError(c, http.StatusBadRequest, "password must be at most 72 bytes")
		return
	}

	user := &models.User{
		Name:   request.Name,
		Login:  request.Login,
		RoleID: request.RoleID,
	}
	if err := h.users.CreateUser(c.Request.Context(), user, request.Password); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create user")
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.users.ListUsers(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list users")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondWithError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.users.GetUser(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not get user")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondWithError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	var request domain.UserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(request.Password) > 72 {
		respondWithError(c, http.StatusBadRequest, "password must be at most 72 bytes")
		return
	}

	user := &models.User{
		Name:   request.Name,
		Login:  request.Login,
		RoleID: request.RoleID,
	}
	if err := h.users.UpdateUser(c.Request.Context(), id, user, request.Password); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not update user")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		respondWithError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.users.DeleteUser(c.Request.Context(), id); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not delete user")
		return
	}
	c.Status(http.StatusNoContent)
}
