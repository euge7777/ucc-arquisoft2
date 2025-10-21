# 📖 LEE ESTO PRIMERO - Proyecto de Microservicios

## 🚨 IMPORTANTE: Estado del Proyecto

Este proyecto está **PARCIALMENTE IMPLEMENTADO** (40% completo).

**Solo 2 de 5 microservicios están completos:**
- ✅ **users-api** - Funcional (sin arquitectura limpia completa)
- ✅ **subscriptions-api** - **EJEMPLO DE REFERENCIA** (arquitectura limpia completa)

**3 microservicios por completar:**
- ⚠️ **activities-api** - Funcional pero falta RabbitMQ
- ❌ **payments-api** - Solo estructura básica
- ❌ **search-api** - Solo estructura básica

---

## 🎯 ¿Qué debo hacer?

### Si eres un desarrollador nuevo en el proyecto:

1. **Lee primero**:
   - ✅ Este archivo (`LEEME_PRIMERO.md`)
   - ✅ `ESTADO_IMPLEMENTACION.md` - Estado detallado de cada microservicio
   - ✅ `subscriptions-api/README.md` - Arquitectura limpia explicada

2. **Estudia el ejemplo de referencia**:
   - ✅ Explora `subscriptions-api/` completo
   - ✅ Entiende la estructura de carpetas
   - ✅ Revisa cómo se hace la inyección de dependencias en `main.go`

3. **Implementa tu microservicio asignado**:
   - ❌ **payments-api** - Equipo de pagos
   - ❌ **search-api** - Equipo de búsqueda

---

## 📚 Documentación del Proyecto

| Documento | Descripción | Cuándo Leerlo |
|-----------|-------------|---------------|
| **`LEEME_PRIMERO.md`** (este archivo) | Punto de entrada | **Primero** |
| **`ESTADO_IMPLEMENTACION.md`** | Estado de cada microservicio | **Segundo** |
| **`ARCHIVOS_A_REFACTORIZAR.md`** | Archivos viejos a eliminar | **⚠️ IMPORTANTE** |
| **`subscriptions-api/README.md`** | Arquitectura limpia explicada | **Tercero - MUY IMPORTANTE** |
| `ARQUITECTURA_MICROSERVICIOS.md` | Patrones generales | Referencia |
| `DIAGRAMA_ENTIDADES.md` | Modelo de datos | Referencia |
| `GUIA_COMPLETA_MICROSERVICIOS.md` | Guía de uso | Cuando esté todo funcionando |
| `TEST_COMMANDS.md` | Comandos de testing | Para probar |
| `docker-compose.new.yml` | Infraestructura completa | Cuando ejecutes todo |

---

## ⚠️ ALERTA: Archivos Viejos Presentes

**IMPORTANTE**: Varios microservicios todavía tienen archivos con estructura VIEJA:

| Microservicio | Archivos Viejos | Qué Hacer |
|---------------|-----------------|-----------|
| users-api | ⚠️ `handlers/`, `models/` | **NO refactorizar** - Ya funciona |
| activities-api | ⚠️ `handlers/`, `models/` | **Opcional** - Solo agregar RabbitMQ |
| payments-api | ⚠️ `handlers/`, `models/`, `services/` | **ELIMINAR y crear desde cero** |
| search-api | ⚠️ `handlers/`, `models/`, `clients/` | **ELIMINAR y refactorizar** |

**Ver `ARCHIVOS_A_REFACTORIZAR.md` para instrucciones detalladas.**

---

## 🏗️ Arquitectura Requerida (OBLIGATORIA)

Todos los microservicios **DEBEN** seguir esta estructura:

```
microservicio-api/
├── cmd/api/
│   └── main.go                    ← 💉 Inyección de Dependencias MANUAL
│
├── internal/
│   ├── domain/                    ← Capa de Dominio
│   │   ├── entities/             # Entidades de negocio
│   │   └── dtos/                 # Data Transfer Objects
│   │
│   ├── repository/                ← Capa de Datos
│   │   ├── *_repository.go            # Interfaces
│   │   └── *_repository_mongo.go      # Implementaciones
│   │
│   ├── services/                  ← Lógica de Negocio
│   │   └── *_service.go          # Depende de Repositories (interfaces)
│   │
│   ├── infrastructure/            ← Servicios Externos
│   │   ├── *_validator.go        # HTTP clients
│   │   └── *_publisher.go        # RabbitMQ, etc.
│   │
│   ├── controllers/               ← Capa HTTP
│   │   └── *_controller.go       # HTTP handlers
│   │
│   ├── middleware/                # CORS, Auth, etc.
│   ├── database/                  # Conexión a BD
│   └── config/                    # Configuración
│
├── Dockerfile
├── .env.example
├── go.mod
└── README.md                      ← Documentar tu microservicio
```

---

## ✅ Conceptos OBLIGATORIOS

### 1. DTOs (Data Transfer Objects)

**❌ NO hacer esto:**
```go
// models/user.go
type User struct {
    ID primitive.ObjectID `bson:"_id" json:"id"`  // ❌ Mezcla BD y API
    Name string `bson:"name" json:"name"`
}
```

**✅ SÍ hacer esto:**
```go
// domain/entities/user.go (para BD)
type User struct {
    ID primitive.ObjectID `bson:"_id"`
    Name string `bson:"name"`
}

// domain/dtos/user_dtos.go (para API)
type UserResponse struct {
    ID string `json:"id"`  // ObjectID convertido a string
    Name string `json:"name"`
}

type CreateUserRequest struct {
    Name string `json:"name" binding:"required"`
}
```

### 2. Repository Pattern

**❌ NO hacer esto:**
```go
// services/user_service.go
func (s *UserService) CreateUser(user User) error {
    collection := s.db.Collection("users")  // ❌ Service accede directamente a BD
    collection.InsertOne(user)
}
```

**✅ SÍ hacer esto:**
```go
// repository/user_repository.go (INTERFACE)
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
}

// repository/user_repository_mongo.go (IMPLEMENTACIÓN)
type UserRepositoryMongo struct {
    collection *mongo.Collection
}

// services/user_service.go
type UserService struct {
    userRepo repository.UserRepository  // ✅ Depende de INTERFACE
}
```

### 3. Dependency Injection Manual

**❌ NO hacer esto:**
```go
// main.go
func main() {
    db := database.Connect()
    service := services.NewUserService()  // ❌ Service crea sus propias dependencias
    handler := handlers.NewUserHandler()
}
```

**✅ SÍ hacer esto:**
```go
// main.go
func main() {
    // 1. DB
    db := database.Connect()

    // 2. Repositories
    userRepo := repository.NewUserRepositoryMongo(db)

    // 3. Infrastructure
    validator := infrastructure.NewValidator()

    // 4. Services (con DI)
    userService := services.NewUserService(userRepo, validator)

    // 5. Controllers (con DI)
    userController := controllers.NewUserController(userService)
}
```

### 4. Controllers vs Handlers

**Usar "controllers"**, NO "handlers":

```go
// controllers/user_controller.go
type UserController struct {
    userService *services.UserService  // DI
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{userService: userService}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
    var req dtos.CreateUserRequest
    ctx.ShouldBindJSON(&req)
    user, err := c.userService.CreateUser(ctx.Request.Context(), req)
    ctx.JSON(http.StatusCreated, user)
}
```

### 5. Infrastructure vs Clients

**Usar "infrastructure"**, NO "clients":

- ❌ `internal/clients/` - Nombre incorrecto
- ✅ `internal/infrastructure/` - Nombre correcto

La carpeta `infrastructure` contiene implementaciones de interfaces definidas en `services`.

---

## 📋 Checklist para Implementar un Microservicio

Usa esto para asegurarte de que no olvidas nada:

### Fase 1: Dominio
- [ ] Crear `internal/domain/entities/` con entidades de BD
- [ ] Crear `internal/domain/dtos/` con DTOs para requests/responses

### Fase 2: Repository
- [ ] Crear `internal/repository/*_repository.go` (interfaces)
- [ ] Crear `internal/repository/*_repository_mongo.go` (implementaciones)

### Fase 3: Services
- [ ] Crear `internal/services/*_service.go`
- [ ] Asegurar que dependen de **interfaces** (repositories, validators, publishers)

### Fase 4: Infrastructure (si necesario)
- [ ] Crear `internal/infrastructure/` para servicios externos
- [ ] Implementar interfaces definidas en services

### Fase 5: Controllers
- [ ] Crear `internal/controllers/*_controller.go`
- [ ] Asegurar que dependen de services

### Fase 6: Main
- [ ] Actualizar `cmd/api/main.go` con DI manual
- [ ] Inicializar: DB → Repos → Infrastructure → Services → Controllers
- [ ] Registrar rutas

### Fase 7: Testing
- [ ] Probar con curl/Postman
- [ ] Verificar health check
- [ ] Opcional: Crear tests unitarios

### Fase 8: Documentación
- [ ] Actualizar README del microservicio
- [ ] Documentar endpoints
- [ ] Actualizar `.env.example`

---

## 🚀 Cómo Empezar (Para Equipos)

### Equipo de Payments API:

1. **Estudiar**:
   ```bash
   cd subscriptions-api/
   # Leer README.md
   # Explorar internal/domain/dtos/
   # Explorar internal/repository/
   # Explorar internal/services/
   # Explorar cmd/api/main.go
   ```

2. **Implementar** siguiendo el mismo patrón:
   ```bash
   cd payments-api/
   mkdir -p internal/domain/entities
   mkdir -p internal/domain/dtos
   mkdir -p internal/repository
   mkdir -p internal/controllers
   # Crear archivos siguiendo subscriptions-api como referencia
   ```

3. **Probar**:
   ```bash
   go mod tidy
   go run cmd/api/main.go
   curl http://localhost:8083/healthz
   ```

### Equipo de Search API:

Mismo proceso que Payments API, pero con las particularidades de:
- Repository para búsqueda (in-memory o Solr)
- RabbitMQ Consumer
- CacheService

---

## 📖 Orden de Lectura Recomendado

Para nuevos desarrolladores:

1. **Día 1 - Teoría**:
   - ✅ `LEEME_PRIMERO.md` (este archivo)
   - ✅ `ESTADO_IMPLEMENTACION.md`
   - ✅ `subscriptions-api/README.md`

2. **Día 2 - Práctica**:
   - ✅ Explorar código de `subscriptions-api/`
   - ✅ Entender flujo de DI en `main.go`
   - ✅ Revisar cómo se usan interfaces

3. **Día 3+ - Implementación**:
   - ✅ Implementar tu microservicio asignado
   - ✅ Usar `subscriptions-api` como referencia constante

---

## ⚠️ Errores Comunes a Evitar

### ❌ NO hagas esto:

1. **Mezclar entities y DTOs**:
   ```go
   type User struct {
       ID bson.ObjectID `bson:"_id" json:"id"`  // ❌ NO
   }
   ```

2. **Services accediendo directamente a BD**:
   ```go
   func (s *Service) GetUser() {
       s.db.Collection("users")  // ❌ NO
   }
   ```

3. **Usar "clients" en lugar de "infrastructure"**:
   ```
   internal/clients/  // ❌ NO
   ```

4. **No usar interfaces en services**:
   ```go
   type Service struct {
       repo UserRepositoryMongo  // ❌ NO (implementación concreta)
   }
   ```

5. **No hacer DI manual en main.go**:
   ```go
   service := NewService()  // ❌ NO (crea sus dependencias)
   ```

### ✅ SÍ haz esto:

1. **Separar entities y DTOs**
2. **Services dependen de repository interfaces**
3. **Usar "infrastructure" para servicios externos**
4. **Siempre usar interfaces en dependencias**
5. **DI manual completo en main.go**

---

## 🆘 ¿Necesitas Ayuda?

1. **Primero**: Revisa `subscriptions-api/` - probablemente ya hay un ejemplo
2. **Segundo**: Lee `ESTADO_IMPLEMENTACION.md` - hay instrucciones detalladas
3. **Tercero**: Consulta con el equipo

---

## 📊 Estado Actual (Resumen)

| Microservicio | Completo | Funcional | Arquitectura Limpia | DI |
|---------------|----------|-----------|---------------------|-----|
| users-api | ✅ | ✅ | ❌ | ❌ |
| subscriptions-api | ✅ | ✅ | ✅ | ✅ |
| activities-api | ⚠️ | ✅ | ❌ | ❌ |
| payments-api | ❌ | ❌ | ❌ | ❌ |
| search-api | ❌ | ❌ | ❌ | ❌ |

**Progreso General**: 40% (2/5 microservicios completos)

---

## 🎯 Objetivo Final

Al terminar, **TODOS** los microservicios deben:
- ✅ Tener DTOs separados
- ✅ Usar Repository Pattern con interfaces
- ✅ Implementar Services con DI
- ✅ Tener Controllers (no handlers)
- ✅ Usar Infrastructure para servicios externos
- ✅ DI manual en main.go
- ✅ README documentado
- ✅ Funcionar con Docker

---

**Buena suerte con la implementación! 🚀**

**Recuerda**: `subscriptions-api/` es tu **mejor amigo** - úsalo como referencia constante.
