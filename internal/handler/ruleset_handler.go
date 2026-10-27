package handler

import (
	"net/http"

	"vigilant/internal/domain"
	"vigilant/internal/models"
	"vigilant/internal/service"

	"github.com/gin-gonic/gin"
)

type RulesetHandler struct{ rulesets service.RulesetService }

func NewRulesetHandler(rulesets service.RulesetService) *RulesetHandler {
	return &RulesetHandler{rulesets: rulesets}
}

func (h *RulesetHandler) CreateRuleset(c *gin.Context) {
	var request domain.RulesetDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ruleset := &models.Ruleset{
		Version: request.Version, ContentChecksum: request.ContentChecksum,
		ConfigSnapshot: request.ConfigSnapshot, ActivatedAt: request.ActivatedAt,
		DeactivatedAt: request.DeactivatedAt,
	}
	if err := h.rulesets.CreateRuleset(c.Request.Context(), ruleset); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create ruleset")
		return
	}
	c.JSON(http.StatusCreated, ruleset)
}

func (h *RulesetHandler) ListRulesets(c *gin.Context) {
	rulesets, err := h.rulesets.ListRulesets(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list rulesets")
		return
	}
	c.JSON(http.StatusOK, rulesets)
}
