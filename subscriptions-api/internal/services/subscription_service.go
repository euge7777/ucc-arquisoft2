package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/gym-management/subscriptions-api/internal/domain/dtos"
	"github.com/yourusername/gym-management/subscriptions-api/internal/domain/entities"
	"github.com/yourusername/gym-management/subscriptions-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubscriptionService - Servicio de lógica de negocio para suscripciones
type SubscriptionService struct {
	subscriptionRepo repository.SubscriptionRepository // DI
	planRepo         repository.PlanRepository         // DI
	userService      UserValidator                     // DI (Interface para validar usuarios)
	eventPublisher   EventPublisher                    // DI (Interface para publicar eventos)
}

// UserValidator - Interface para validar usuarios (abstrae users-api)
type UserValidator interface {
	ValidateUser(ctx context.Context, userID string) (bool, error)
}

// EventPublisher - Interface para publicar eventos (abstrae RabbitMQ)
type EventPublisher interface {
	PublishSubscriptionEvent(action, subscriptionID string, data map[string]interface{}) error
}

// NewSubscriptionService - Constructor con DI
func NewSubscriptionService(
	subscriptionRepo repository.SubscriptionRepository,
	planRepo repository.PlanRepository,
	userService UserValidator,
	eventPublisher EventPublisher,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		userService:      userService,
		eventPublisher:   eventPublisher,
	}
}

// CreateSubscription - Crea una nueva suscripción
func (s *SubscriptionService) CreateSubscription(ctx context.Context, req dtos.CreateSubscriptionRequest) (*dtos.SubscriptionResponse, error) {
	// 1. Validar usuario existe
	valid, err := s.userService.ValidateUser(ctx, req.UsuarioID)
	if err != nil || !valid {
		return nil, fmt.Errorf("usuario no válido: %w", err)
	}

	// 2. Obtener plan
	planObjID, err := primitive.ObjectIDFromHex(req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("ID de plan inválido")
	}

	plan, err := s.planRepo.FindByID(ctx, planObjID)
	if err != nil {
		return nil, fmt.Errorf("plan no encontrado: %w", err)
	}

	if !plan.Activo {
		return nil, fmt.Errorf("el plan no está activo")
	}

	// 3. Calcular fechas
	now := time.Now()
	fechaVencimiento := now.AddDate(0, 0, plan.DuracionDias)

	// 4. Crear suscripción
	subscription := &entities.Subscription{
		ID:               primitive.NewObjectID(),
		UsuarioID:        req.UsuarioID,
		PlanID:           planObjID,
		SucursalOrigenID: req.SucursalOrigenID,
		FechaInicio:      now,
		FechaVencimiento: fechaVencimiento,
		Estado:           "pendiente_pago",
		Metadata: entities.Metadata{
			MetodoPagoPreferido: req.MetodoPago,
			AutoRenovacion:      req.AutoRenovacion,
			Notas:               req.Notas,
		},
		HistorialRenovaciones: []entities.Renovacion{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// 5. Guardar en repositorio
	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	// 6. Publicar evento
	eventData := map[string]interface{}{
		"usuario_id": subscription.UsuarioID,
		"plan_id":    subscription.PlanID.Hex(),
		"estado":     subscription.Estado,
	}
	s.eventPublisher.PublishSubscriptionEvent("create", subscription.ID.Hex(), eventData)

	// 7. Mapear a DTO de respuesta
	return s.mapSubscriptionToResponse(subscription, plan.Nombre), nil
}

// GetSubscriptionByID - Obtiene una suscripción por ID
func (s *SubscriptionService) GetSubscriptionByID(ctx context.Context, id string) (*dtos.SubscriptionResponse, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido")
	}

	subscription, err := s.subscriptionRepo.FindByID(ctx, objID)
	if err != nil {
		return nil, err
	}

	// Enriquecer con nombre del plan
	plan, _ := s.planRepo.FindByID(ctx, subscription.PlanID)
	planNombre := ""
	if plan != nil {
		planNombre = plan.Nombre
	}

	return s.mapSubscriptionToResponse(subscription, planNombre), nil
}

// GetActiveSubscriptionByUserID - Obtiene la suscripción activa de un usuario
func (s *SubscriptionService) GetActiveSubscriptionByUserID(ctx context.Context, userID string) (*dtos.SubscriptionResponse, error) {
	subscription, err := s.subscriptionRepo.FindActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Enriquecer con nombre del plan
	plan, _ := s.planRepo.FindByID(ctx, subscription.PlanID)
	planNombre := ""
	if plan != nil {
		planNombre = plan.Nombre
	}

	return s.mapSubscriptionToResponse(subscription, planNombre), nil
}

// UpdateSubscriptionStatus - Actualiza el estado de una suscripción
func (s *SubscriptionService) UpdateSubscriptionStatus(ctx context.Context, id string, req dtos.UpdateSubscriptionStatusRequest) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido")
	}

	if err := s.subscriptionRepo.UpdateStatus(ctx, objID, req.Estado, req.PagoID); err != nil {
		return err
	}

	// Publicar evento
	eventData := map[string]interface{}{
		"estado":  req.Estado,
		"pago_id": req.PagoID,
	}
	s.eventPublisher.PublishSubscriptionEvent("update", id, eventData)

	return nil
}

// CancelSubscription - Cancela una suscripción
func (s *SubscriptionService) CancelSubscription(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("ID inválido")
	}

	if err := s.subscriptionRepo.UpdateStatus(ctx, objID, "cancelada", ""); err != nil {
		return err
	}

	// Publicar evento
	s.eventPublisher.PublishSubscriptionEvent("delete", id, nil)

	return nil
}

// mapSubscriptionToResponse - Helper para mapear entidad a DTO
func (s *SubscriptionService) mapSubscriptionToResponse(subscription *entities.Subscription, planNombre string) *dtos.SubscriptionResponse {
	var renovaciones []dtos.RenovacionResponse
	for _, r := range subscription.HistorialRenovaciones {
		renovaciones = append(renovaciones, dtos.RenovacionResponse{
			Fecha:  r.Fecha,
			PagoID: r.PagoID,
			Monto:  r.Monto,
		})
	}

	return &dtos.SubscriptionResponse{
		ID:                    subscription.ID.Hex(),
		UsuarioID:             subscription.UsuarioID,
		PlanID:                subscription.PlanID.Hex(),
		PlanNombre:            planNombre,
		SucursalOrigenID:      subscription.SucursalOrigenID,
		FechaInicio:           subscription.FechaInicio,
		FechaVencimiento:      subscription.FechaVencimiento,
		Estado:                subscription.Estado,
		PagoID:                subscription.PagoID,
		AutoRenovacion:        subscription.Metadata.AutoRenovacion,
		MetodoPagoPreferido:   subscription.Metadata.MetodoPagoPreferido,
		Notas:                 subscription.Metadata.Notas,
		HistorialRenovaciones: renovaciones,
		CreatedAt:             subscription.CreatedAt,
		UpdatedAt:             subscription.UpdatedAt,
	}
}
