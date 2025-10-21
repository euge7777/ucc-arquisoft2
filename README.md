# Sistema de Gestión de Gimnasio - Arquitectura de Microservicios

Sistema de gestión de gimnasio implementado con **arquitectura de microservicios** en Go.

## 🚨 LEE ESTO PRIMERO

**⚠️ Este proyecto está PARCIALMENTE IMPLEMENTADO**

### Estado Real:

| Microservicio | ¿Funciona? | Arquitectura Correcta | Archivos Viejos |
|---------------|------------|----------------------|-----------------|
| users-api | ✅ SÍ | ❌ NO | ⚠️ `handlers/`, `models/` |
| subscriptions-api | ✅ SÍ | ✅ **SÍ** (Ejemplo) | ✅ Ninguno |
| activities-api | ✅ SÍ | ❌ NO | ⚠️ `handlers/`, `models/` |
| payments-api | ❌ NO | ❌ NO | ⚠️ Estructura vieja |
| search-api | ❌ NO | ❌ NO | ⚠️ Estructura vieja |

**Solo 1 de 5 microservicios** tiene arquitectura limpia completa: **`subscriptions-api`**

---

## 📖 Documentación (Leer en ESTE orden)

### 🎯 Documentos CRÍTICOS:

1. **[`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)** ⭐⭐⭐
   - **EMPIEZA AQUÍ**
   - Estado real de cada microservicio
   - Qué funciona y qué no

2. **[`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)** ⭐⭐
   - Archivos viejos a eliminar
   - Instrucciones por microservicio

3. **[`subscriptions-api/README.md`](subscriptions-api/README.md)** ⭐⭐
   - **Ejemplo de referencia**
   - Arquitectura limpia explicada
   - DTOs, Repository, DI

### 📚 Documentos de Referencia:

4. [`LEEME_PRIMERO.md`](LEEME_PRIMERO.md) - Conceptos y arquitectura
5. [`ESTADO_IMPLEMENTACION.md`](ESTADO_IMPLEMENTACION.md) - Estado detallado
6. [`ARQUITECTURA_MICROSERVICIOS.md`](ARQUITECTURA_MICROSERVICIOS.md) - Patrones
7. [`DIAGRAMA_ENTIDADES.md`](DIAGRAMA_ENTIDADES.md) - Modelo de datos

---

## 🏗️ Arquitectura

```
Frontend (React)
     │
     ├─→ users-api (8080)        ✅ Funciona (estructura vieja)
     ├─→ subscriptions-api (8081) ✅ Funciona (arquitectura correcta) ⭐
     ├─→ activities-api (8082)   ✅ Funciona (estructura vieja)
     ├─→ payments-api (8083)     ❌ No funciona
     └─→ search-api (8084)       ❌ No funciona
          │
          ├─→ RabbitMQ (5672)
          └─→ Memcached (11211)
```

### Bases de Datos:
- **MySQL** - users-api, activities-api
- **MongoDB** - subscriptions-api, payments-api

---

## 🚀 Inicio Rápido

### 1. Levantar Infraestructura

```bash
# Levantar bases de datos y RabbitMQ
docker-compose -f docker-compose.new.yml up -d mysql mongodb rabbitmq memcached
```

### 2. Ejecutar Microservicios Funcionales

```bash
# users-api
cd users-api
go run cmd/api/main.go  # Puerto 8080

# subscriptions-api (EJEMPLO DE REFERENCIA)
cd subscriptions-api
go run cmd/api/main.go  # Puerto 8081

# activities-api
cd activities-api
go run cmd/api/main.go  # Puerto 8082
```

### 3. Verificar

```bash
curl http://localhost:8080/healthz  # users-api
curl http://localhost:8081/healthz  # subscriptions-api
curl http://localhost:8082/healthz  # activities-api
```

---

## ⚠️ Microservicios que NO funcionan

### payments-api ❌
**Estado**: Solo estructura básica con archivos viejos

**Acción requerida**:
1. Leer [`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)
2. Leer [`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)
3. Eliminar archivos viejos
4. Implementar desde cero usando `subscriptions-api` como base

### search-api ❌
**Estado**: Solo estructura básica con archivos viejos

**Acción requerida**:
1. Leer [`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)
2. Leer [`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)
3. Eliminar archivos viejos (`handlers/`, `clients/`, `models/`)
4. Refactorizar usando `subscriptions-api` como base

---

## 📦 Arquitectura Correcta (subscriptions-api)

**Solo este microservicio tiene la implementación correcta:**

```
subscriptions-api/
├── cmd/api/main.go                # ✅ DI manual
├── internal/
│   ├── domain/
│   │   ├── entities/              # ✅ Entidades de BD
│   │   └── dtos/                  # ✅ DTOs API
│   ├── repository/                # ✅ Interfaces + Mongo
│   ├── services/                  # ✅ Lógica con DI
│   ├── infrastructure/            # ✅ Servicios externos
│   └── controllers/               # ✅ HTTP handlers
```

**Ver [`subscriptions-api/README.md`](subscriptions-api/README.md) para detalles.**

---

## ❌ Estructura Incorrecta (resto)

**users-api, activities-api, payments-api, search-api tienen:**

```
microservicio-api/
├── internal/
│   ├── handlers/      ❌ Debería ser "controllers"
│   ├── models/        ❌ Debería estar separado en entities/ y dtos/
│   ├── clients/       ❌ Debería ser "infrastructure"
│   └── services/      ⚠️ Sin DI
```

**Ver [`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md) para detalles.**

---

## 🎯 Para Equipos de Desarrollo

### Si vas a trabajar en payments-api o search-api:

1. **NO ejecutes `go run` directamente** - no funcionará
2. **Primero lee**:
   - [`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)
   - [`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)
   - [`subscriptions-api/README.md`](subscriptions-api/README.md)
3. **Elimina archivos viejos** según las instrucciones
4. **Implementa desde cero** usando `subscriptions-api` como referencia

### Si vas a trabajar en users-api o activities-api:

- ✅ Ya funcionan, puedes usarlos
- ⚠️ Tienen estructura vieja (pero funcional)
- Decisión de equipo si refactorizar o dejar como están

---

## 📊 Progreso Real

- **Funcionalidad**: 50% (3 de 5 microservicios funcionan)
- **Arquitectura Limpia**: 20% (solo 1 de 5 correctos)
- **Documentación**: 100% ✅

---

## 🆘 ¿Por dónde empiezo?

### Si eres nuevo:
1. Lee [`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md) (5 min)
2. Lee [`subscriptions-api/README.md`](subscriptions-api/README.md) (15 min)
3. Explora el código de `subscriptions-api/` (30 min)

### Si vas a implementar payments-api o search-api:
1. Lee [`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)
2. Sigue las instrucciones específicas de tu microservicio
3. Usa `subscriptions-api/` como referencia constante

---

## 🔗 Enlaces Rápidos

- [Resumen Honesto](RESUMEN_HONESTO.md) ⭐
- [Archivos a Refactorizar](ARCHIVOS_A_REFACTORIZAR.md) ⭐
- [Ejemplo de Referencia](subscriptions-api/README.md) ⭐
- [Estado Detallado](ESTADO_IMPLEMENTACION.md)
- [Guía General](LEEME_PRIMERO.md)
- [Arquitectura](ARQUITECTURA_MICROSERVICIOS.md)
- [Modelo de Datos](DIAGRAMA_ENTIDADES.md)

---

## 📝 Tecnologías

- **Backend**: Go 1.23
- **Bases de Datos**: MySQL 8.0, MongoDB 7.0
- **Mensajería**: RabbitMQ 3.12
- **Caché**: Memcached 1.6
- **Framework Web**: Gin
- **Contenedores**: Docker & Docker Compose

---

## ⚖️ Licencia

Proyecto académico - Universidad Católica de Córdoba

---

**Recuerda**: Solo `subscriptions-api` tiene la arquitectura correcta. Úsalo como referencia.
