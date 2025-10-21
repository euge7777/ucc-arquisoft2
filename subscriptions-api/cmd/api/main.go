package main

import (
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

	eventPublisher, err := clients.NewRabbitMQEventPublisher(cfg.RabbitMQURL, cfg.RabbitMQExchange)
	if err != nil {
		log.Printf("⚠️  Warning: No se pudo conectar a RabbitMQ: %v", err)
		// En desarrollo, continuar sin RabbitMQ
		// En producción, esto sería un error fatal
	}
	if eventPublisher != nil {
		defer eventPublisher.Close()
	}

	// 5. Inicializar Services (Lógica de Negocio) con DI
	planService := services.NewPlanService(planRepo)
	subscriptionService := services.NewSubscriptionService(
		subscriptionRepo,
		planRepo,
		usersValidator,
		eventPublisher,
	)

	// 6. Inicializar Controllers (Capa HTTP) con DI
	planController := controllers.NewPlanController(planService)
	subscriptionController := controllers.NewSubscriptionController(subscriptionService)

	// 7. Configurar Gin Router
	router := gin.Default()
	router.Use(middleware.CORS())

	// 8. Registrar Rutas
	registerRoutes(router, planController, subscriptionController)

	// 9. Iniciar servidor
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
) {
	// Health check
	router.GET("/healthz", subscriptionController.HealthCheck)

	// Rutas de planes
	planRoutes := router.Group("/plans")
	{
		planRoutes.POST("", planController.CreatePlan)
		planRoutes.GET("", planController.ListPlans)
		planRoutes.GET("/:id", planController.GetPlan)
	}

	// Rutas de suscripciones
	subscriptionRoutes := router.Group("/subscriptions")
	{
		subscriptionRoutes.POST("", subscriptionController.CreateSubscription)
		subscriptionRoutes.GET("/:id", subscriptionController.GetSubscription)
		subscriptionRoutes.GET("/active/:user_id", subscriptionController.GetActiveSubscriptionByUser)
		subscriptionRoutes.PATCH("/:id/status", subscriptionController.UpdateSubscriptionStatus)
		subscriptionRoutes.DELETE("/:id", subscriptionController.CancelSubscription)
	}
}
