package services

import (
	"context"
	"fmt"
	"time"

	"github.com/yourusername/payments-api/internal/database"
	"github.com/yourusername/payments-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentService struct {
	db *database.MongoDB
}

func NewPaymentService(db *database.MongoDB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req models.CreatePaymentRequest) (*models.Payment, error) {
	payment := models.Payment{
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

	collection := s.db.GetCollection("payments")
	_, err := collection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID string) (*models.Payment, error) {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return nil, fmt.Errorf("ID de pago inválido")
	}

	collection := s.db.GetCollection("payments")
	var payment models.Payment

	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&payment)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("pago no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (s *PaymentService) GetPaymentsByUser(ctx context.Context, userID string) ([]models.Payment, error) {
	collection := s.db.GetCollection("payments")

	cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *PaymentService) GetPaymentsByEntity(ctx context.Context, entityType, entityID string) ([]models.Payment, error) {
	collection := s.db.GetCollection("payments")

	filter := bson.M{
		"entity_type": entityType,
		"entity_id":   entityID,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID string, req models.UpdatePaymentStatusRequest) error {
	objID, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return fmt.Errorf("ID de pago inválido")
	}

	collection := s.db.GetCollection("payments")

	update := bson.M{
		"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		},
	}

	// Si el estado es "completed", agregar processed_at
	if req.Status == "completed" {
		now := time.Now()
		update["$set"].(bson.M)["processed_at"] = now
	}

	// Si hay transaction_id, agregarlo
	if req.TransactionID != "" {
		update["$set"].(bson.M)["transaction_id"] = req.TransactionID
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("pago no encontrado")
	}

	return nil
}

func (s *PaymentService) GetPaymentsByStatus(ctx context.Context, status string) ([]models.Payment, error) {
	collection := s.db.GetCollection("payments")

	cursor, err := collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []models.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}

	return payments, nil
}

// ProcessPayment simula el procesamiento de un pago
// En producción, aquí se integraría con Stripe, MercadoPago, etc.
func (s *PaymentService) ProcessPayment(ctx context.Context, paymentID string) error {
	// Simular procesamiento
	time.Sleep(100 * time.Millisecond)

	// Actualizar a completado
	req := models.UpdatePaymentStatusRequest{
		Status:        "completed",
		TransactionID: fmt.Sprintf("TXN-%s", paymentID),
	}

	return s.UpdatePaymentStatus(ctx, paymentID, req)
}
