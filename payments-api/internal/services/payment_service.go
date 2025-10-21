package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/payments-api/internal/domain/dtos"
	"github.com/yourusername/payments-api/internal/domain/entities"
	"github.com/yourusername/payments-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PaymentServiceNew - Servicio con Dependency Injection
type PaymentServiceNew struct {
	paymentRepo repository.PaymentRepository
}

// NewPaymentServiceNew - Constructor con DI
func NewPaymentServiceNew(paymentRepo repository.PaymentRepository) *PaymentServiceNew {
	return &PaymentServiceNew{
		paymentRepo: paymentRepo,
	}
}

// CreatePayment crea un nuevo pago
func (s *PaymentServiceNew) CreatePayment(ctx context.Context, req dtos.CreatePaymentRequest) (dtos.PaymentResponse, error) {
	payment := entities.Payment{
		ID:             primitive.NewObjectID(),
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		UserID:         req.UserID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "pending",
		PaymentMethod:  req.PaymentMethod,
		PaymentGateway: req.PaymentGateway,
		Metadata:       req.Metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.paymentRepo.Create(ctx, &payment); err != nil {
		return dtos.PaymentResponse{}, fmt.Errorf("error al crear pago: %w", err)
	}

	return dtos.ToPaymentResponse(
		payment.ID,
		payment.EntityType,
		payment.EntityID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethod,
		payment.PaymentGateway,
		payment.TransactionID,
		payment.Metadata,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.ProcessedAt,
	), nil
}

// GetPaymentByID obtiene un pago por su ID
func (s *PaymentServiceNew) GetPaymentByID(ctx context.Context, paymentID string) (dtos.PaymentResponse, error) {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return dtos.PaymentResponse{}, fmt.Errorf("ID de pago inválido")
	}

	payment, err := s.paymentRepo.FindByID(ctx, objID)
	if err != nil {
		return dtos.PaymentResponse{}, err
	}

	return dtos.ToPaymentResponse(
		payment.ID,
		payment.EntityType,
		payment.EntityID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.PaymentMethod,
		payment.PaymentGateway,
		payment.TransactionID,
		payment.Metadata,
		payment.CreatedAt,
		payment.UpdatedAt,
		payment.ProcessedAt,
	), nil
}

// GetPaymentsByUser obtiene todos los pagos de un usuario
func (s *PaymentServiceNew) GetPaymentsByUser(ctx context.Context, userID string) ([]dtos.PaymentResponse, error) {
	payments, err := s.paymentRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = dtos.ToPaymentResponse(
			payment.ID,
			payment.EntityType,
			payment.EntityID,
			payment.UserID,
			payment.Amount,
			payment.Currency,
			payment.Status,
			payment.PaymentMethod,
			payment.PaymentGateway,
			payment.TransactionID,
			payment.Metadata,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ProcessedAt,
		)
	}

	return responses, nil
}

// GetPaymentsByEntity obtiene pagos asociados a una entidad
func (s *PaymentServiceNew) GetPaymentsByEntity(ctx context.Context, entityType, entityID string) ([]dtos.PaymentResponse, error) {
	payments, err := s.paymentRepo.FindByEntity(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = dtos.ToPaymentResponse(
			payment.ID,
			payment.EntityType,
			payment.EntityID,
			payment.UserID,
			payment.Amount,
			payment.Currency,
			payment.Status,
			payment.PaymentMethod,
			payment.PaymentGateway,
			payment.TransactionID,
			payment.Metadata,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ProcessedAt,
		)
	}

	return responses, nil
}

// GetPaymentsByStatus obtiene pagos por estado
func (s *PaymentServiceNew) GetPaymentsByStatus(ctx context.Context, status string) ([]dtos.PaymentResponse, error) {
	payments, err := s.paymentRepo.FindByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	responses := make([]dtos.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = dtos.ToPaymentResponse(
			payment.ID,
			payment.EntityType,
			payment.EntityID,
			payment.UserID,
			payment.Amount,
			payment.Currency,
			payment.Status,
			payment.PaymentMethod,
			payment.PaymentGateway,
			payment.TransactionID,
			payment.Metadata,
			payment.CreatedAt,
			payment.UpdatedAt,
			payment.ProcessedAt,
		)
	}

	return responses, nil
}

// UpdatePaymentStatus actualiza el estado de un pago
func (s *PaymentServiceNew) UpdatePaymentStatus(ctx context.Context, paymentID string, req dtos.UpdatePaymentStatusRequest) error {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return fmt.Errorf("ID de pago inválido")
	}

	return s.paymentRepo.UpdateStatus(ctx, objID, req.Status, req.TransactionID)
}

// ProcessPayment simula el procesamiento de un pago
func (s *PaymentServiceNew) ProcessPayment(ctx context.Context, paymentID string) error {
	// Simular procesamiento
	time.Sleep(100 * time.Millisecond)

	// Actualizar a completado
	req := dtos.UpdatePaymentStatusRequest{
		Status:        "completed",
		TransactionID: fmt.Sprintf("TXN-%s", paymentID),
	}

	return s.UpdatePaymentStatus(ctx, paymentID, req)
}
