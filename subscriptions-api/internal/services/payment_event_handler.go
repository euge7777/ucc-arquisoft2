package services

import (
	"context"
	"fmt"
	"log"

	"github.com/yourusername/gym-management/subscriptions-api/internal/clients"
)

// PaymentEventHandlerService - Implementación del handler de eventos de pago
type PaymentEventHandlerService struct {
	subscriptionService *SubscriptionService
}

// NewPaymentEventHandlerService - Constructor
func NewPaymentEventHandlerService(subscriptionService *SubscriptionService) *PaymentEventHandlerService {
	return &PaymentEventHandlerService{
		subscriptionService: subscriptionService,
	}
}

// HandlePaymentCompleted - Maneja eventos de pago completado
func (h *PaymentEventHandlerService) HandlePaymentCompleted(ctx context.Context, event clients.PaymentEvent) error {
	log.Printf("🔄 Procesando pago completado: PaymentID=%s, SubscriptionID=%s, UserID=%s, Amount=%.2f %s\n",
		event.PaymentID, event.EntityID, event.UserID, event.Amount, event.Currency)

	// Activar la suscripción
	err := h.subscriptionService.ActivateSubscriptionFromPayment(ctx, event.EntityID, event.PaymentID)
	if err != nil {
		return fmt.Errorf("error activando suscripción %s: %w", event.EntityID, err)
	}

	log.Printf("✅ Suscripción %s activada exitosamente para el usuario %s\n", event.EntityID, event.UserID)
	return nil
}
