package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gym-management/search-api/internal/clients"
	"github.com/yourusername/gym-management/search-api/internal/config"
	"github.com/yourusername/gym-management/search-api/internal/controllers"
	"github.com/yourusername/gym-management/search-api/internal/middleware"
	"github.com/yourusername/gym-management/search-api/internal/services"
)

func main() {
	// 1. Cargar configuración
	cfg := config.LoadConfig()

	// 2. Crear servicios con DI
	searchService := services.NewSearchService()
	cacheService := services.NewCacheService(
		cfg.MemcachedServers,
		cfg.CacheTTL,
		cfg.LocalCacheTTL,
	)

	// Iniciar limpieza periódica del caché local
	cacheService.StartCleanupRoutine(5 * time.Minute)

	// 3. Conectar a RabbitMQ como consumidor (Client externo)
	rabbitConsumer, err := clients.NewRabbitMQConsumer(
		cfg.RabbitMQURL,
		cfg.RabbitMQExchange,
		cfg.RabbitMQQueue,
		searchService,
		cacheService,
	)
	if err != nil {
		log.Printf("⚠️  Warning: No se pudo conectar a RabbitMQ: %v", err)
		// Continuar sin RabbitMQ para desarrollo
	} else {
		defer rabbitConsumer.Close()

		// Iniciar consumidor
		if err := rabbitConsumer.Start(); err != nil {
			log.Printf("⚠️  Warning: Error iniciando consumidor: %v", err)
		}
	}

	// 4. Crear controllers con DI
	searchController := controllers.NewSearchController(searchService, cacheService)

	// 5. Configurar Gin Router
	router := gin.Default()
	router.Use(middleware.CORS())

	// 6. Registrar Rutas
	registerRoutes(router, searchController)

	// 7. Iniciar servidor
	log.Printf("🚀 Search API corriendo en puerto %s", cfg.Port)
	log.Println("📦 Arquitectura: Controllers → Services")
	log.Println("💉 Dependency Injection: Activada")
	log.Printf("🔍 Sistema de búsqueda listo (in-memory mode)")
	log.Printf("💾 Caché de dos niveles activado (Local: %ds, Memcached: %ds)", cfg.LocalCacheTTL, cfg.CacheTTL)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}

// registerRoutes - Registra todas las rutas HTTP
func registerRoutes(router *gin.Engine, searchController *controllers.SearchController) {
	// Health check
	router.GET("/healthz", searchController.HealthCheck)

	// Rutas de búsqueda
	searchRoutes := router.Group("/search")
	{
		searchRoutes.POST("", searchController.Search)              // Búsqueda avanzada
		searchRoutes.GET("", searchController.QuickSearch)          // Búsqueda rápida (query params)
		searchRoutes.GET("/stats", searchController.GetStats)       // Estadísticas del índice
		searchRoutes.GET("/:id", searchController.GetDocument)      // Obtener documento por ID
		searchRoutes.POST("/index", searchController.IndexDocument) // Indexar documento manualmente
		searchRoutes.DELETE("/:id", searchController.DeleteDocument) // Eliminar documento
	}
}
