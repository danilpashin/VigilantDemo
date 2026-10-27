package handler

import (
	"errors"
	"net/http"
	"uuid"

	"vigilant/internal/domain"
	"vigilant/internal/models"
	"vigilant/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CardHandler struct{ cards service.CardService }

func NewCardHandler(cards service.CardService) *CardHandler {
	return &CardHandler{cards: cards}
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	var request domain.CardDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.ID != "" {
		if _, err := uuid.Parse(request.ID); err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid card id")
			return
		}
	}
	if request.CustomerID != "" {
		if _, err := uuid.Parse(request.CustomerID); err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid customer id")
			return
		}
	} else {
		respondWithError(c, http.StatusBadRequest, "no customer id")
		return
	}

	card := &models.Card{
		ID: request.ID, CustomerID: request.CustomerID, Token: request.Token,
		Last4: request.Last4, Status: request.Status, ExpiresAt: request.ExpiresAt,
	}
	if err := h.cards.CreateCard(c.Request.Context(), card); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create card")
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *CardHandler) ListCards(c *gin.Context) {
	cards, err := h.cards.ListCards(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list cards")
		return
	}
	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) GetCard(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid card id")
		return
	}
	card, err := h.cards.GetCard(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "card not found")
		return
	}
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not get card")
		return
	}
	c.JSON(http.StatusOK, card)
}

func (h *CardHandler) UpdateCard(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid card id")
		return
	}
	var request domain.CardDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if request.CustomerID == "" {
		respondWithError(c, http.StatusBadRequest, "no customer id")
		return
	}
	if _, err := uuid.Parse(request.CustomerID); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	card := &models.Card{
		CustomerID: request.CustomerID, Token: request.Token, Last4: request.Last4,
		Status: request.Status, ExpiresAt: request.ExpiresAt,
	}
	if err := h.cards.UpdateCard(c.Request.Context(), id, card); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "card not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not update card")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid card id")
		return
	}
	if err := h.cards.DeleteCard(c.Request.Context(), id); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "card not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not delete card")
		return
	}
	c.Status(http.StatusNoContent)
}
