package handler

import (
	"errors"
	"net/http"
	"regexp"
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
	customers service.CustomerService
	cards     service.CardService
}

var uuidPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

func NewHandler(
	users service.UserService,
	auditLogs service.AuditLogService,
	rulesets service.RulesetService,
	roles service.RoleService,
	entityServices ...any,
) *Handler {
	h := &Handler{users: users, auditLogs: auditLogs, rulesets: rulesets, roles: roles}
	for _, svc := range entityServices {
		if customerService, ok := svc.(service.CustomerService); ok {
			h.customers = customerService
		}
		if cardService, ok := svc.(service.CardService); ok {
			h.cards = cardService
		}
	}
	return h
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

func (h *Handler) CreateCustomer(c *gin.Context) {
	var request domain.CustomerDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.ID != "" && !uuidPattern.MatchString(request.ID) {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	customer := &models.Customer{
		ID: request.ID, Status: request.Status,
		HomeCountry: request.HomeCountry, HomeCity: request.HomeCity,
	}
	if err := h.customers.CreateCustomer(c.Request.Context(), customer); err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not create customer")
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func (h *Handler) ListCustomers(c *gin.Context) {
	customers, err := h.customers.ListCustomers(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list customers")
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *Handler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	customer, err := h.customers.GetCustomer(c.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "customer not found")
		return
	}
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not get customer")
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *Handler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	var request domain.CustomerDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	customer := &models.Customer{
		Status: request.Status, HomeCountry: request.HomeCountry, HomeCity: request.HomeCity,
	}
	if err := h.customers.UpdateCustomer(c.Request.Context(), id, customer); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "customer not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not update customer")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
		return
	}
	if err := h.customers.DeleteCustomer(c.Request.Context(), id); errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(c, http.StatusNotFound, "customer not found")
		return
	} else if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not delete customer")
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreateCard(c *gin.Context) {
	var request domain.CardDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.ID != "" && !uuidPattern.MatchString(request.ID) {
		respondWithError(c, http.StatusBadRequest, "invalid card id")
		return
	}
	if !uuidPattern.MatchString(request.CustomerID) {
		respondWithError(c, http.StatusBadRequest, "invalid customer id")
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

func (h *Handler) ListCards(c *gin.Context) {
	cards, err := h.cards.ListCards(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list cards")
		return
	}
	c.JSON(http.StatusOK, cards)
}

func (h *Handler) GetCard(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
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

func (h *Handler) UpdateCard(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
		respondWithError(c, http.StatusBadRequest, "invalid card id")
		return
	}
	var request domain.CardDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if !uuidPattern.MatchString(request.CustomerID) {
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

func (h *Handler) DeleteCard(c *gin.Context) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
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
