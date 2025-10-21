package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/payments-api/internal/models"
	"github.com/yourusername/payments-api/internal/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req models.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.service.CreatePayment(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.service.GetPaymentByID(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPaymentsByUser(c *gin.Context) {
	userID := c.Param("user_id")

	payments, err := h.service.GetPaymentsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) GetPaymentsByEntity(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID := c.Query("entity_id")

	if entityType == "" || entityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type y entity_id son requeridos"})
		return
	}

	payments, err := h.service.GetPaymentsByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")

	var req models.UpdatePaymentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdatePaymentStatus(c.Request.Context(), paymentID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Estado actualizado correctamente"})
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	paymentID := c.Param("id")

	err := h.service.ProcessPayment(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pago procesado correctamente"})
}

func (h *PaymentHandler) GetPaymentsByStatus(c *gin.Context) {
	status := c.Query("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status es requerido"})
		return
	}

	payments, err := h.service.GetPaymentsByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *PaymentHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "payments-api",
	})
}
