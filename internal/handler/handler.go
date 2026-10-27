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

type Handler struct {
	users     service.UserService
	auditLogs service.AuditLogService
	rulesets  service.RulesetService
	roles     service.RoleService
}

func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

func NewHandler(
	users service.UserService,
	auditLogs service.AuditLogService,
	rulesets service.RulesetService,
	roles service.RoleService,
) *Handler {
	return &Handler{users: users, auditLogs: auditLogs, rulesets: rulesets, roles: roles}
}

func (h *Handler) CreateUser(c *gin.Context) {
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

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.users.ListUsers(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list users")
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
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

func (h *Handler) UpdateUser(c *gin.Context) {
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

func (h *Handler) DeleteUser(c *gin.Context) {
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

func (h *Handler) CreateAuditLog(c *gin.Context) {
	var request domain.AuditLogDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.Details == nil {
		request.Details = map[string]any{}
	}

	auditLog := &models.AuditLog{
		UserID:     request.UserID,
		Action:     request.Action,
		EntityType: request.EntityType,
		EntityID:   request.EntityID,
		Details:    request.Details,
	}
	if err := h.auditLogs.CreateAuditLog(c.Request.Context(), auditLog); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create audit log")
		return
	}
	c.JSON(http.StatusCreated, auditLog)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	auditLogs, err := h.auditLogs.ListAuditLogs(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list audit logs")
		return
	}
	c.JSON(http.StatusOK, auditLogs)
}

func (h *Handler) CreateRuleset(c *gin.Context) {
	var request domain.RulesetDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	ruleset := &models.Ruleset{
		Version:         request.Version,
		ContentChecksum: request.ContentChecksum,
		ConfigSnapshot:  request.ConfigSnapshot,
		ActivatedAt:     request.ActivatedAt,
		DeactivatedAt:   request.DeactivatedAt,
	}
	if err := h.rulesets.CreateRuleset(c.Request.Context(), ruleset); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create ruleset")
		return
	}
	c.JSON(http.StatusCreated, ruleset)
}

func (h *Handler) ListRulesets(c *gin.Context) {
	rulesets, err := h.rulesets.ListRulesets(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list rulesets")
		return
	}
	c.JSON(http.StatusOK, rulesets)
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.roles.ListRoles(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list roles")
		return
	}
	c.JSON(http.StatusOK, roles)
}
