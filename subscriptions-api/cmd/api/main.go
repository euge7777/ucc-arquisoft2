package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gym-management/subscriptions-api/internal/clients"
	"github.com/yourusername/gym-management/subscriptions-api/internal/config"
	"github.com/yourusername/gym-management/subscriptions-api/internal/controllers"
	"github.com/yourusername/gym-management/subscriptions-api/internal/dao"
	"github.com/yourusername/gym-management/subscriptions-api/internal/database"
	"github.com/yourusername/gym-management/subscriptions-api/internal/middleware"
	"github.com/yourusername/gym-management/subscriptions-api/internal/services"
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
	planRepo := dao.NewPlanRepositoryMongo(mongoDB.Database)
	subscriptionRepo := dao.NewSubscriptionRepositoryMongo(mongoDB.Database)

	// 4. Inicializar Clients (Servicios Externos) con DI
	usersValidator := clients.NewUsersAPIValidator(cfg.UsersAPIURL)

	// Intentar conectar a RabbitMQ, si falla usar NullEventPublisher
	var eventPublisher services.EventPublisher
	rabbitPublisher, err := clients.NewRabbitMQEventPublisher(cfg.RabbitMQURL, cfg.RabbitMQExchange)
	if err != nil {
		log.Printf("⚠️  Warning: No se pudo conectar a RabbitMQ: %v", err)
		log.Println("⚠️  Usando NullEventPublisher como fallback")
		eventPublisher = clients.NewNullEventPublisher()
	} else {
		eventPublisher = rabbitPublisher
		defer rabbitPublisher.Close()
	}

	// 5. Inicializar Services (Lógica de Negocio) con DI
	planService := services.NewPlanService(planRepo)
	subscriptionService := services.NewSubscriptionService(
		subscriptionRepo,
		planRepo,
		usersValidator,
		eventPublisher,
	)
	healthService := services.NewHealthService(mongoDB.Client, eventPublisher)

	// 6. Inicializar Payment Event Handler y Consumer
	paymentHandler := services.NewPaymentEventHandlerService(subscriptionService)

	// Intentar conectar al consumidor de eventos de pago
	paymentConsumer, err := clients.NewPaymentEventConsumer(
		cfg.RabbitMQURL,
		cfg.RabbitMQExchange,
		"subscriptions-api.payment-events", // nombre de la cola
	)
	if err != nil {
		log.Printf("⚠️  Warning: No se pudo conectar al consumidor de eventos de pago: %v", err)
		log.Println("⚠️  Las suscripciones NO se activarán automáticamente al completar pagos")
	} else {
		defer paymentConsumer.Close()

		// Iniciar el consumo de eventos en background
		ctx := context.Background()
		if err := paymentConsumer.StartConsuming(ctx, paymentHandler); err != nil {
			log.Printf("❌ Error iniciando consumidor de eventos: %v", err)
		} else {
			log.Println("✅ Consumidor de eventos de pago iniciado (escuchando payment.completed.subscription)")
		}
	}

	// 7. Inicializar Controllers (Capa HTTP) con DI
	planController := controllers.NewPlanController(planService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService, healthService)

	// 8. Configurar Gin Router
	router := gin.Default()
	router.Use(middleware.CORS())

	// 9. Registrar Rutas
	registerRoutes(router, planController, subscriptionController, cfg)

	// 10. Iniciar servidor
	log.Printf("🚀 Subscriptions API corriendo en puerto %s", cfg.Port)
	log.Println("📦 Arquitectura: Controllers → Services → Repositories")
	log.Println("💉 Dependency Injection: Activada")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// registerRoutes - Registra todas las rutas HTTP
func registerRoutes(
	router *gin.Engine,
	planController *controllers.PlanController,
	subscriptionController *controllers.SubscriptionController,
	cfg *config.Config,
) {
	// Health check (público)
	router.GET("/healthz", subscriptionController.HealthCheck)

	// Rutas públicas de planes (solo lectura)
	publicPlanRoutes := router.Group("/plans")
	{
		publicPlanRoutes.GET("", planController.ListPlans)
		publicPlanRoutes.GET("/:id", planController.GetPlan)
	}

	// Rutas protegidas de planes (solo admins)
	protectedPlanRoutes := router.Group("/plans")
	protectedPlanRoutes.Use(middleware.JWTAuth(cfg.JWTSecret))
	protectedPlanRoutes.Use(middleware.RequireRole("admin"))
	{
		protectedPlanRoutes.POST("", planController.CreatePlan)
	}

	// Rutas protegidas de suscripciones (requieren autenticación)
	subscriptionRoutes := router.Group("/subscriptions")
	subscriptionRoutes.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		subscriptionRoutes.POST("", subscriptionController.CreateSubscription)
		subscriptionRoutes.GET("/:id", subscriptionController.GetSubscription)
		subscriptionRoutes.GET("/active/:user_id", subscriptionController.GetActiveSubscriptionByUser)
		subscriptionRoutes.PATCH("/:id/status", subscriptionController.UpdateSubscriptionStatus)
		subscriptionRoutes.DELETE("/:id", subscriptionController.CancelSubscription)
	}
}
