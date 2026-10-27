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

type CustomerHandler struct{ customers service.CustomerService }

func NewCustomerHandler(customers service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customers: customers}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var request domain.CustomerDTO
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if request.ID != "" {
		if _, err := uuid.Parse(request.ID); err != nil {
			respondWithError(c, http.StatusBadRequest, "invalid customer id")
			return
		}
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

func (h *CustomerHandler) ListCustomers(c *gin.Context) {
	customers, err := h.customers.ListCustomers(c.Request.Context())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "could not list customers")
		return
	}
	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
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

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
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

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
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
