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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := &models.User{
		Name:         request.Name,
		Login:        request.Login,
		PasswordHash: request.PasswordHash,
		RoleID:       request.RoleID,
	}
	if err := h.users.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.users.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.users.GetUser(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var request domain.UserDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := &models.User{
		Name:         request.Name,
		Login:        request.Login,
		PasswordHash: request.PasswordHash,
		RoleID:       request.RoleID,
	}
	if err := h.users.UpdateUser(c.Request.Context(), id, user); errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.users.DeleteUser(c.Request.Context(), id); errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateAuditLog(c *gin.Context) {
	var request domain.AuditLogDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create audit log"})
		return
	}
	c.JSON(http.StatusCreated, auditLog)
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	auditLogs, err := h.auditLogs.ListAuditLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list audit logs"})
		return
	}
	c.JSON(http.StatusOK, auditLogs)
}

func (h *Handler) CreateRuleset(c *gin.Context) {
	var request domain.RulesetDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create ruleset"})
		return
	}
	c.JSON(http.StatusCreated, ruleset)
}

func (h *Handler) ListRulesets(c *gin.Context) {
	rulesets, err := h.rulesets.ListRulesets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list rulesets"})
		return
	}
	c.JSON(http.StatusOK, rulesets)
}

func (h *Handler) ListRoles(c *gin.Context) {
	roles, err := h.roles.ListRoles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}
