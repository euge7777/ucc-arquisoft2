# Diagrama de Entidades - Sistema de Gestión de Gimnasio

## Arquitectura de 5 Microservicios

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        FRONTEND (React)                                  │
│  Login | Registro | Home/Búsqueda | Detalles | Inscripción | Admin      │
└─────────────────────────────────────────────────────────────────────────┘
                    │ HTTP/JSON │ HTTP/JSON │ HTTP/JSON
                    ▼           ▼           ▼
    ┌───────────────┴───────────┴───────────┴───────────────┐
    │                                                        │
    ▼                    ▼                    ▼             ▼
┌─────────┐      ┌──────────────┐    ┌──────────────┐  ┌─────────┐
│ users   │      │subscriptions │    │ activities   │  │ search  │
│  api    │◄────►│     api      │◄──►│     api      │  │  api    │
│ (MySQL) │      │  (MongoDB)   │    │   (MySQL)    │  │ (Solr)  │
└─────────┘      └──────┬───────┘    └──────────────┘  └────▲────┘
                        │                                    │
                        ▼                                    │
                 ┌─────────────┐                            │
                 │  payments   │                            │
                 │     api     │                            │
                 │  (MongoDB)  │                            │
                 └─────────────┘                            │
                                                            │
        ┌───────────────────────────────────────────────────┘
        │
    ┌───▼────┐          ┌─────────────┐
    │RabbitMQ│◄────────►│ Memcached   │
    │        │          │   +CCache   │
    └────────┘          └─────────────┘
```

---

## DIAGRAMA ENTIDAD-RELACIÓN DETALLADO

### 1. USERS-API (MySQL)

```
┌─────────────────────────────────────┐
│           USUARIOS                  │
├─────────────────────────────────────┤
│ PK  id (INT AUTO_INCREMENT)         │
│     nombre (VARCHAR 30)             │
│     apellido (VARCHAR 30)           │
│ UK  username (VARCHAR 30)           │
│ UK  email (VARCHAR 100)             │
│     password (CHAR 64) -- SHA256    │
│     is_admin (BOOLEAN)              │
│     sucursal_origen_id (INT) *ref   │
│     fecha_registro (TIMESTAMP)      │
│     created_at (TIMESTAMP)          │
│     updated_at (TIMESTAMP)          │
└─────────────────────────────────────┘
        │
        │ *Referencia lógica (no FK física)
        ▼
    (activities-api: sucursales)
```

---

### 2. SUBSCRIPTIONS-API (MongoDB)

```
┌─────────────────────────────────────────────────────┐
│         PLANES (Collection)                          │
├─────────────────────────────────────────────────────┤
│ _id: ObjectId                                        │
│ nombre: "Plan Básico" | "Plan Premium"              │
│ descripcion: String                                  │
│ precio_mensual: Number (Decimal128)                 │
│ tipo_acceso: "limitado" | "completo"                │
│ duracion_dias: Number (30, 90, 365)                 │
│ activo: Boolean                                      │
│ actividades_permitidas: [String] -- IDs actividades │
│ created_at: Date                                     │
│ updated_at: Date                                     │
└─────────────────────────────────────────────────────┘
                        │
                        │ 1:N
                        ▼
┌──────────────────────────────────────────────────────────┐
│         SUSCRIPCIONES (Collection)                       │
├──────────────────────────────────────────────────────────┤
│ _id: ObjectId                                            │
│ usuario_id: String -- Ref a users-api (HTTP)            │
│ plan_id: ObjectId -- Ref a planes                       │
│ sucursal_origen_id: String -- Ref opcional              │
│ fecha_inicio: Date                                       │
│ fecha_vencimiento: Date                                  │
│ estado: "activa"|"vencida"|"cancelada"|"pendiente_pago" │
│ pago_id: String -- Ref a payments-api                   │
│ metadata: {                                              │
│   auto_renovacion: Boolean,                              │
│   metodo_pago_preferido: String,                         │
│   notas: String                                          │
│ }                                                        │
│ historial_renovaciones: [{                               │
│   fecha: Date,                                           │
│   pago_id: String,                                       │
│   monto: Number                                          │
│ }]                                                       │
│ created_at: Date                                         │
│ updated_at: Date                                         │
└──────────────────────────────────────────────────────────┘
        │                           │
        │ *Referencia HTTP          │ *Referencia lógica
        ▼                           ▼
    (users-api)              (payments-api)
```

**Eventos RabbitMQ publicados:**
- `subscription.created` → {action: "create", type: "subscription", id: "..."}
- `subscription.updated` → {action: "update", type: "subscription", id: "..."}
- `subscription.cancelled` → {action: "delete", type: "subscription", id: "..."}

---

### 3. ACTIVITIES-API (MySQL)

```
┌─────────────────────────────────────┐
│           SUCURSALES                │
├─────────────────────────────────────┤
│ PK  id (INT AUTO_INCREMENT)         │
│     nombre (VARCHAR 100)            │
│     direccion (VARCHAR 255)         │
│     telefono (VARCHAR 20)           │
│     created_at (TIMESTAMP)          │
└─────────────────────────────────────┘
                │
                │ 1:N
                ▼
┌─────────────────────────────────────────────┐
│           ACTIVIDADES                        │
├─────────────────────────────────────────────┤
│ PK  id (INT AUTO_INCREMENT)                 │
│     titulo (VARCHAR 50)                     │
│     descripcion (VARCHAR 255)               │
│     cupo (INT)                              │
│     dia (ENUM: Lunes-Domingo)               │
│     horario_inicio (TIME)                   │
│     horario_final (TIME)                    │
│     foto_url (VARCHAR 511)                  │
│     instructor (VARCHAR 50)                 │
│     categoria (VARCHAR 40)                  │
│ FK  sucursal_id (INT)                       │
│     requiere_plan_premium (BOOLEAN)         │
│     created_at (TIMESTAMP)                  │
│     updated_at (TIMESTAMP)                  │
└─────────────────────────────────────────────┘
                │
                │ 1:N
                ▼
┌──────────────────────────────────────────────────┐
│           INSCRIPCIONES                          │
├──────────────────────────────────────────────────┤
│ PK  id (INT AUTO_INCREMENT)                      │
│     usuario_id (INT) -- Ref a users-api *lógica │
│ FK  actividad_id (INT)                           │
│     suscripcion_id (STRING) -- Ref MongoDB      │
│     fecha_inscripcion (TIMESTAMP)                │
│     is_activa (BOOLEAN)                          │
│     created_at (TIMESTAMP)                       │
│     updated_at (TIMESTAMP)                       │
│                                                  │
│ UK  UNIQUE(usuario_id, actividad_id)            │
└──────────────────────────────────────────────────┘
        │                   │
        │ *Ref HTTP         │ *Ref HTTP
        ▼                   ▼
    (users-api)      (subscriptions-api)
```

**Validaciones antes de crear inscripción:**
1. Usuario existe (HTTP GET users-api/users/:id)
2. Suscripción activa (HTTP GET subscriptions-api/subscriptions/active/:user_id)
3. Plan cubre actividad (si requiere_plan_premium, validar tipo_acceso)
4. Cupo disponible (query local)

**Eventos RabbitMQ publicados:**
- `inscription.created` → {action: "create", type: "inscription", id: "..."}
- `inscription.deleted` → {action: "delete", type: "inscription", id: "..."}

---

### 4. PAYMENTS-API (MongoDB) - GENÉRICO

```
┌────────────────────────────────────────────────────────┐
│         PAYMENTS (Collection)                          │
├────────────────────────────────────────────────────────┤
│ _id: ObjectId                                          │
│ entity_type: String -- "subscription", "inscription",  │
│                         "plan_upgrade", "penalty", ... │
│ entity_id: String -- ID de la entidad (cualquiera)    │
│ user_id: String -- ID del usuario que paga            │
│ amount: Number (Decimal128)                            │
│ currency: String -- "USD", "ARS", "EUR"               │
│ status: "pending"|"completed"|"failed"|"refunded"      │
│ payment_method: String -- "credit_card", "cash", etc   │
│ payment_gateway: String -- "stripe", "mercadopago"     │
│ transaction_id: String -- ID externo del gateway      │
│ metadata: {                                            │
│   -- Información específica del dominio (flexible)     │
│   plan_nombre: String,                                 │
│   duracion_dias: Number,                               │
│   descripcion: String,                                 │
│   ... cualquier dato adicional                         │
│ }                                                      │
│ created_at: Date                                       │
│ updated_at: Date                                       │
│ processed_at: Date -- Fecha de completado             │
└────────────────────────────────────────────────────────┘
```

**Características:**
- ✅ Totalmente desacoplado del dominio
- ✅ Reutilizable en cualquier proyecto
- ✅ Solo necesita saber: quién pagó, cuánto, por qué (entity_type/id)
- ✅ Metadata flexible para información adicional

**Ejemplo uso en gimnasio:**
```json
{
  "entity_type": "subscription",
  "entity_id": "507f1f77bcf86cd799439011",
  "user_id": "123",
  "amount": 100.00,
  "currency": "USD",
  "status": "completed",
  "payment_method": "credit_card",
  "metadata": {
    "plan_nombre": "Plan Premium",
    "duracion_dias": 30,
    "sucursal": "Centro"
  }
}
```

**Ejemplo uso en e-commerce:**
```json
{
  "entity_type": "order",
  "entity_id": "order_12345",
  "user_id": "456",
  "amount": 250.00,
  "currency": "USD",
  "status": "completed",
  "payment_method": "paypal",
  "metadata": {
    "order_items": ["producto1", "producto2"],
    "shipping_address": "Calle Falsa 123"
  }
}
```

---

### 5. SEARCH-API (Apache Solr)

```
┌────────────────────────────────────────────────────────────┐
│         GYM_SEARCH (Solr Core)                             │
├────────────────────────────────────────────────────────────┤
│ id: String (unique) -- Prefijo: "act_", "plan_", "sub_"   │
│ type: String -- "activity", "plan", "subscription"        │
│                                                            │
│ -- Campos de Actividad --                                 │
│ titulo: String (indexed, stored)                          │
│ descripcion: Text (indexed)                               │
│ categoria: String (facet)                                 │
│ instructor: String (facet)                                │
│ dia: String (facet)                                       │
│ horario_inicio: String                                    │
│ horario_final: String                                     │
│ sucursal_id: String (facet)                               │
│ sucursal_nombre: String                                   │
│ requiere_premium: Boolean (facet)                         │
│ cupo_disponible: Int (range query)                        │
│                                                            │
│ -- Campos de Plan --                                      │
│ plan_nombre: String (indexed, stored)                     │
│ plan_precio: Float (range query)                          │
│ plan_tipo_acceso: String (facet)                          │
│                                                            │
│ -- Campos de Suscripción (para admin) --                  │
│ usuario_id: String                                        │
│ estado: String (facet)                                    │
│                                                            │
│ -- Metadata --                                            │
│ created_at: Date (range query)                            │
│ updated_at: Date                                          │
└────────────────────────────────────────────────────────────┘
```

**Consumidores de RabbitMQ:**
- Escucha eventos de `subscriptions-api`
- Escucha eventos de `activities-api`
- Indexa/actualiza/elimina documentos en Solr

**Capas de Caché:**
```
Request → CCache (local, 30s TTL)
            ↓ MISS
          Memcached (distribuido, 60s TTL)
            ↓ MISS
          Solr (búsqueda)
            ↓
          Guarda en Memcached + CCache
            ↓
          Response
```

---

## FLUJOS DE DATOS CLAVE

### Flujo 1: Crear Suscripción
```
┌─────────┐  1. POST /subscriptions    ┌──────────────┐
│Frontend │──────────────────────────►│subscriptions │
└─────────┘                            │     api      │
                                       └──────┬───────┘
                                              │
                     ┌────────────────────────┼────────────────┐
                     │                        │                │
              2. Valida usuario        3. Crea sub      4. Registra pago
                     │                        │                │
                     ▼                        ▼                ▼
              ┌─────────┐             ┌──────────┐      ┌─────────┐
              │ users   │             │ MongoDB  │      │payments │
              │  api    │             │          │      │  api    │
              └─────────┘             └──────────┘      └─────────┘
                                              │
                                       5. Publica evento
                                              │
                                              ▼
                                       ┌──────────┐
                                       │RabbitMQ  │
                                       └────┬─────┘
                                            │
                                       6. Consume e indexa
                                            │
                                            ▼
                                       ┌─────────┐
                                       │ search  │
                                       │  api    │
                                       └─────────┘
```

### Flujo 2: Crear Inscripción
```
┌─────────┐  1. POST /inscripciones   ┌──────────────┐
│Frontend │──────────────────────────►│ activities   │
└─────────┘                            │     api      │
                                       └──────┬───────┘
                                              │
                     ┌────────────────────────┼───────────────────┐
                     │                        │                   │
              2. Valida usuario      3. Valida suscripción  4. Crea inscripción
                     │                        │                   │
                     ▼                        ▼                   ▼
              ┌─────────┐             ┌──────────────┐     ┌─────────┐
              │ users   │             │subscriptions │     │ MySQL   │
              │  api    │             │     api      │     │         │
              └─────────┘             └──────────────┘     └─────────┘
                                              │
                                       5. Publica evento
                                              │
                                              ▼
                                       ┌──────────┐
                                       │RabbitMQ  │
                                       └────┬─────┘
                                            │
                                       6. Actualiza índice
                                            │
                                            ▼
                                       ┌─────────┐
                                       │ search  │
                                       │  api    │
                                       └─────────┘
```

### Flujo 3: Búsqueda de Actividades
```
┌─────────┐  1. GET /search?q=yoga    ┌─────────┐
│Frontend │──────────────────────────►│ search  │
└─────────┘                            │  api    │
                                       └────┬────┘
                                            │
                                       2. Busca caché
                                            │
                     ┌──────────────────────┼──────────────┐
                     │                      │              │
              3. CCache (local)      4. Memcached    5. Solr (source)
                     │                      │              │
                     ▼                      ▼              ▼
              ┌─────────┐             ┌──────────┐  ┌─────────┐
              │In-Memory│             │Memcached │  │  Solr   │
              │30s TTL  │             │60s TTL   │  │ Search  │
              └─────────┘             └──────────┘  └─────────┘
                                            │
                                       6. Guarda en cachés
                                            │
                                            ▼
                                       7. Response JSON
```

---

## RELACIONES CROSS-MICROSERVICIO

### Referencias Lógicas (HTTP):
```
users-api ──────────────┐
    ▲                   │ Valida usuario existe
    │                   │
    │                   ▼
    │            subscriptions-api ──────┐
    │                   ▲                │ Valida suscripción activa
    │                   │                │
    │                   │                ▼
    │                   │         activities-api
    │                   │
    └───────────────────┘
         Valida usuario
```

### Comunicación Asíncrona (RabbitMQ):
```
subscriptions-api ───┐
                     │ Publica eventos
                     ├──► subscription.created
                     ├──► subscription.updated
                     └──► subscription.cancelled
                          │
                          ▼
                     ┌──────────┐
                     │RabbitMQ  │
                     │ Exchange │
                     └────┬─────┘
                          │
                          ▼
                     search-api (Consumer)
                          │
                          ├──► Indexa en Solr
                          └──► Invalida caché

activities-api ──────┐
                     │ Publica eventos
                     ├──► inscription.created
                     └──► inscription.deleted
                          │
                          ▼
                     ┌──────────┐
                     │RabbitMQ  │
                     │ Exchange │
                     └────┬─────┘
                          │
                          ▼
                     search-api (Consumer)
                          │
                          └──► Actualiza cupos en Solr
```

---

## RESUMEN DE ENTIDADES POR BASE DE DATOS

### MySQL (users-api):
- **usuarios** (1 tabla)

### MySQL (activities-api):
- **sucursales** (1 tabla)
- **actividades** (1 tabla)
- **inscripciones** (1 tabla)

### MongoDB (subscriptions-api):
- **planes** (1 colección)
- **suscripciones** (1 colección)

### MongoDB (payments-api):
- **payments** (1 colección - genérica)

### Solr (search-api):
- **gym_search** (1 core - documentos mixtos con campo "type")

**Total: 8 entidades distribuidas en 5 microservicios**
