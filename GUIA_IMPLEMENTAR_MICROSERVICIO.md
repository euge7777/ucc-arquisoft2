# 🚀 GUÍA COMPLETA: Cómo Implementar un Microservicio Desde Cero

Esta guía te muestra **PASO A PASO** cómo crear un microservicio nuevo usando `users-api` como template.

**Ejemplo:** Vamos a crear `subscriptions-api` (MongoDB)

---

## 📋 Tabla de Contenidos

1. [Pre-requisitos](#pre-requisitos)
2. [Paso 1: Crear Estructura de Carpetas](#paso-1-crear-estructura-de-carpetas)
3. [Paso 2: Domain Models (Modelos de Negocio)](#paso-2-domain-models)
4. [Paso 3: DAO Models (Modelos de Base de Datos)](#paso-3-dao-models)
5. [Paso 4: Repository (Acceso a Datos)](#paso-4-repository)
6. [Paso 5: Services (Lógica de Negocio)](#paso-5-services)
7. [Paso 6: Controllers (Handlers HTTP)](#paso-6-controllers)
8. [Paso 7: Middleware](#paso-7-middleware)
9. [Paso 8: Config](#paso-8-config)
10. [Paso 9: Main.go (Dependency Injection)](#paso-9-maingo)
11. [Paso 10: Archivos de Configuración](#paso-10-archivos-de-configuración)
12. [Paso 11: Probar el Microservicio](#paso-11-probar)
13. [Paso 12: Docker](#paso-12-docker)
14. [Checklist Final](#checklist-final)

---

## Pre-requisitos

- ✅ Go 1.22+ instalado
- ✅ Docker instalado (opcional)
- ✅ `users-api` funcionando (como referencia)
- ✅ Conocer el negocio del microservicio que vas a crear

---

## Paso 1: Crear Estructura de Carpetas

### 1.1. Crear carpetas base

```bash
cd C:\Users\eli_v\ucc-arquisoft2

# Crear microservicio nuevo
mkdir subscriptions-api
cd subscriptions-api

# Crear estructura completa
mkdir -p cmd/api
mkdir -p internal/config
mkdir -p internal/domain
mkdir -p internal/dao
mkdir -p internal/repository
mkdir -p internal/services
mkdir -p internal/controllers
mkdir -p internal/middleware
mkdir -p internal/clients  # Para RabbitMQ, HTTP calls, etc.
```

### 1.2. Verificar estructura

```
subscriptions-api/
├── cmd/
│   └── api/
│       └── main.go          (TODO: crear después)
├── internal/
│   ├── config/
│   ├── domain/
│   ├── dao/
│   ├── repository/
│   ├── services/
│   ├── controllers/
│   ├── middleware/
│   └── clients/
├── go.mod                   (TODO: crear después)
├── .env.example
├── Dockerfile
└── README.md
```

---

## Paso 2: Domain Models (Modelos de Negocio)

**Archivo:** `internal/domain/plan.go`

Los modelos de dominio son **independientes de la base de datos**.

### 2.1. Identificar las entidades

Para `subscriptions-api`:
- **Plan** - Definición de un plan (Básico, Premium)
- **Suscripción** - Instancia de un usuario suscrito a un plan

### 2.2. Crear modelo de dominio

```go
// internal/domain/plan.go
package domain

import "time"

// Plan representa un plan de suscripción (negocio)
type Plan struct {
	ID          string    `json:"id"`
	Nombre      string    `json:"nombre"`
	Descripcion string    `json:"descripcion"`
	Precio      float64   `json:"precio"`
	TipoAcceso  string    `json:"tipo_acceso"` // "basico" o "completo"
	Duracion    int       `json:"duracion"`    // Días
	Activo      bool      `json:"activo"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PlanCreate - DTO para crear plan
type PlanCreate struct {
	Nombre      string  `json:"nombre" binding:"required"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio" binding:"required,min=0"`
	TipoAcceso  string  `json:"tipo_acceso" binding:"required,oneof=basico completo"`
	Duracion    int     `json:"duracion" binding:"required,min=1"`
}

// PlanUpdate - DTO para actualizar plan
type PlanUpdate struct {
	Nombre      string  `json:"nombre" binding:"required"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio" binding:"required,min=0"`
	TipoAcceso  string  `json:"tipo_acceso" binding:"required,oneof=basico completo"`
	Duracion    int     `json:"duracion" binding:"required,min=1"`
	Activo      bool    `json:"activo"`
}

// PlanResponse - DTO para respuesta HTTP
type PlanResponse struct {
	ID          string  `json:"id"`
	Nombre      string  `json:"nombre"`
	Descripcion string  `json:"descripcion"`
	Precio      float64 `json:"precio"`
	TipoAcceso  string  `json:"tipo_acceso"`
	Duracion    int     `json:"duracion"`
	Activo      bool    `json:"activo"`
}

// ToResponse convierte de Plan a PlanResponse
func (p Plan) ToResponse() PlanResponse {
	return PlanResponse{
		ID:          p.ID,
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		Precio:      p.Precio,
		TipoAcceso:  p.TipoAcceso,
		Duracion:    p.Duracion,
		Activo:      p.Activo,
	}
}
```

### 2.3. Repetir para cada entidad

Crear también:
- `internal/domain/suscripcion.go`
- `internal/domain/...` (otras entidades)

**Reglas importantes:**
- ✅ NO usar tags de GORM/MongoDB aquí
- ✅ Usar tipos de Go nativos (string, int, float64, time.Time)
- ✅ Crear DTOs: Create, Update, Response
- ✅ Crear método `ToResponse()`

---

## Paso 3: DAO Models (Modelos de Base de Datos)

**Archivo:** `internal/dao/Plan.go`

Los modelos DAO son **específicos de la base de datos**.

### 3.1. Para MongoDB

```go
// internal/dao/Plan.go
package dao

import (
	"subscriptions-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Plan representa el modelo de MongoDB
type Plan struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Nombre      string             `bson:"nombre"`
	Descripcion string             `bson:"descripcion"`
	Precio      float64            `bson:"precio"`
	TipoAcceso  string             `bson:"tipo_acceso"` // "basico" o "completo"
	Duracion    int                `bson:"duracion"`
	Activo      bool               `bson:"activo"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// CollectionName especifica el nombre de la colección
func (Plan) CollectionName() string {
	return "planes"
}

// ToDomain convierte de DAO (MongoDB) a Domain (negocio)
func (p Plan) ToDomain() domain.Plan {
	return domain.Plan{
		ID:          p.ID.Hex(),
		Nombre:      p.Nombre,
		Descripcion: p.Descripcion,
		Precio:      p.Precio,
		TipoAcceso:  p.TipoAcceso,
		Duracion:    p.Duracion,
		Activo:      p.Activo,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// FromDomain convierte de Domain (negocio) a DAO (MongoDB)
func PlanFromDomain(domainPlan domain.Plan) Plan {
	var id primitive.ObjectID
	if domainPlan.ID != "" {
		id, _ = primitive.ObjectIDFromHex(domainPlan.ID)
	}

	return Plan{
		ID:          id,
		Nombre:      domainPlan.Nombre,
		Descripcion: domainPlan.Descripcion,
		Precio:      domainPlan.Precio,
		TipoAcceso:  domainPlan.TipoAcceso,
		Duracion:    domainPlan.Duracion,
		Activo:      domainPlan.Activo,
		CreatedAt:   domainPlan.CreatedAt,
		UpdatedAt:   domainPlan.UpdatedAt,
	}
}
```

### 3.2. Para MySQL (si fuera el caso)

```go
// internal/dao/Plan.go (MySQL)
package dao

import (
	"subscriptions-api/internal/domain"
	"time"
)

// Plan representa el modelo de MySQL con GORM
type Plan struct {
	ID          uint       `gorm:"column:id;primaryKey;autoIncrement"`
	Nombre      string     `gorm:"type:varchar(100);not null"`
	Descripcion string     `gorm:"type:text"`
	Precio      float64    `gorm:"type:decimal(10,2);not null"`
	TipoAcceso  string     `gorm:"type:enum('basico','completo');not null"`
	Duracion    int        `gorm:"type:int;not null"`
	Activo      bool       `gorm:"default:true;not null"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"` // Soft delete
}

// TableName especifica el nombre de la tabla
func (Plan) TableName() string {
	return "planes"
}

// ToDomain y FromDomain (similar a MongoDB)
```

**Reglas importantes:**
- ✅ Usar tags de la BD (bson, gorm)
- ✅ Siempre tener `ToDomain()` y `FromDomain()`
- ✅ ObjectID para MongoDB, uint para MySQL

---

## Paso 4: Repository (Acceso a Datos)

**Archivo:** `internal/repository/planes_mongo.go`

El repository es la **única capa que habla con la BD**.

### 4.1. Definir interfaz

```go
// internal/repository/planes_mongo.go
package repository

import (
	"context"
	"subscriptions-api/internal/domain"
)

// PlanesRepository define la interfaz del repositorio
type PlanesRepository interface {
	Create(ctx context.Context, plan domain.Plan) (domain.Plan, error)
	GetByID(ctx context.Context, id string) (domain.Plan, error)
	List(ctx context.Context) ([]domain.Plan, error)
	Update(ctx context.Context, id string, plan domain.Plan) (domain.Plan, error)
	Delete(ctx context.Context, id string) error
}
```

### 4.2. Implementar repositorio MongoDB

```go
// internal/repository/planes_mongo.go
package repository

import (
	"context"
	"errors"
	"fmt"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/dao"
	"subscriptions-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPlanesRepository implementa PlanesRepository usando MongoDB
type MongoPlanesRepository struct {
	collection *mongo.Collection
}

// NewMongoPlanesRepository crea una nueva instancia del repository
func NewMongoPlanesRepository(cfg config.MongoConfig) *MongoPlanesRepository {
	// Conectar a MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to MongoDB: %v", err))
	}

	// Ping
	if err := client.Ping(ctx, nil); err != nil {
		panic(fmt.Sprintf("Error pinging MongoDB: %v", err))
	}

	collection := client.Database(cfg.Database).Collection("planes")

	fmt.Println("✅ Connected to MongoDB successfully (Planes)")

	return &MongoPlanesRepository{
		collection: collection,
	}
}

// Create inserta un nuevo plan
func (r *MongoPlanesRepository) Create(ctx context.Context, plan domain.Plan) (domain.Plan, error) {
	planDAO := dao.PlanFromDomain(plan)
	planDAO.ID = primitive.NewObjectID()
	planDAO.CreatedAt = time.Now()
	planDAO.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, planDAO)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("error creating plan: %w", err)
	}

	return planDAO.ToDomain(), nil
}

// GetByID obtiene un plan por ID
func (r *MongoPlanesRepository) GetByID(ctx context.Context, id string) (domain.Plan, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Plan{}, errors.New("invalid ID format")
	}

	var planDAO dao.Plan
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&planDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Plan{}, errors.New("plan not found")
		}
		return domain.Plan{}, fmt.Errorf("error getting plan: %w", err)
	}

	return planDAO.ToDomain(), nil
}

// List obtiene todos los planes
func (r *MongoPlanesRepository) List(ctx context.Context) ([]domain.Plan, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error listing plans: %w", err)
	}
	defer cursor.Close(ctx)

	var planesDAO []dao.Plan
	if err := cursor.All(ctx, &planesDAO); err != nil {
		return nil, fmt.Errorf("error decoding plans: %w", err)
	}

	// Convertir a Domain
	planes := make([]domain.Plan, len(planesDAO))
	for i, planDAO := range planesDAO {
		planes[i] = planDAO.ToDomain()
	}

	return planes, nil
}

// Update actualiza un plan existente
func (r *MongoPlanesRepository) Update(ctx context.Context, id string, plan domain.Plan) (domain.Plan, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Plan{}, errors.New("invalid ID format")
	}

	planDAO := dao.PlanFromDomain(plan)
	planDAO.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"nombre":      planDAO.Nombre,
			"descripcion": planDAO.Descripcion,
			"precio":      planDAO.Precio,
			"tipo_acceso": planDAO.TipoAcceso,
			"duracion":    planDAO.Duracion,
			"activo":      planDAO.Activo,
			"updated_at":  planDAO.UpdatedAt,
		},
	}

	result := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var updatedPlanDAO dao.Plan
	if err := result.Decode(&updatedPlanDAO); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Plan{}, errors.New("plan not found")
		}
		return domain.Plan{}, fmt.Errorf("error updating plan: %w", err)
	}

	return updatedPlanDAO.ToDomain(), nil
}

// Delete elimina un plan
func (r *MongoPlanesRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("error deleting plan: %w", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("plan not found")
	}

	return nil
}
```

**Reglas importantes:**
- ✅ Siempre definir interfaz primero
- ✅ Convertir entre DAO ↔ Domain
- ✅ Manejar errores específicos (not found, invalid ID)
- ✅ Usar context para timeouts

---

## Paso 5: Services (Lógica de Negocio)

**Archivo:** `internal/services/planes.go`

El service contiene **validaciones y lógica de negocio**.

```go
// internal/services/planes.go
package services

import (
	"context"
	"errors"
	"fmt"
	"subscriptions-api/internal/domain"
	"subscriptions-api/internal/repository"
)

// PlanesService define la interfaz del servicio
type PlanesService interface {
	Create(ctx context.Context, planCreate domain.PlanCreate) (domain.PlanResponse, error)
	GetByID(ctx context.Context, id string) (domain.PlanResponse, error)
	List(ctx context.Context) ([]domain.PlanResponse, error)
	Update(ctx context.Context, id string, planUpdate domain.PlanUpdate) (domain.PlanResponse, error)
	Delete(ctx context.Context, id string) error
}

// PlanesServiceImpl implementa PlanesService
type PlanesServiceImpl struct {
	repository repository.PlanesRepository
}

// NewPlanesService crea una nueva instancia del servicio
func NewPlanesService(repo repository.PlanesRepository) *PlanesServiceImpl {
	return &PlanesServiceImpl{
		repository: repo,
	}
}

// Create crea un nuevo plan
func (s *PlanesServiceImpl) Create(ctx context.Context, planCreate domain.PlanCreate) (domain.PlanResponse, error) {
	// Validaciones de negocio
	if err := s.validatePlanCreate(planCreate); err != nil {
		return domain.PlanResponse{}, err
	}

	// Crear dominio
	plan := domain.Plan{
		Nombre:      planCreate.Nombre,
		Descripcion: planCreate.Descripcion,
		Precio:      planCreate.Precio,
		TipoAcceso:  planCreate.TipoAcceso,
		Duracion:    planCreate.Duracion,
		Activo:      true, // Por defecto activo
	}

	createdPlan, err := s.repository.Create(ctx, plan)
	if err != nil {
		return domain.PlanResponse{}, fmt.Errorf("error creating plan: %w", err)
	}

	return createdPlan.ToResponse(), nil
}

// GetByID obtiene un plan por ID
func (s *PlanesServiceImpl) GetByID(ctx context.Context, id string) (domain.PlanResponse, error) {
	plan, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.PlanResponse{}, fmt.Errorf("plan not found: %w", err)
	}

	return plan.ToResponse(), nil
}

// List obtiene todos los planes
func (s *PlanesServiceImpl) List(ctx context.Context) ([]domain.PlanResponse, error) {
	planes, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing plans: %w", err)
	}

	// Convertir a Response DTO
	responses := make([]domain.PlanResponse, len(planes))
	for i, plan := range planes {
		responses[i] = plan.ToResponse()
	}

	return responses, nil
}

// Update actualiza un plan existente
func (s *PlanesServiceImpl) Update(ctx context.Context, id string, planUpdate domain.PlanUpdate) (domain.PlanResponse, error) {
	// Validaciones de negocio
	if err := s.validatePlanUpdate(planUpdate); err != nil {
		return domain.PlanResponse{}, err
	}

	// Crear dominio
	plan := domain.Plan{
		Nombre:      planUpdate.Nombre,
		Descripcion: planUpdate.Descripcion,
		Precio:      planUpdate.Precio,
		TipoAcceso:  planUpdate.TipoAcceso,
		Duracion:    planUpdate.Duracion,
		Activo:      planUpdate.Activo,
	}

	updatedPlan, err := s.repository.Update(ctx, id, plan)
	if err != nil {
		return domain.PlanResponse{}, fmt.Errorf("error updating plan: %w", err)
	}

	return updatedPlan.ToResponse(), nil
}

// Delete elimina un plan
func (s *PlanesServiceImpl) Delete(ctx context.Context, id string) error {
	// TODO: Validar que no haya suscripciones activas con este plan
	// suscripcionesActivas, err := s.suscripcionesRepo.CountByPlanID(ctx, id)
	// if suscripcionesActivas > 0 {
	//     return errors.New("no se puede eliminar un plan con suscripciones activas")
	// }

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("error deleting plan: %w", err)
	}

	return nil
}

// validatePlanCreate valida los datos para crear un plan
func (s *PlanesServiceImpl) validatePlanCreate(planCreate domain.PlanCreate) error {
	if planCreate.Nombre == "" {
		return errors.New("el nombre del plan no puede estar vacío")
	}

	if planCreate.Precio < 0 {
		return errors.New("el precio no puede ser negativo")
	}

	if planCreate.Duracion <= 0 {
		return errors.New("la duración debe ser mayor a 0 días")
	}

	if planCreate.TipoAcceso != "basico" && planCreate.TipoAcceso != "completo" {
		return errors.New("tipo_acceso debe ser 'basico' o 'completo'")
	}

	return nil
}

// validatePlanUpdate valida los datos para actualizar un plan
func (s *PlanesServiceImpl) validatePlanUpdate(planUpdate domain.PlanUpdate) error {
	// Mismas validaciones que Create
	return s.validatePlanCreate(domain.PlanCreate{
		Nombre:      planUpdate.Nombre,
		Descripcion: planUpdate.Descripcion,
		Precio:      planUpdate.Precio,
		TipoAcceso:  planUpdate.TipoAcceso,
		Duracion:    planUpdate.Duracion,
	})
}
```

**Reglas importantes:**
- ✅ Siempre definir interfaz primero
- ✅ Validaciones de negocio aquí (NO en controller)
- ✅ Usar DTOs (Create, Update, Response)
- ✅ Retornar Response DTOs (nunca Domain directo)

---

## Paso 6: Controllers (Handlers HTTP)

**Archivo:** `internal/controllers/planes.go`

El controller maneja **peticiones y respuestas HTTP**.

```go
// internal/controllers/planes.go
package controllers

import (
	"net/http"
	"subscriptions-api/internal/domain"
	"subscriptions-api/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

// PlanesController maneja las peticiones HTTP de planes
type PlanesController struct {
	service services.PlanesService
}

// NewPlanesController crea una nueva instancia del controller
func NewPlanesController(service services.PlanesService) *PlanesController {
	return &PlanesController{
		service: service,
	}
}

// List obtiene todos los planes
// GET /planes
func (c *PlanesController) List(ctx *gin.Context) {
	planes, err := c.service.List(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener planes"})
		return
	}

	ctx.JSON(http.StatusOK, planes)
}

// GetByID obtiene un plan por ID
// GET /planes/:id
func (c *PlanesController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	plan, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener plan"})
		}
		return
	}

	ctx.JSON(http.StatusOK, plan)
}

// Create crea un nuevo plan
// POST /planes (admin only)
func (c *PlanesController) Create(ctx *gin.Context) {
	var planCreate domain.PlanCreate
	if err := ctx.ShouldBindJSON(&planCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto", "details": err.Error()})
		return
	}

	createdPlan, err := c.service.Create(ctx.Request.Context(), planCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear plan", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdPlan)
}

// Update actualiza un plan existente
// PUT /planes/:id (admin only)
func (c *PlanesController) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var planUpdate domain.PlanUpdate
	if err := ctx.ShouldBindJSON(&planUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Datos con formato incorrecto", "details": err.Error()})
		return
	}

	updatedPlan, err := c.service.Update(ctx.Request.Context(), id, planUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar plan", "details": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, updatedPlan)
}

// Delete elimina un plan
// DELETE /planes/:id (admin only)
func (c *PlanesController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Plan no encontrado"})
		} else if strings.Contains(err.Error(), "suscripciones activas") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar plan", "details": err.Error()})
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
```

**Reglas importantes:**
- ✅ Solo parsear HTTP y llamar al service
- ✅ NO poner lógica de negocio aquí
- ✅ Manejar errores específicos (404, 400, 500)
- ✅ Usar `ctx.Request.Context()` para context

---

## Paso 7: Middleware

**Copiar de `users-api`:**

```bash
# Copiar CORS
cp ../users-api/internal/middleware/cors.go internal/middleware/

# Copiar JWT
cp ../users-api/internal/middleware/jwt.go internal/middleware/
```

**NO modificar nada** - El middleware debe ser idéntico para que los JWT funcionen en todos los servicios.

---

## Paso 8: Config

**Archivo:** `internal/config/config.go`

```go
// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Mongo    MongoConfig
	JWT      JWTConfig
	RabbitMQ RabbitMQConfig // Si usas RabbitMQ
}

type MongoConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret string
}

type RabbitMQConfig struct {
	Host  string
	Port  string
	User  string
	Pass  string
	Queue string
}

func Load() Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	return Config{
		Port: getEnv("PORT", "8081"),
		Mongo: MongoConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DB", "gym_db"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "my-secret-key"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:  getEnv("RABBITMQ_HOST", "localhost"),
			Port:  getEnv("RABBITMQ_PORT", "5672"),
			User:  getEnv("RABBITMQ_USER", "admin"),
			Pass:  getEnv("RABBITMQ_PASS", "admin"),
			Queue: getEnv("RABBITMQ_QUEUE", "gym-events"),
		},
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
```

---

## Paso 9: Main.go (Dependency Injection)

**Archivo:** `cmd/api/main.go`

Este es el **corazón del microservicio** - aquí conectas todas las capas.

```go
// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/controllers"
	"subscriptions-api/internal/middleware"
	"subscriptions-api/internal/repository"
	"subscriptions-api/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Cargar configuración
	cfg := config.Load()

	// ========== CAPA DE DATOS (REPOSITORY) ==========
	// Crear repositorios con dependency injection
	planesRepo := repository.NewMongoPlanesRepository(cfg.Mongo)
	// suscripcionesRepo := repository.NewMongoSuscripcionesRepository(cfg.Mongo)

	// TODO: Si usas RabbitMQ
	// rabbitmqClient := clients.NewRabbitMQClient(cfg.RabbitMQ)

	// ========== CAPA DE NEGOCIO (SERVICES) ==========
	// Crear servicios con dependency injection
	planesService := services.NewPlanesService(planesRepo)
	// suscripcionesService := services.NewSuscripcionesService(suscripcionesRepo, planesRepo, rabbitmqClient)

	// ========== CAPA DE PRESENTACIÓN (CONTROLLERS) ==========
	// Crear controllers con dependency injection
	planesController := controllers.NewPlanesController(planesService)
	// suscripcionesController := controllers.NewSuscripcionesController(suscripcionesService)

	// ========== CONFIGURACIÓN DE GIN ==========
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/healthz", healthCheckHandler)

	// ========== RUTAS PÚBLICAS ==========
	// Planes (solo lectura sin auth)
	router.GET("/planes", planesController.List)
	router.GET("/planes/:id", planesController.GetByID)

	// ========== RUTAS PROTEGIDAS (REQUIEREN JWT) ==========
	protected := router.Group("/")
	protected.Use(middleware.JWTAuthMiddleware(cfg.JWT.Secret))
	{
		// Suscripciones (requieren autenticación)
		// protected.GET("/suscripciones", suscripcionesController.List)
		// protected.POST("/suscripciones", suscripcionesController.Create)
		// protected.GET("/suscripciones/active/:user_id", suscripcionesController.GetActiveByUser)
	}

	// ========== RUTAS DE ADMIN (REQUIEREN JWT + ADMIN) ==========
	adminOnly := protected.Group("/")
	adminOnly.Use(middleware.AdminOnlyMiddleware())
	{
		// Planes (CRUD completo solo admin)
		adminOnly.POST("/planes", planesController.Create)
		adminOnly.PUT("/planes/:id", planesController.Update)
		adminOnly.DELETE("/planes/:id", planesController.Delete)
	}

	// ========== INICIAR SERVIDOR ==========
	port := cfg.Port
	log.Printf("🚀 Subscriptions API running on port %s", port)
	log.Printf("📋 Endpoints disponibles:")
	log.Printf("   GET    /healthz")
	log.Printf("   GET    /planes")
	log.Printf("   GET    /planes/:id")
	log.Printf("   POST   /planes (admin)")
	log.Printf("   PUT    /planes/:id (admin)")
	log.Printf("   DELETE /planes/:id (admin)")

	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func healthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "subscriptions-api",
		"version": "1.0.0",
	})
}
```

**Orden del Dependency Injection:**

```
1. Config → Se carga primero
2. Repositories → Reciben Config
3. Services → Reciben Repositories
4. Controllers → Reciben Services
5. Router → Recibe Controllers
```

---

## Paso 10: Archivos de Configuración

### 10.1. go.mod

```bash
cd subscriptions-api
go mod init subscriptions-api

# Instalar dependencias
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v4
go get github.com/joho/godotenv
go get go.mongodb.org/mongo-driver/mongo
```

### 10.2. .env.example

```env
# Server Configuration
PORT=8081

# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=gym_db

# JWT Configuration (debe ser el mismo que users-api)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# RabbitMQ Configuration (opcional)
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=admin
RABBITMQ_PASS=admin
RABBITMQ_QUEUE=gym-events
```

### 10.3. Dockerfile

```dockerfile
# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o subscriptions-api ./cmd/api

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root/

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/subscriptions-api .

# Expose port
EXPOSE 8081

# Run
CMD ["./subscriptions-api"]
```

### 10.4. README.md

```markdown
# Subscriptions API

Microservicio para gestionar planes y suscripciones del gimnasio.

**Puerto:** 8081
**Base de datos:** MongoDB

## Endpoints

- `GET /planes` - Lista todos los planes
- `GET /planes/:id` - Obtiene un plan por ID
- `POST /planes` - Crea un plan (admin)
- `PUT /planes/:id` - Actualiza un plan (admin)
- `DELETE /planes/:id` - Elimina un plan (admin)

## Ejecutar

```bash
go run cmd/api/main.go
```
```

---

## Paso 11: Probar

### 11.1. Probar compilación

```bash
cd subscriptions-api
go mod download
go run cmd/api/main.go
```

### 11.2. Probar endpoints

```bash
# Health check
curl http://localhost:8081/healthz

# Listar planes (vacío al principio)
curl http://localhost:8081/planes

# Crear plan (necesitas JWT de admin)
curl -X POST http://localhost:8081/planes \
  -H "Authorization: Bearer <token_admin>" \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Básico",
    "descripcion": "Acceso a actividades básicas",
    "precio": 50.00,
    "tipo_acceso": "basico",
    "duracion": 30
  }'
```

---

## Paso 12: Docker

### 12.1. Agregar a docker-compose.new.yml

```yaml
  subscriptions-api:
    build:
      context: ./subscriptions-api
      dockerfile: Dockerfile
    container_name: gym-subscriptions-api
    environment:
      PORT: 8081
      MONGO_URI: mongodb://mongo:27017
      MONGO_DB: gym_db
      JWT_SECRET: my-super-secret-jwt-key
      RABBITMQ_HOST: rabbitmq
      RABBITMQ_PORT: 5672
      RABBITMQ_USER: admin
      RABBITMQ_PASS: admin
      RABBITMQ_QUEUE: gym-events
    ports:
      - "8081:8081"
    depends_on:
      mongo:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
    networks:
      - gym-network
    restart: unless-stopped
```

### 12.2. Probar con Docker

```bash
cd ..
docker-compose -f docker-compose.new.yml up --build subscriptions-api
```

---

## Checklist Final

### Estructura

- [ ] Carpetas creadas: domain, dao, repository, services, controllers, middleware, config
- [ ] `cmd/api/main.go` creado
- [ ] `go.mod` creado e inicializado

### Código

- [ ] Domain models creados (con DTOs: Create, Update, Response)
- [ ] DAO models creados (con tags de BD y conversiones ToDomain/FromDomain)
- [ ] Repository creado (con interfaz y implementación)
- [ ] Service creado (con interfaz, implementación y validaciones)
- [ ] Controller creado (con handlers HTTP)
- [ ] Middleware copiado (cors.go, jwt.go)
- [ ] Config creado
- [ ] Main.go con DI completo

### Archivos

- [ ] `.env.example` creado
- [ ] `Dockerfile` creado
- [ ] `README.md` creado

### Testing

- [ ] `go mod download` funciona
- [ ] `go run cmd/api/main.go` compila y ejecuta
- [ ] Health check funciona (`curl /healthz`)
- [ ] Endpoints públicos funcionan
- [ ] Endpoints protegidos requieren JWT
- [ ] Endpoints admin requieren JWT + admin

### Docker

- [ ] Agregado a `docker-compose.new.yml`
- [ ] `docker-compose up --build <servicio>` funciona
- [ ] Se conecta a las BDs correctamente
- [ ] Logs muestran "Connected successfully"

---

## 🎯 Resumen de Capas

```
HTTP Request
     ↓
[Controller] ← Parsea HTTP, llama Service
     ↓
[Service] ← Validaciones de negocio, llama Repository
     ↓
[Repository] ← Acceso a BD, convierte DAO ↔ Domain
     ↓
[DAO] ← Modelo de BD (MongoDB/MySQL)
     ↓
Base de Datos
```

**Flujo de datos:**

```
HTTP JSON → DTO → Domain → DAO → BD
           ↑              ↓
      Controller     Repository
                ↓
           Service (validaciones)
```

---

## 📚 Referencia Rápida

### MongoDB vs MySQL

| Aspecto | MongoDB | MySQL |
|---------|---------|-------|
| **ID** | `primitive.ObjectID` | `uint` |
| **Tags** | `bson:"nombre"` | `gorm:"column:nombre"` |
| **Importar** | `go.mongodb.org/mongo-driver/mongo` | `gorm.io/gorm` |
| **Colección/Tabla** | `CollectionName()` | `TableName()` |
| **Queries** | `bson.M{"_id": id}` | `.Where("id = ?", id)` |

### Dependencias comunes

```bash
# Gin (HTTP framework)
go get github.com/gin-gonic/gin

# JWT
go get github.com/golang-jwt/jwt/v4

# .env
go get github.com/joho/godotenv

# MongoDB
go get go.mongodb.org/mongo-driver/mongo

# MySQL/GORM
go get gorm.io/gorm
go get gorm.io/driver/mysql

# RabbitMQ
go get github.com/streadway/amqp
```

---

🎉 **¡Listo! Con esta guía tus compañeros pueden crear cualquier microservicio nuevo.**

**Archivos de referencia:**
- `users-api/` - Template completo con MySQL
- `activities-api/` - Template completo con MySQL
- Esta guía - Para MongoDB

**Siguiente paso:** Crear `subscriptions-api` siguiendo esta guía paso a paso.
