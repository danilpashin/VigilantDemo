package handler

import (
	"net/http"

	"vigilant/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct{ roles service.RoleService }

func NewRoleHandler(roles service.RoleService) *RoleHandler {
	return &RoleHandler{roles: roles}
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	roles, err := h.roles.ListRoles(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list roles")
		return
	}
	c.JSON(http.StatusOK, roles)
}
