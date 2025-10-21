# Subscriptions API - Arquitectura Limpia con DI

Microservicio de gestión de suscripciones y planes con **arquitectura limpia completa** y **dependency injection**.

## 🏗️ Arquitectura

```
cmd/api/
└── main.go                    ← Inyección de Dependencias Manual

internal/
├── domain/                    ← Capa de Dominio
│   ├── entities/             # Entidades de negocio (Plan, Subscription)
│   └── dtos/                 # Data Transfer Objects (Request/Response)
│
├── repository/                ← Capa de Datos
│   ├── plan_repository.go            # Interface (Inversión de Dependencias)
│   ├── plan_repository_mongo.go      # Implementación MongoDB
│   ├── subscription_repository.go    # Interface
│   └── subscription_repository_mongo.go  # Implementación MongoDB
│
├── services/                  ← Capa de Lógica de Negocio
│   ├── plan_service.go        # Depende de Repository (interface)
│   └── subscription_service.go  # Depende de Repository + UserValidator + EventPublisher
│
├── infrastructure/            ← Capa de Infraestructura
│   ├── users_api_validator.go     # Implementa UserValidator (HTTP)
│   └── rabbitmq_event_publisher.go  # Implementa EventPublisher (RabbitMQ)
│
├── controllers/               ← Capa HTTP
│   ├── plan_controller.go     # Depende de PlanService
│   └── subscription_controller.go  # Depende de SubscriptionService
│
├── middleware/                # CORS, Auth, etc.
├── config/                    # Configuración
└── database/                  # Conexión MongoDB
```

## 🔑 Conceptos Clave

### 1. DTOs (Data Transfer Objects)

Los DTOs **separan** las entidades de dominio de las requests/responses HTTP:

- **Entities** (`internal/domain/entities/`): Modelos de dominio que se mapean a MongoDB
- **DTOs** (`internal/domain/dtos/`): Modelos para API (requests, responses, queries)

```go
// Entity (dominio)
type Plan struct {
    ID primitive.ObjectID
    Nombre string
    ...
}

// DTO Request
type CreatePlanRequest struct {
    Nombre string `json:"nombre" binding:"required"`
    ...
}

// DTO Response
type PlanResponse struct {
    ID string `json:"id"`  // ObjectID → string
    Nombre string `json:"nombre"`
    ...
}
```

### 2. Repository Pattern con Interfaces

Los repositories **abstraen** el acceso a datos usando interfaces:

```go
// Interface (en repository/)
type PlanRepository interface {
    Create(ctx context.Context, plan *entities.Plan) error
    FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Plan, error)
    ...
}

// Implementación MongoDB (en repository/)
type PlanRepositoryMongo struct {
    collection *mongo.Collection
}

func NewPlanRepositoryMongo(db *mongo.Database) PlanRepository {
    return &PlanRepositoryMongo{...}
}
```

**Ventajas**:
- ✅ **Testeable**: Se puede mockear el repository en tests
- ✅ **Intercambiable**: Cambiar de MongoDB a PostgreSQL solo requiere crear `PlanRepositoryPostgres`
- ✅ **SOLID**: Inversión de dependencias (Dependency Inversion Principle)

### 3. Services con Dependency Injection

Los services reciben sus dependencias como **interfaces** (no implementaciones):

```go
type SubscriptionService struct {
    subscriptionRepo repository.SubscriptionRepository  // Interface
    planRepo         repository.PlanRepository          // Interface
    userService      UserValidator                      // Interface
    eventPublisher   EventPublisher                     // Interface
}

// Constructor con DI
func NewSubscriptionService(
    subscriptionRepo repository.SubscriptionRepository,
    planRepo repository.PlanRepository,
    userService UserValidator,
    eventPublisher EventPublisher,
) *SubscriptionService {
    return &SubscriptionService{...}
}
```

**Ventajas**:
- ✅ **Desacoplamiento**: Services no conocen implementaciones concretas
- ✅ **Testeable**: Se pueden inyectar mocks en tests
- ✅ **Flexible**: Cambiar implementaciones sin modificar services

### 4. Infrastructure (Servicios Externos)

La carpeta `infrastructure/` contiene **implementaciones** de interfaces definidas en `services/`:

```go
// Interface definida en services/
type UserValidator interface {
    ValidateUser(ctx context.Context, userID string) (bool, error)
}

// Implementación en infrastructure/
type UsersAPIValidator struct {
    baseURL string
    client  *http.Client
}

func (u *UsersAPIValidator) ValidateUser(ctx context.Context, userID string) (bool, error) {
    // Llama a users-api via HTTP
}
```

**Por qué NO usar "clients/"**:
- ❌ "clients" sugiere que es solo una biblioteca HTTP
- ✅ "infrastructure" indica que es parte de la infraestructura (puede ser HTTP, gRPC, file system, etc.)
- ✅ Más alineado con **Domain-Driven Design** y **Arquitectura Hexagonal**

### 5. Controllers vs Handlers

**Controllers** (usado en este proyecto):
- Capa HTTP que depende de Services
- Separa claramente routing de lógica de negocio

```go
type PlanController struct {
    planService *services.PlanService  // DI
}

func (c *PlanController) CreatePlan(ctx *gin.Context) {
    var req dtos.CreatePlanRequest
    ctx.ShouldBindJSON(&req)

    plan, err := c.planService.CreatePlan(ctx.Request.Context(), req)
    ctx.JSON(http.StatusCreated, plan)
}
```

### 6. Inyección de Dependencias Manual

En `cmd/api/main.go` se hace la **composición manual** de todas las dependencias:

```go
func main() {
    // 1. Inicializar DB
    mongoDB, _ := database.NewMongoDB(...)

    // 2. Inicializar Repositories
    planRepo := repository.NewPlanRepositoryMongo(mongoDB.Database)
    subscriptionRepo := repository.NewSubscriptionRepositoryMongo(mongoDB.Database)

    // 3. Inicializar Infrastructure
    usersValidator := infrastructure.NewUsersAPIValidator(usersAPIURL)
    eventPublisher, _ := infrastructure.NewRabbitMQEventPublisher(...)

    // 4. Inicializar Services (con DI)
    planService := services.NewPlanService(planRepo)
    subscriptionService := services.NewSubscriptionService(
        subscriptionRepo,
        planRepo,
        usersValidator,
        eventPublisher,
    )

    // 5. Inicializar Controllers (con DI)
    planController := controllers.NewPlanController(planService)
    subscriptionController := controllers.NewSubscriptionController(subscriptionService)

    // 6. Registrar rutas
    router := gin.Default()
    registerRoutes(router, planController, subscriptionController)
}
```

## 📦 Endpoints

```bash
# Planes
POST   /plans              - Crear plan
GET    /plans              - Listar planes (query: ?activo=true)
GET    /plans/:id          - Obtener plan por ID

# Suscripciones
POST   /subscriptions                  - Crear suscripción
GET    /subscriptions/:id              - Obtener suscripción
GET    /subscriptions/active/:user_id  - Suscripción activa del usuario
PATCH  /subscriptions/:id/status       - Actualizar estado
DELETE /subscriptions/:id              - Cancelar suscripción

# Health
GET    /healthz            - Health check
```

## 🧪 Testing

Para testear este microservicio, crear mocks de las interfaces:

```go
// mocks/plan_repository_mock.go
type PlanRepositoryMock struct {
    mock.Mock
}

func (m *PlanRepositoryMock) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.Plan, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*entities.Plan), args.Error(1)
}

// services/plan_service_test.go
func TestCreatePlan(t *testing.T) {
    // Arrange
    mockRepo := new(PlanRepositoryMock)
    service := NewPlanService(mockRepo)

    // Act & Assert
    ...
}
```

## 🚀 Ejecución

```bash
# Local
go mod tidy
go run cmd/api/main.go

# Docker
docker build -t subscriptions-api .
docker run -p 8081:8081 --env-file .env subscriptions-api
```

## 📝 Ejemplo de Uso

```bash
# 1. Crear plan
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Premium",
    "precio_mensual": 100.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true
  }'

# 2. Crear suscripción
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": "5",
    "plan_id": "507f1f77bcf86cd799439011",
    "metodo_pago": "credit_card"
  }'
```

## 🎯 Próximos Pasos

Para equipos que implementen otros microservicios, usar esta estructura como referencia:

1. ✅ Separar **Entities** de **DTOs**
2. ✅ Crear **Repositories** con interfaces
3. ✅ Implementar **Services** con DI
4. ✅ Crear **Infrastructure** para servicios externos
5. ✅ Crear **Controllers** (no handlers)
6. ✅ Hacer DI manual en **main.go**

## 📚 Patrones Implementados

- ✅ **Repository Pattern** - Abstracción de datos
- ✅ **Dependency Injection** - Desacoplamiento
- ✅ **DTO Pattern** - Separación de capas
- ✅ **Service Layer** - Lógica de negocio
- ✅ **Dependency Inversion** - SOLID Principles
- ✅ **Clean Architecture** - Capas bien definidas
