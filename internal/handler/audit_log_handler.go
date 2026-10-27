package handler

import (
	"net/http"

	"vigilant/internal/domain"
	"vigilant/internal/models"
	"vigilant/internal/service"

	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct{ auditLogs service.AuditLogService }

func NewAuditLogHandler(auditLogs service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogs: auditLogs}
}

func (h *AuditLogHandler) CreateAuditLog(c *gin.Context) {
	var request domain.AuditLogDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.Details == nil {
		request.Details = map[string]any{}
	}
	auditLog := &models.AuditLog{
		UserID: request.UserID, Action: request.Action, EntityType: request.EntityType,
		EntityID: request.EntityID, Details: request.Details,
	}
	if err := h.auditLogs.CreateAuditLog(c.Request.Context(), auditLog); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create audit log")
		return
	}
	c.JSON(http.StatusCreated, auditLog)
}

func (h *AuditLogHandler) ListAuditLogs(c *gin.Context) {
	auditLogs, err := h.auditLogs.ListAuditLogs(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list audit logs")
		return
	}
	c.JSON(http.StatusOK, auditLogs)
}
