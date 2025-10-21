package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/gym-management/search-api/internal/config"
	"github.com/yourusername/gym-management/search-api/internal/consumers"
	"github.com/yourusername/gym-management/search-api/internal/handlers"
	"github.com/yourusername/gym-management/search-api/internal/middleware"
	"github.com/yourusername/gym-management/search-api/internal/services"
)

func main() {
	// Cargar configuración
	cfg := config.LoadConfig()

	// Crear servicios
	searchService := services.NewSearchService()
	cacheService := services.NewCacheService(
		cfg.MemcachedServers,
		cfg.CacheTTL,
		cfg.LocalCacheTTL,
	)

	// Iniciar limpieza periódica del caché local
	cacheService.StartCleanupRoutine(5 * time.Minute)

	// Conectar a RabbitMQ como consumidor
	rabbitConsumer, err := consumers.NewRabbitMQConsumer(
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

	// Crear handlers
	searchHandler := handlers.NewSearchHandler(searchService, cacheService)

	// Configurar Gin
	router := gin.Default()
	router.Use(middleware.CORS())

	// Health check
	router.GET("/healthz", searchHandler.HealthCheck)

	// Rutas de búsqueda
	searchRoutes := router.Group("/search")
	{
		searchRoutes.POST("", searchHandler.Search)           // Búsqueda avanzada
		searchRoutes.GET("", searchHandler.QuickSearch)       // Búsqueda rápida (query params)
		searchRoutes.GET("/stats", searchHandler.GetStats)    // Estadísticas del índice
		searchRoutes.GET("/:id", searchHandler.GetDocument)   // Obtener documento por ID
		searchRoutes.POST("/index", searchHandler.IndexDocument) // Indexar documento manualmente
		searchRoutes.DELETE("/:id", searchHandler.DeleteDocument) // Eliminar documento
	}

	// Iniciar servidor
	log.Printf("🚀 Search API corriendo en puerto %s", cfg.Port)
	log.Printf("🔍 Sistema de búsqueda listo (in-memory mode)")
	log.Printf("💾 Caché de dos niveles activado (Local: %ds, Memcached: %ds)", cfg.LocalCacheTTL, cfg.CacheTTL)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
