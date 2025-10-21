package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/payments-api/internal/config"
	"github.com/yourusername/payments-api/internal/database"
	"github.com/yourusername/payments-api/internal/handlers"
	"github.com/yourusername/payments-api/internal/middleware"
	"github.com/yourusername/payments-api/internal/services"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Conectar a MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("❌ Error conectando a MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Crear servicios
	paymentService := services.NewPaymentService(mongoDB)

	// Crear handlers
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Configurar Gin
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/healthz", paymentHandler.HealthCheck)

	// Rutas de pagos
	paymentRoutes := router.Group("/payments")
	{
		paymentRoutes.POST("", paymentHandler.CreatePayment)
		paymentRoutes.GET("/:id", paymentHandler.GetPayment)
		paymentRoutes.GET("/user/:user_id", paymentHandler.GetPaymentsByUser)
		paymentRoutes.GET("/entity", paymentHandler.GetPaymentsByEntity) // Query: ?entity_type=subscription&entity_id=123
		paymentRoutes.GET("/status", paymentHandler.GetPaymentsByStatus) // Query: ?status=pending
		paymentRoutes.PATCH("/:id/status", paymentHandler.UpdatePaymentStatus)
		paymentRoutes.POST("/:id/process", paymentHandler.ProcessPayment)
	}

	// Iniciar servidor
	log.Printf("🚀 Payments API corriendo en puerto %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
