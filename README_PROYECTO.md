# Sistema de Gestión de Gimnasio - Arquitectura de Microservicios

## 🎯 Resumen del Proyecto

Sistema completo de gestión de gimnasio implementado con **arquitectura de microservicios** en Go, con 5 servicios independientes, comunicación asíncrona vía RabbitMQ, y sistema de caché de dos niveles.

## 📊 Arquitectura Completa

```
                    ┌─────────────────────┐
                    │     FRONTEND        │
                    │   (React/Next.js)   │
                    └──────────┬──────────┘
                               │ HTTP/JSON
        ┌──────────────────────┼──────────────────────┐
        │                      │                      │
        ▼                      ▼                      ▼
   ┌─────────┐          ┌──────────┐          ┌──────────┐
   │ users   │          │subscrip- │          │activities│
   │  api    │◄────────►│ tions    │◄────────►│   api    │
   │ :8080   │  HTTP    │  api     │  HTTP    │  :8082   │
   │ MySQL   │          │ :8081    │          │  MySQL   │
   └─────────┘          │ MongoDB  │          └──────────┘
                        └─────┬────┘                │
                              │                     │
                              ▼                     │
                        ┌─────────┐                │
                        │payments │                │
                        │  api    │                │
                        │ :8083   │                │
                        │MongoDB  │                │
                        └─────────┘                │
                                                   │
        ┌──────────────────────────────────────────┘
        │ RabbitMQ Events
        ▼
   ┌─────────┐          ┌──────────┐
   │ search  │◄────────►│Memcached │
   │  api    │  Cache   │ :11211   │
   │ :8084   │          └──────────┘
   └─────────┘
        ▲
        │ Consume
   ┌────▼────┐
   │RabbitMQ │
   │ :5672   │
   └─────────┘
```

## 🚀 Microservicios Implementados

| Servicio | Puerto | Base de Datos | Estado | Descripción |
|----------|--------|---------------|--------|-------------|
| **users-api** | 8080 | MySQL | ✅ Completo | Autenticación, JWT, CRUD usuarios |
| **subscriptions-api** | 8081 | MongoDB | ✅ Nuevo | Planes y suscripciones + RabbitMQ |
| **activities-api** | 8082 | MySQL | ✅ Migrado | Actividades, sucursales, inscripciones |
| **payments-api** | 8083 | MongoDB | ✅ Nuevo | API genérica de pagos (reutilizable) |
| **search-api** | 8084 | In-Memory* | ✅ Nuevo | Búsqueda + Caché 2 niveles + RabbitMQ consumer |

\* Migrable a Apache Solr en producción

## 📁 Estructura del Proyecto

```
ucc-arquisoft2/
│
├── users-api/                   ✅ Autenticación y usuarios
│   ├── cmd/api/
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
├── subscriptions-api/           ✅ Planes y suscripciones
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── services/
│   │   └── clients/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── activities-api/              ✅ Actividades e inscripciones
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── repository/
│   │   └── services/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── payments-api/                ✅ API genérica de pagos
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── database/
│   │   ├── handlers/
│   │   ├── models/
│   │   └── services/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── search-api/                  ✅ Búsqueda con caché
│   ├── cmd/api/
│   ├── internal/
│   │   ├── config/
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── services/
│   │   └── consumers/
│   ├── Dockerfile
│   ├── go.mod
│   └── README.md
│
├── docker-compose.new.yml       ✅ Infraestructura completa
│
├── DIAGRAMA_ENTIDADES.md        📚 Modelo de datos
├── ARQUITECTURA_MICROSERVICIOS.md  📚 Patrones y decisiones
└── GUIA_COMPLETA_MICROSERVICIOS.md 📚 Guía de implementación
```

## 🔧 Tecnologías Utilizadas

### Backend
- **Go 1.23** - Todos los microservicios
- **Gin** - Framework web HTTP

### Bases de Datos
- **MySQL 8.0** - users-api, activities-api
- **MongoDB 7.0** - subscriptions-api, payments-api

### Mensajería y Caché
- **RabbitMQ 3.12** - Comunicación asíncrona entre microservicios
- **Memcached 1.6** - Caché distribuido
- **CCache (in-memory)** - Caché local

### Infraestructura
- **Docker & Docker Compose** - Contenedores y orquestación
- **Apache Solr 9** (opcional) - Motor de búsqueda

## 🏃 Inicio Rápido

### Opción 1: Docker Compose (Recomendado)

```bash
# Levantar todos los servicios
docker-compose -f docker-compose.new.yml up -d

# Ver logs
docker-compose -f docker-compose.new.yml logs -f

# Ver estado
docker-compose -f docker-compose.new.yml ps

# Detener todo
docker-compose -f docker-compose.new.yml down
```

### Opción 2: Ejecución Local

Necesitas tener instalado:
- Go 1.23+
- MySQL 8.0
- MongoDB 7.0
- RabbitMQ 3.12
- Memcached 1.6

```bash
# Terminal 1: users-api
cd users-api
go run cmd/api/main.go

# Terminal 2: subscriptions-api
cd subscriptions-api
go run cmd/api/main.go

# Terminal 3: activities-api
cd activities-api
go run cmd/api/main.go

# Terminal 4: payments-api
cd payments-api
go run cmd/api/main.go

# Terminal 5: search-api
cd search-api
go run cmd/api/main.go
```

### Health Checks

```bash
curl http://localhost:8080/healthz  # users-api
curl http://localhost:8081/healthz  # subscriptions-api
curl http://localhost:8082/healthz  # activities-api
curl http://localhost:8083/healthz  # payments-api
curl http://localhost:8084/healthz  # search-api
```

## 📚 Documentación

### Guías Principales

1. **[DIAGRAMA_ENTIDADES.md](DIAGRAMA_ENTIDADES.md)**
   - Modelo de datos completo
   - Relaciones entre entidades
   - Esquemas de bases de datos

2. **[ARQUITECTURA_MICROSERVICIOS.md](ARQUITECTURA_MICROSERVICIOS.md)**
   - Patrones de diseño
   - Decisiones arquitectónicas
   - Comunicación entre servicios

3. **[GUIA_COMPLETA_MICROSERVICIOS.md](GUIA_COMPLETA_MICROSERVICIOS.md)**
   - Implementación paso a paso
   - Ejemplos de uso (curl)
   - Troubleshooting

### READMEs por Microservicio

- `users-api/README.md` - API de usuarios
- `subscriptions-api/README.md` - API de suscripciones
- `activities-api/README.md` - API de actividades
- `payments-api/README.md` - API de pagos
- `search-api/README.md` - API de búsqueda

## 🔄 Flujos de Datos

### Flujo 1: Crear Suscripción

```
1. Usuario → POST /subscriptions → subscriptions-api
2. subscriptions-api valida usuario con users-api (HTTP)
3. subscriptions-api valida plan existe y está activo
4. Crea suscripción con estado "pendiente_pago"
5. Publica evento a RabbitMQ: subscription.create
6. search-api consume evento y indexa
```

### Flujo 2: Búsqueda con Caché

```
1. Usuario → GET /search?q=yoga → search-api
2. Busca en CCache local (30s TTL)
   ├─ HIT → Return + Header "X-Cache: HIT"
   └─ MISS → 3
3. Busca en Memcached (60s TTL)
   ├─ HIT → Guarda en CCache → Return
   └─ MISS → 4
4. Ejecuta búsqueda real (in-memory/Solr)
5. Guarda en Memcached + CCache
6. Return + Header "X-Cache: MISS"
```

### Flujo 3: Crear Inscripción

```
1. Usuario → POST /inscripciones → activities-api
2. activities-api valida usuario con users-api
3. activities-api valida suscripción activa con subscriptions-api
4. Valida cupo disponible
5. Crea inscripción
6. Publica evento a RabbitMQ: inscription.create
7. search-api actualiza cupo disponible en índice
```

## 🎨 Características Destacadas

### ✅ Patrones Implementados

- **Arquitectura Limpia** (handlers, services, repositories)
- **Event-Driven Architecture** (RabbitMQ)
- **API Gateway Pattern** (cada microservicio expone su API)
- **Cache-Aside Pattern** (caché de dos niveles)
- **Circuit Breaker** (manejo de errores en llamadas HTTP)

### ✅ Seguridad

- **JWT Authentication** (users-api)
- **Password Hashing** (SHA-256)
- **Validación de Contraseñas Fuertes**
- **CORS Configurado**

### ✅ Observabilidad

- **Health Checks** en todos los servicios
- **Logs Estructurados**
- **Headers de Caché** (`X-Cache: HIT/MISS`)

### ✅ Escalabilidad

- **Stateless Services** (pueden escalar horizontalmente)
- **Base de Datos por Servicio**
- **Caché Distribuido** (Memcached)
- **Message Queue** (RabbitMQ para desacoplamiento)

## 🚧 Próximos Pasos

### Corto Plazo

- [ ] Implementar frontend (React/Next.js)
- [ ] Agregar tests unitarios y de integración
- [ ] Migrar search-api a Apache Solr
- [ ] Implementar métricas (Prometheus + Grafana)

### Mediano Plazo

- [ ] API Gateway (Kong/Traefik)
- [ ] Service Discovery (Consul)
- [ ] Distributed Tracing (Jaeger)
- [ ] Autenticación OAuth2

### Largo Plazo

- [ ] Migrar a Kubernetes
- [ ] CI/CD completo (GitHub Actions)
- [ ] Monitoreo avanzado (ELK Stack)

## 👥 Equipo y Contribuciones

Este proyecto fue desarrollado como parte de la materia **Arquitectura de Software II** en la Universidad Católica de Córdoba.

### Cómo Contribuir

1. Cada equipo debe implementar su microservicio asignado
2. Seguir la estructura de carpetas definida
3. Documentar en el README del microservicio
4. Probar integraciones con otros servicios
5. Actualizar docker-compose si es necesario

## 📞 Soporte

Para preguntas o problemas:
1. Revisar `GUIA_COMPLETA_MICROSERVICIOS.md` (sección Troubleshooting)
2. Verificar logs: `docker-compose logs <servicio>`
3. Consultar README del microservicio específico

## 📄 Licencia

Este proyecto es de uso académico para la Universidad Católica de Córdoba.

---

## 🚨 IMPORTANTE: Estado REAL del Proyecto

**⚠️ LEE ESTO PRIMERO**: **[`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)** ⭐

### Estado Real:

| Microservicio | ¿Funciona? | Arquitectura Limpia | Archivos Viejos |
|---------------|------------|---------------------|-----------------|
| users-api | ✅ SÍ | ❌ NO | ⚠️ handlers/, models/ |
| subscriptions-api | ✅ SÍ | ✅ **SÍ** | ✅ Ninguno |
| activities-api | ✅ SÍ | ❌ NO | ⚠️ handlers/, models/ |
| payments-api | ❌ NO | ❌ NO | ⚠️ handlers/, models/, services/ |
| search-api | ❌ NO | ❌ NO | ⚠️ handlers/, models/, clients/ |

**Solo 1 de 5 tiene arquitectura correcta**: `subscriptions-api`

**Progreso real**: 50% funcionalidad, 20% arquitectura limpia

---

### 📖 Documentos ESENCIALES (leer en orden):

1. **[`RESUMEN_HONESTO.md`](RESUMEN_HONESTO.md)** ⭐⭐⭐ - **EMPIEZA AQUÍ**
2. **[`ARCHIVOS_A_REFACTORIZAR.md`](ARCHIVOS_A_REFACTORIZAR.md)** ⭐⭐ - Qué eliminar
3. **[`subscriptions-api/README.md`](subscriptions-api/README.md)** ⭐⭐ - Ejemplo de referencia
4. [`ESTADO_IMPLEMENTACION.md`](ESTADO_IMPLEMENTACION.md) - Estado detallado
5. [`LEEME_PRIMERO.md`](LEEME_PRIMERO.md) - Conceptos generales

---

**🎯 Microservicio de Referencia**: `subscriptions-api/` es el ÚNICO con arquitectura correcta

---

## 📚 Documentación Esencial (En Orden de Lectura)

1. **`LEEME_PRIMERO.md`** ⭐ - **Empieza aquí**
2. **`ESTADO_IMPLEMENTACION.md`** ⭐ - Estado detallado de cada microservicio
3. **`subscriptions-api/README.md`** ⭐ - Arquitectura limpia explicada
4. `ARQUITECTURA_MICROSERVICIOS.md` - Patrones generales
5. `DIAGRAMA_ENTIDADES.md` - Modelo de datos
6. `GUIA_COMPLETA_MICROSERVICIOS.md` - Guía de uso
7. `TEST_COMMANDS.md` - Comandos de testing

---
