package main

import (
	"activities-api/internal/config"
	"activities-api/internal/controllers"
	"activities-api/internal/middleware"
	"activities-api/internal/repository"
	"activities-api/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Cargar configuración
	cfg := config.Load()

	// ========== CAPA DE DATOS (REPOSITORY) ==========
	// Crear repositorio de actividades (comparte conexión DB)
	actividadesRepo := repository.NewMySQLActividadesRepository(cfg.MySQL)
	if actividadesRepo == nil {
		log.Fatal("Failed to initialize actividades repository")
	}

	// Crear repositorio de inscripciones (comparte la misma DB)
	inscripcionesRepo := repository.NewMySQLInscripcionesRepository(actividadesRepo.GetDB())

	// TODO: Cuando el equipo implemente Sucursales:
	// sucursalesRepo := repository.NewMySQLSucursalesRepository(actividadesRepo.GetDB())

	// ========== CAPA DE NEGOCIO (SERVICES) ==========
	// Crear servicios con dependency injection
	actividadesService := services.NewActividadesService(actividadesRepo)
	inscripcionesService := services.NewInscripcionesService(inscripcionesRepo, actividadesRepo)
	// TODO: sucursalesService := services.NewSucursalesService(sucursalesRepo)

	// TODO: Cuando el equipo implemente RabbitMQ:
	// rabbitmqClient := clients.NewRabbitMQClient(cfg.RabbitMQ)
	// inscripcionesService := services.NewInscripcionesService(inscripcionesRepo, actividadesRepo, rabbitmqClient)

	// ========== CAPA DE PRESENTACIÓN (CONTROLLERS) ==========
	// Crear controllers con dependency injection
	actividadesController := controllers.NewActividadesController(actividadesService)
	inscripcionesController := controllers.NewInscripcionesController(inscripcionesService)
	// TODO: sucursalesController := controllers.NewSucursalesController(sucursalesService)

	// ========== CONFIGURACIÓN DE GIN ==========
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/healthz", healthCheckHandler)

	// ========== RUTAS PÚBLICAS ==========
	// Actividades (solo lectura sin auth)
	router.GET("/actividades", actividadesController.List)
	router.GET("/actividades/buscar", actividadesController.Search)
	router.GET("/actividades/:id", actividadesController.GetByID)

	// TODO: Sucursales (solo lectura sin auth)
	// router.GET("/sucursales", sucursalesController.List)
	// router.GET("/sucursales/:id", sucursalesController.GetByID)

	// ========== RUTAS PROTEGIDAS (REQUIEREN JWT) ==========
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		// Inscripciones (requieren autenticación)
		protected.GET("/inscripciones", inscripcionesController.List)
		protected.POST("/inscripciones", inscripcionesController.Create)
		protected.DELETE("/inscripciones", inscripcionesController.Deactivate)
	}

	// ========== RUTAS DE ADMIN (REQUIEREN JWT + ADMIN) ==========
	adminOnly := protected.Group("/")
	adminOnly.Use(middleware.AdminOnlyMiddleware())
	{
		// Actividades (CRUD completo solo admin)
		adminOnly.POST("/actividades", actividadesController.Create)
		adminOnly.PUT("/actividades/:id", actividadesController.Update)
		adminOnly.DELETE("/actividades/:id", actividadesController.Delete)

		// TODO: Sucursales (CRUD completo solo admin)
		// adminOnly.POST("/sucursales", sucursalesController.Create)
		// adminOnly.PUT("/sucursales/:id", sucursalesController.Update)
		// adminOnly.DELETE("/sucursales/:id", sucursalesController.Delete)
	}

	// ========== INICIAR SERVIDOR ==========
	port := cfg.Port
	log.Printf("🚀 Activities API running on port %s", port)
	log.Printf("📋 Endpoints disponibles:")
	log.Printf("   GET    /healthz")
	log.Printf("   GET    /actividades")
	log.Printf("   GET    /actividades/buscar?id=&titulo=&horario=&categoria=")
	log.Printf("   GET    /actividades/:id")
	log.Printf("   POST   /actividades (admin)")
	log.Printf("   PUT    /actividades/:id (admin)")
	log.Printf("   DELETE /actividades/:id (admin)")
	log.Printf("   GET    /inscripciones (auth)")
	log.Printf("   POST   /inscripciones (auth)")
	log.Printf("   DELETE /inscripciones (auth)")

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "activities-api",
		"version": "1.0.0",
	})
}
