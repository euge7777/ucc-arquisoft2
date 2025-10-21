package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/payments-api/internal/config"
	"github.com/yourusername/payments-api/internal/controllers"
	"github.com/yourusername/payments-api/internal/dao"
	"github.com/yourusername/payments-api/internal/database"
	"github.com/yourusername/payments-api/internal/middleware"
	"github.com/yourusername/payments-api/internal/services"
)

func main() {
	// 1. Cargar configuración
	cfg := config.LoadConfig()

	// 2. Conectar a MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("❌ Error conectando a MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// 3. Inicializar DAOs (Implementaciones de Repository) con DI
	paymentRepo := dao.NewPaymentRepositoryMongo(mongoDB.Database)

	// 4. Inicializar Services (Lógica de Negocio) con DI
	paymentService := services.NewPaymentServiceNew(paymentRepo)

	// 5. Inicializar Controllers (Capa HTTP) con DI
	paymentController := controllers.NewPaymentController(paymentService)

	// 6. Configurar Gin Router
	router := gin.Default()
	router.Use(middleware.CORS())

	// 7. Registrar Rutas
	registerRoutes(router, paymentController)

	// 8. Iniciar servidor
	log.Printf("🚀 Payments API corriendo en puerto %s", cfg.Port)
	log.Println("📦 Arquitectura: Controllers → Services → Repositories")
	log.Println("💉 Dependency Injection: Activada")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// registerRoutes - Registra todas las rutas HTTP
func registerRoutes(router *gin.Engine, paymentController *controllers.PaymentController) {
	// Health check
	router.GET("/healthz", paymentController.HealthCheck)

	// Rutas de pagos
	paymentRoutes := router.Group("/payments")
	{
		paymentRoutes.POST("", paymentController.CreatePayment)
		paymentRoutes.GET("/:id", paymentController.GetPayment)
		paymentRoutes.GET("/user/:user_id", paymentController.GetPaymentsByUser)
		paymentRoutes.GET("/entity", paymentController.GetPaymentsByEntity)   // Query: ?entity_type=subscription&entity_id=123
		paymentRoutes.GET("/status", paymentController.GetPaymentsByStatus)   // Query: ?status=pending
		paymentRoutes.PATCH("/:id/status", paymentController.UpdatePaymentStatus)
		paymentRoutes.POST("/:id/process", paymentController.ProcessPayment)
	}
}
