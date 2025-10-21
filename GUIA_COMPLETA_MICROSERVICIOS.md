# Guía Completa de Implementación - Arquitectura de Microservicios

## Resumen del Sistema

El sistema de gestión de gimnasio está compuesto por **5 microservicios independientes**:

| Microservicio | Puerto | Base de Datos | Función Principal |
|---------------|--------|---------------|-------------------|
| `users-api` | 8080 | MySQL | Autenticación y gestión de usuarios |
| `subscriptions-api` | 8081 | MongoDB | Planes y suscripciones |
| `activities-api` | 8082 | MySQL | Actividades e inscripciones |
| `payments-api` | 8083 | MongoDB | Pagos (genérico, reutilizable) |
| `search-api` | 8084 | In-Memory/Solr | Búsqueda con caché |

## Infraestructura Completa

```
┌─────────────────────────────────────────────────────────────────┐
│                         FRONTEND                                 │
│                     (React / Next.js)                            │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTP/JSON
        ┌────────────────┼────────────────┐
        │                │                │
        ▼                ▼                ▼
   ┌─────────┐     ┌──────────┐    ┌──────────┐     ┌─────────┐
   │ users   │     │subscrip- │    │activities│     │ search  │
   │  api    │◄───►│ tions    │◄──►│   api    │     │  api    │
   │ :8080   │     │  api     │    │  :8082   │     │ :8084   │
   │ MySQL   │     │ :8081    │    │  MySQL   │     │In-Memory│
   └─────────┘     │ MongoDB  │    └──────────┘     └────▲────┘
                   └─────┬────┘                          │
                         │                               │
                         ▼                               │
                   ┌─────────┐                          │
                   │payments │                          │
                   │  api    │                          │
                   │ :8083   │                          │
                   │MongoDB  │                          │
                   └─────────┘                          │
                                                        │
        ┌───────────────────────────────────────────────┘
        │
    ┌───▼────┐          ┌─────────┐
    │RabbitMQ│◄────────►│Memcached│
    │ :5672  │          │ :11211  │
    └────────┘          └─────────┘
```

## 1. Microservicios ya Implementados

### ✅ users-api (Completo)

**Ubicación**: `users-api/`

**Características**:
- Registro y login con JWT
- CRUD de usuarios
- Middleware de autenticación
- Validación de contraseñas fuertes

**Endpoints Principales**:
```
POST   /register          - Registrar usuario
POST   /login             - Login con JWT
GET    /users             - Listar usuarios (admin)
GET    /users/:id         - Obtener usuario
PUT    /users/:id         - Actualizar usuario
DELETE /users/:id         - Eliminar usuario
GET    /healthz           - Health check
```

**Testing**:
```bash
# Registrar usuario
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan",
    "apellido": "Pérez",
    "username": "juanp",
    "email": "juan@example.com",
    "password": "Password123"
  }'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "juanp",
    "password": "Password123"
  }'
```

---

### ✅ activities-api (Migrado)

**Ubicación**: `activities-api/`

**Características**:
- Gestión de sucursales
- Gestión de actividades
- Inscripciones con validaciones
- Arquitectura limpia (handlers, services, repositories)

**Endpoints Principales**:
```
# Sucursales
GET    /sucursales
POST   /sucursales
GET    /sucursales/:id

# Actividades
GET    /actividades
POST   /actividades
GET    /actividades/:id
PUT    /actividades/:id
DELETE /actividades/:id

# Inscripciones
POST   /inscripciones
GET    /inscripciones/usuario/:id
DELETE /inscripciones/:id
```

**Testing**:
```bash
# Crear sucursal
curl -X POST http://localhost:8082/sucursales \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Sede Centro",
    "direccion": "Av. Principal 123",
    "telefono": "555-1234"
  }'

# Crear actividad
curl -X POST http://localhost:8082/actividades \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Yoga Matutino",
    "descripcion": "Clase de yoga para todos los niveles",
    "cupo": 20,
    "dia": "Lunes",
    "horario_inicio": "08:00:00",
    "horario_final": "09:00:00",
    "sucursal_id": 1,
    "instructor": "María López",
    "categoria": "Fitness"
  }'
```

---

## 2. Nuevos Microservicios Creados

### 🆕 subscriptions-api

**Ubicación**: `subscriptions-api/`

**Características**:
- Gestión de planes (básico, premium)
- Suscripciones de usuarios
- Validación con `users-api`
- Publicación de eventos a RabbitMQ
- MongoDB para flexibilidad

**Endpoints**:
```
# Planes
POST   /plans              - Crear plan
GET    /plans              - Listar planes
GET    /plans/:id          - Obtener plan

# Suscripciones
POST   /subscriptions      - Crear suscripción
GET    /subscriptions/:id  - Obtener suscripción
GET    /subscriptions/active/:user_id  - Suscripción activa del usuario
PATCH  /subscriptions/:id/status       - Actualizar estado
DELETE /subscriptions/:id  - Cancelar suscripción
```

**Flujo de Creación de Suscripción**:
```
1. Usuario hace POST /subscriptions
2. subscriptions-api valida usuario con users-api
3. subscriptions-api valida plan existe y está activo
4. Calcula fecha_vencimiento = fecha_inicio + plan.duracion_dias
5. Crea suscripción con estado "pendiente_pago"
6. Publica evento a RabbitMQ: subscription.create
7. search-api indexa la suscripción
```

**Testing**:
```bash
# Crear plan
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Premium",
    "descripcion": "Acceso completo a todas las actividades",
    "precio_mensual": 100.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true
  }'

# Crear suscripción
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": "5",
    "plan_id": "507f1f77bcf86cd799439011",
    "metodo_pago": "credit_card"
  }'
```

---

### 🆕 payments-api (Genérico)

**Ubicación**: `payments-api/`

**Características**:
- **100% Agnóstico del dominio**
- Sirve para gimnasio, e-commerce, SaaS, etc.
- Campo `metadata` flexible
- Estados: pending, completed, failed, refunded

**Endpoints**:
```
POST   /payments                    - Crear pago
GET    /payments/:id                - Obtener pago
GET    /payments/user/:user_id      - Pagos de usuario
GET    /payments/entity?entity_type=X&entity_id=Y  - Pagos de entidad
GET    /payments/status?status=pending             - Pagos por estado
PATCH  /payments/:id/status         - Actualizar estado
POST   /payments/:id/process        - Procesar pago
```

**Modelo Genérico**:
```json
{
  "entity_type": "subscription",  // Qué se está pagando
  "entity_id": "507f...",         // ID de esa entidad
  "user_id": "5",                 // Quién paga
  "amount": 100.00,
  "currency": "USD",
  "status": "pending",
  "payment_method": "credit_card",
  "metadata": {                   // Información adicional
    "plan_nombre": "Premium",
    "duracion_dias": 30
  }
}
```

**Testing**:
```bash
# Crear pago
curl -X POST http://localhost:8083/payments \
  -H "Content-Type: application/json" \
  -d '{
    "entity_type": "subscription",
    "entity_id": "507f1f77bcf86cd799439011",
    "user_id": "5",
    "amount": 100.00,
    "currency": "USD",
    "payment_method": "credit_card",
    "metadata": {
      "plan_nombre": "Plan Premium"
    }
  }'

# Procesar pago (simula aprobación)
curl -X POST http://localhost:8083/payments/65a7b3c1d2e3f4g5h6i7j8k9/process
```

---

### 🆕 search-api

**Ubicación**: `search-api/`

**Características**:
- Búsqueda avanzada con filtros
- **Caché de dos niveles**:
  - CCache local (in-memory): 30s TTL
  - Memcached distribuido: 60s TTL
- Consumidor de RabbitMQ para indexación automática
- Implementación in-memory (migrable a Solr)

**Endpoints**:
```
POST   /search             - Búsqueda avanzada
GET    /search?q=yoga      - Búsqueda rápida
GET    /search/stats       - Estadísticas del índice
GET    /search/:id         - Obtener documento
POST   /search/index       - Indexar manualmente
DELETE /search/:id         - Eliminar documento
```

**Flujo de Búsqueda con Caché**:
```
1. Request → search-api
2. Busca en CCache local (30s)
   ├─ HIT → Return + Header "X-Cache: HIT"
   └─ MISS → 3
3. Busca en Memcached (60s)
   ├─ HIT → Guarda en CCache → Return
   └─ MISS → 4
4. Ejecuta búsqueda real
5. Guarda en Memcached + CCache
6. Return + Header "X-Cache: MISS"
```

**Testing**:
```bash
# Búsqueda rápida
curl "http://localhost:8084/search?q=yoga&type=activity"

# Búsqueda avanzada
curl -X POST http://localhost:8084/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "yoga",
    "type": "activity",
    "filters": {
      "categoria": "fitness",
      "dia": "Lunes"
    },
    "page": 1,
    "page_size": 10
  }'

# Ver estadísticas
curl http://localhost:8084/search/stats
```

**Eventos Consumidos**:
- `activity.create/update/delete`
- `plan.create/update`
- `subscription.create/update/delete`

---

## 3. Comunicación entre Microservicios

### HTTP Síncrono (Validaciones)

```go
// subscriptions-api valida usuario con users-api
usersClient := NewUsersClient("http://users-api:8080")
valid, err := usersClient.ValidateUser(userID)

// activities-api valida suscripción activa
subsClient := NewSubscriptionsClient("http://subscriptions-api:8081")
subscription, err := subsClient.GetActiveSubscription(userID)
```

### RabbitMQ Asíncrono (Eventos)

**Publicadores**:
- `subscriptions-api` → publica: `plan.*`, `subscription.*`
- `activities-api` → publica: `activity.*`, `inscription.*`

**Consumidores**:
- `search-api` → consume todos los eventos y actualiza índice

**Formato de Evento**:
```json
{
  "action": "create",
  "type": "subscription",
  "id": "507f1f77bcf86cd799439011",
  "timestamp": "2025-01-15T10:00:00Z",
  "data": {
    "usuario_id": "5",
    "plan_id": "...",
    "estado": "activa"
  }
}
```

---

## 4. Ejecución del Sistema Completo

### Usando Docker Compose

```bash
# Levantar todos los servicios
docker-compose -f docker-compose.new.yml up -d

# Ver logs
docker-compose -f docker-compose.new.yml logs -f

# Ver estado de servicios
docker-compose -f docker-compose.new.yml ps

# Detener todo
docker-compose -f docker-compose.new.yml down
```

### Ejecución Local (Desarrollo)

**Terminal 1: users-api**
```bash
cd users-api
go mod tidy
go run cmd/api/main.go
# Corre en :8080
```

**Terminal 2: subscriptions-api**
```bash
cd subscriptions-api
go mod tidy
go run cmd/api/main.go
# Corre en :8081
```

**Terminal 3: activities-api**
```bash
cd activities-api
go mod tidy
go run cmd/api/main.go
# Corre en :8082
```

**Terminal 4: payments-api**
```bash
cd payments-api
go mod tidy
go run cmd/api/main.go
# Corre en :8083
```

**Terminal 5: search-api**
```bash
cd search-api
go mod tidy
go run cmd/api/main.go
# Corre en :8084
```

---

## 5. Health Checks

Verificar que todos los servicios estén corriendo:

```bash
# users-api
curl http://localhost:8080/healthz

# subscriptions-api
curl http://localhost:8081/healthz

# activities-api
curl http://localhost:8082/healthz

# payments-api
curl http://localhost:8083/healthz

# search-api
curl http://localhost:8084/healthz
```

---

## 6. Flujo Completo de Uso

### Caso 1: Usuario se Registra y Crea Suscripción

```bash
# 1. Registrar usuario
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Ana",
    "apellido": "García",
    "username": "anag",
    "email": "ana@example.com",
    "password": "Password123"
  }'
# Respuesta: { "id_usuario": 10, ... }

# 2. Crear plan (admin)
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Mensual",
    "precio_mensual": 50.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true
  }'
# Respuesta: { "id": "65a...", ... }

# 3. Crear suscripción
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": "10",
    "plan_id": "65a...",
    "metodo_pago": "credit_card"
  }'
# Respuesta: { "id": "65b...", "estado": "pendiente_pago", ... }

# 4. Crear pago
curl -X POST http://localhost:8083/payments \
  -H "Content-Type: application/json" \
  -d '{
    "entity_type": "subscription",
    "entity_id": "65b...",
    "user_id": "10",
    "amount": 50.00,
    "currency": "USD",
    "payment_method": "credit_card"
  }'
# Respuesta: { "id": "65c...", "status": "pending", ... }

# 5. Procesar pago
curl -X POST http://localhost:8083/payments/65c.../process
# Respuesta: { "message": "Pago procesado correctamente" }

# 6. Actualizar estado de suscripción
curl -X PATCH http://localhost:8081/subscriptions/65b.../status \
  -H "Content-Type: application/json" \
  -d '{
    "estado": "activa",
    "pago_id": "65c..."
  }'
# Respuesta: { "message": "Estado actualizado correctamente" }
```

### Caso 2: Usuario se Inscribe a Actividad

```bash
# 1. Buscar actividades de yoga
curl "http://localhost:8084/search?q=yoga&type=activity"

# 2. Inscribirse a actividad
curl -X POST http://localhost:8082/inscripciones \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 10,
    "actividad_id": 5
  }'
# Respuesta: { "id": 20, "usuario_id": 10, "actividad_id": 5, ... }
```

---

## 7. Arquitectura de Directorios

```
ucc-arquisoft2/
├── users-api/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── subscriptions-api/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── services/
│   │   └── clients/ (RabbitMQ, Users API)
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── activities-api/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── payments-api/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   └── services/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── search-api/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handlers/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── services/ (search, cache)
│   │   └── consumers/ (RabbitMQ)
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── docker-compose.new.yml
├── DIAGRAMA_ENTIDADES.md
├── ARQUITECTURA_MICROSERVICIOS.md
└── GUIA_COMPLETA_MICROSERVICIOS.md (este archivo)
```

---

## 8. Migraciones y Próximos Pasos

### Migrar search-api a Apache Solr

Actualmente usa implementación in-memory. Para producción:

1. **Habilitar Solr en docker-compose**:
   ```yaml
   # Descomentar servicio 'solr' en docker-compose.new.yml
   ```

2. **Reemplazar SearchService**:
   ```go
   // internal/services/search_service.go
   import "github.com/rtt/Go-Solr"

   solr := gosolr.NewSolrInterface("http://solr:8983/solr", "gym_search")
   ```

3. **El resto de la arquitectura permanece igual** (caché, consumer, handlers)

### Implementar Frontend

El frontend debería comunicarse con todos los microservicios:

```javascript
// config.js
export const API_URLS = {
  users: 'http://localhost:8080',
  subscriptions: 'http://localhost:8081',
  activities: 'http://localhost:8082',
  payments: 'http://localhost:8083',
  search: 'http://localhost:8084'
}
```

### Agregar API Gateway (Opcional)

Para simplificar el frontend, considera usar:
- **Kong**
- **Traefik**
- **NGINX** como reverse proxy

---

## 9. Troubleshooting

### Problema: Microservicio no se conecta a la base de datos

```bash
# Verificar que la base de datos esté corriendo
docker ps | grep mysql
docker ps | grep mongodb

# Ver logs del microservicio
docker logs gym-users-api
```

### Problema: RabbitMQ no recibe eventos

```bash
# Verificar que RabbitMQ esté corriendo
docker logs gym-rabbitmq

# Acceder a la UI de RabbitMQ
# http://localhost:15672 (usuario: admin, password: admin)

# Verificar que el exchange existe
# Exchanges → gym_events
```

### Problema: Caché no funciona

```bash
# Verificar Memcached
docker logs gym-memcached

# Probar conexión
telnet localhost 11211
```

---

## 10. Resumen

✅ **5 Microservicios Implementados**
- users-api
- subscriptions-api
- activities-api
- payments-api
- search-api

✅ **Infraestructura Completa**
- MySQL (users, activities)
- MongoDB (subscriptions, payments)
- RabbitMQ (eventos)
- Memcached (caché)
- Solr (opcional, search usa in-memory)

✅ **Patrones Implementados**
- Arquitectura limpia (handlers, services, repositories)
- Comunicación síncrona (HTTP)
- Comunicación asíncrona (RabbitMQ)
- Caché de dos niveles
- Health checks
- Docker & Docker Compose

✅ **Listo para Desarrollo y Producción**
