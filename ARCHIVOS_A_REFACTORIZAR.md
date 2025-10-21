# Archivos Viejos que HAY QUE ELIMINAR/REFACTORIZAR

## ⚠️ IMPORTANTE: Estructura Vieja vs Nueva

Varios microservicios todavía tienen la **estructura VIEJA** con:
- ❌ `internal/handlers/` (debería ser `controllers/`)
- ❌ `internal/clients/` (debería ser `infrastructure/`)
- ❌ `internal/models/` mezclando entities y DTOs

---

## 1. users-api ⚠️

### Archivos con Estructura VIEJA:

```
users-api/internal/
├── handlers/              ❌ ELIMINAR - Reemplazar por controllers/
│   ├── user_handler.go
│   └── auth_handler.go
│
├── models/                ❌ REFACTORIZAR - Separar en entities/ y dtos/
│   └── user.go
│
└── clients/               ❌ NO EXISTE (está bien, no tiene clients)
```

### ✅ Lo que HAY QUE HACER:

```bash
cd users-api

# 1. Crear nueva estructura
mkdir -p internal/domain/entities
mkdir -p internal/domain/dtos
mkdir -p internal/controllers

# 2. Mover y refactorizar
# - Copiar internal/models/user.go a internal/domain/entities/user.go
# - Crear internal/domain/dtos/user_dtos.go (separar DTOs)
# - Refactorizar internal/handlers/ → internal/controllers/
# - Actualizar cmd/api/main.go con DI manual

# 3. Eliminar archivos viejos
rm -rf internal/handlers/
rm -rf internal/models/
```

### Decisión Recomendada:
**NO refactorizar users-api** - Ya funciona y no es prioridad. Dejarlo como está y enfocarse en los que faltan.

---

## 2. activities-api ⚠️

### Archivos con Estructura VIEJA:

```
activities-api/internal/
├── handlers/              ❌ ELIMINAR - Reemplazar por controllers/
│   ├── sucursal_handler.go
│   ├── actividad_handler.go
│   └── inscripcion_handler.go
│
├── models/                ❌ REFACTORIZAR - Separar en entities/ y dtos/
│   ├── sucursal.go
│   ├── actividad.go
│   └── inscripcion.go
│
├── clients/               ❌ NO EXISTE (pero falta RabbitMQ publisher)
│
├── repository/            ✅ MANTENER (ya tiene repository pattern)
│   ├── sucursal_repository.go
│   ├── actividad_repository.go
│   └── inscripcion_repository.go
│
└── services/              ✅ MANTENER (pero refactorizar para DI)
    ├── sucursal_service.go
    ├── actividad_service.go
    └── inscripcion_service.go
```

### ✅ Lo que HAY QUE HACER:

```bash
cd activities-api

# 1. Crear nueva estructura
mkdir -p internal/domain/entities
mkdir -p internal/domain/dtos
mkdir -p internal/controllers
mkdir -p internal/infrastructure

# 2. Mover y refactorizar
# - Mover internal/models/*.go → internal/domain/entities/
# - Crear internal/domain/dtos/ con DTOs separados
# - Refactorizar internal/handlers/ → internal/controllers/
# - Crear internal/infrastructure/rabbitmq_publisher.go
# - Actualizar services/ para aceptar interfaces (DI)
# - Actualizar cmd/api/main.go con DI manual

# 3. Eliminar archivos viejos
rm -rf internal/handlers/
rm -rf internal/models/
```

### Prioridad:
**MEDIA** - Funciona pero falta RabbitMQ y arquitectura limpia.

---

## 3. subscriptions-api ✅

### Estado: ✅ **YA ESTÁ REFACTORIZADO**

```
subscriptions-api/internal/
├── domain/
│   ├── entities/          ✅ Correcto
│   └── dtos/              ✅ Correcto
├── repository/            ✅ Correcto (con interfaces)
├── services/              ✅ Correcto (con DI)
├── infrastructure/        ✅ Correcto (no "clients")
├── controllers/           ✅ Correcto (no "handlers")
├── middleware/            ✅ Correcto
├── database/              ✅ Correcto
└── config/                ✅ Correcto
```

**NO tocar** - Este es el **EJEMPLO DE REFERENCIA**.

---

## 4. payments-api ⚠️

### Archivos con Estructura VIEJA:

```
payments-api/internal/
├── handlers/              ❌ ELIMINAR - Reemplazar por controllers/
│   └── payment_handler.go
│
├── models/                ❌ REFACTORIZAR - Separar en entities/ y dtos/
│   └── payment.go
│
├── services/              ❌ REFACTORIZAR - Implementar con DI
│   └── payment_service.go
│
└── clients/               ❌ NO EXISTE (no tiene clients externos)
```

### ✅ Lo que HAY QUE HACER:

```bash
cd payments-api

# 1. ELIMINAR archivos viejos (son solo estructura básica)
rm -rf internal/handlers/
rm -rf internal/models/
rm -rf internal/services/

# 2. Crear estructura nueva DESDE CERO siguiendo subscriptions-api
mkdir -p internal/domain/entities
mkdir -p internal/domain/dtos
mkdir -p internal/repository
mkdir -p internal/services
mkdir -p internal/controllers

# 3. Implementar desde cero usando subscriptions-api como referencia
# Ver ESTADO_IMPLEMENTACION.md para pasos detallados
```

### Prioridad:
**ALTA** ⚠️ - TODO por hacer. Empezar desde cero siguiendo `subscriptions-api`.

---

## 5. search-api ⚠️

### Archivos con Estructura VIEJA:

```
search-api/internal/
├── handlers/              ❌ ELIMINAR - Reemplazar por controllers/
│   └── search_handler.go
│
├── models/                ❌ REFACTORIZAR - Separar en entities/ y dtos/
│   └── search.go
│
├── services/              ⚠️ REFACTORIZAR - Implementar con DI
│   ├── search_service.go
│   └── cache_service.go
│
├── clients/               ❌ ELIMINAR - Reemplazar por infrastructure/
│   └── rabbitmq_client.go
│
└── consumers/             ⚠️ REFACTORIZAR - Implementar con DI
    └── rabbitmq_consumer.go
```

### ✅ Lo que HAY QUE HACER:

```bash
cd search-api

# 1. ELIMINAR archivos viejos
rm -rf internal/handlers/
rm -rf internal/models/
rm -rf internal/clients/

# 2. Crear estructura nueva
mkdir -p internal/domain/entities
mkdir -p internal/domain/dtos
mkdir -p internal/repository
mkdir -p internal/controllers
mkdir -p internal/infrastructure

# 3. Refactorizar lo existente
# - Mover consumers/ a infrastructure/ (es parte de infraestructura)
# - Refactorizar services/ con DI
# - Crear repository pattern para búsqueda
# - Crear controllers/
# - Actualizar main.go con DI

# 4. Estructura final:
# internal/
#   ├── domain/entities/
#   ├── domain/dtos/
#   ├── repository/              # SearchRepository interface
#   ├── services/                # SearchService con DI
#   ├── infrastructure/          # RabbitMQConsumer, CacheService
#   ├── controllers/
#   └── ...
```

### Prioridad:
**ALTA** ⚠️ - Requiere refactorización completa.

---

## 📋 Plan de Acción por Equipo

### Equipo 1: users-api
**Decisión**: ✅ **NO REFACTORIZAR**
- Ya funciona correctamente
- No es prioridad
- Dejar `handlers/` y `models/` como están

### Equipo 2: subscriptions-api
**Estado**: ✅ **YA COMPLETO**
- Es el ejemplo de referencia
- No tocar

### Equipo 3: activities-api
**Decisión**: ⚠️ **Opcional - BAJA PRIORIDAD**
- Ya funciona
- Solo falta agregar RabbitMQ publisher
- La refactorización completa es opcional

**Si deciden refactorizar**:
1. Crear `internal/domain/entities/` y mover `models/`
2. Crear `internal/domain/dtos/`
3. Renombrar `handlers/` a `controllers/`
4. Crear `internal/infrastructure/rabbitmq_publisher.go`
5. Actualizar `main.go` con DI

### Equipo 4: payments-api
**Decisión**: ⚠️ **IMPLEMENTAR DESDE CERO - ALTA PRIORIDAD**

**Pasos**:
1. ❌ Eliminar `internal/handlers/`, `internal/models/`, `internal/services/`
2. ✅ Copiar estructura de `subscriptions-api/internal/`
3. ✅ Adaptar a Payment (en vez de Plan/Subscription)
4. ✅ Implementar siguiendo el README de subscriptions-api

**Archivos a crear** (usar subscriptions-api como base):
```
internal/
├── domain/
│   ├── entities/payment.go
│   └── dtos/payment_dtos.go
├── repository/
│   ├── payment_repository.go
│   └── payment_repository_mongo.go
├── services/
│   └── payment_service.go
├── controllers/
│   └── payment_controller.go
└── ...
```

### Equipo 5: search-api
**Decisión**: ⚠️ **REFACTORIZAR COMPLETO - ALTA PRIORIDAD**

**Pasos**:
1. ❌ Eliminar `internal/handlers/`, `internal/models/`, `internal/clients/`
2. ✅ Crear estructura nueva siguiendo subscriptions-api
3. ⚠️ Refactorizar `services/` con DI
4. ⚠️ Mover `consumers/` a `infrastructure/`
5. ✅ Crear `repository/` para abstracción de búsqueda
6. ✅ Crear `controllers/`
7. ✅ Actualizar `main.go` con DI

---

## 🔍 Verificación Rápida

Para saber si un microservicio está correctamente refactorizado:

### ❌ Estructura VIEJA (INCORRECTA):
```
internal/
├── handlers/          ❌
├── models/            ❌
├── clients/           ❌
└── services/          ⚠️ (sin DI)
```

### ✅ Estructura NUEVA (CORRECTA):
```
internal/
├── domain/
│   ├── entities/      ✅
│   └── dtos/          ✅
├── repository/        ✅ (con interfaces)
├── services/          ✅ (con DI)
├── infrastructure/    ✅ (no "clients")
├── controllers/       ✅ (no "handlers")
├── middleware/
├── database/
└── config/
```

---

## 📊 Resumen de Estado Real

| Microservicio | handlers | clients | models | Necesita Refactorizar |
|---------------|----------|---------|--------|----------------------|
| users-api | ❌ Tiene | ✅ No tiene | ❌ Tiene | ⚠️ Opcional (NO hacerlo) |
| subscriptions-api | ✅ No tiene | ✅ No tiene | ✅ No tiene | ✅ YA ESTÁ BIEN |
| activities-api | ❌ Tiene | ✅ No tiene | ❌ Tiene | ⚠️ Opcional |
| payments-api | ❌ Tiene | ✅ No tiene | ❌ Tiene | ❌ SÍ - DESDE CERO |
| search-api | ❌ Tiene | ❌ Tiene | ❌ Tiene | ❌ SÍ - REFACTORIZAR |

---

## 🎯 Prioridades Finales

1. **ALTA** - `payments-api`: Implementar desde cero
2. **ALTA** - `search-api`: Refactorizar completo
3. **MEDIA** - `activities-api`: Solo agregar RabbitMQ (refactorización opcional)
4. **BAJA** - `users-api`: No tocar (ya funciona)

---

**Conclusión**: Solo `subscriptions-api` tiene la arquitectura correcta. Los demás tienen archivos viejos que hay que eliminar/refactorizar.
