package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Payment representa un pago genérico
// Este modelo es completamente agnóstico del dominio
type Payment struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	EntityType      string                 `bson:"entity_type" json:"entity_type" binding:"required"` // subscription, inscription, plan_upgrade, etc.
	EntityID        string                 `bson:"entity_id" json:"entity_id" binding:"required"`
	UserID          string                 `bson:"user_id" json:"user_id" binding:"required"`
	Amount          float64                `bson:"amount" json:"amount" binding:"required,gt=0"`
	Currency        string                 `bson:"currency" json:"currency" binding:"required"` // USD, ARS, EUR
	Status          string                 `bson:"status" json:"status"`                        // pending, completed, failed, refunded
	PaymentMethod   string                 `bson:"payment_method" json:"payment_method"`        // credit_card, debit_card, cash, transfer
	PaymentGateway  string                 `bson:"payment_gateway,omitempty" json:"payment_gateway,omitempty"` // stripe, mercadopago, manual
	TransactionID   string                 `bson:"transaction_id,omitempty" json:"transaction_id,omitempty"`
	Metadata        map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"` // Información adicional específica del dominio
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
	ProcessedAt     *time.Time             `bson:"processed_at,omitempty" json:"processed_at,omitempty"`
}

// CreatePaymentRequest representa la petición para crear un pago
type CreatePaymentRequest struct {
	EntityType     string                 `json:"entity_type" binding:"required"`
	EntityID       string                 `json:"entity_id" binding:"required"`
	UserID         string                 `json:"user_id" binding:"required"`
	Amount         float64                `json:"amount" binding:"required,gt=0"`
	Currency       string                 `json:"currency" binding:"required"`
	PaymentMethod  string                 `json:"payment_method" binding:"required"`
	PaymentGateway string                 `json:"payment_gateway,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// UpdatePaymentStatusRequest representa la petición para actualizar el estado de un pago
type UpdatePaymentStatusRequest struct {
	Status        string `json:"status" binding:"required,oneof=pending completed failed refunded"`
	TransactionID string `json:"transaction_id,omitempty"`
}
