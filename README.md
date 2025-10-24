# Sistema de Gestión de Gimnasio - Arquitectura de Microservicios

Sistema de gestión de gimnasio implementado con **arquitectura de microservicios** en Go.

## 🚀 Inicio Rápido

### 1. Levantar Infraestructura

```bash
# Levantar bases de datos y servicios
docker-compose -f docker-compose.new.yml up -d mysql mongodb rabbitmq memcached
```

### 2. Ejecutar Microservicios

```bash
# users-api
cd users-api
go run cmd/api/main.go  # Puerto 8080

# subscriptions-api
cd subscriptions-api
go run cmd/api/main.go  # Puerto 8081

# activities-api
cd activities-api
go run cmd/api/main.go  # Puerto 8082

# payments-api
cd payments-api
go run cmd/api/main.go  # Puerto 8083

# search-api
cd search-api
go run cmd/api/main.go  # Puerto 8084
```

### 3. Verificar Health Checks

```bash
curl http://localhost:8080/healthz  # users-api
curl http://localhost:8081/healthz  # subscriptions-api
curl http://localhost:8082/healthz  # activities-api
curl http://localhost:8083/healthz  # payments-api
curl http://localhost:8084/healthz  # search-api
```

---

## 🏗️ Arquitectura

```
Frontend (React)
     │
     ├─→ users-api (8080)         MySQL
     ├─→ subscriptions-api (8081)  MongoDB + RabbitMQ
     ├─→ activities-api (8082)    MySQL + RabbitMQ
     ├─→ payments-api (8083)      MongoDB
     └─→ search-api (8084)        In-Memory + RabbitMQ + Memcached
```

### Microservicios

| Servicio | Puerto | Base de Datos | Estado | Descripción |
|----------|--------|---------------|--------|-------------|
| **users-api** | 8080 | MySQL | ✅ Funcional | Autenticación, JWT, CRUD usuarios |
| **subscriptions-api** | 8081 | MongoDB | ✅ Funcional | Planes y suscripciones + eventos |
| **activities-api** | 8082 | MySQL | ✅ Funcional | Actividades, sucursales, inscripciones |
| **payments-api** | 8083 | MongoDB | ✅ Funcional | Pagos genéricos, gateways múltiples |
| **search-api** | 8084 | In-Memory | ✅ Funcional | Búsqueda con caché de 2 niveles |

---

## 📁 Estructura del Proyecto

```
ucc-arquisoft2/
│
├── users-api/              # Autenticación y gestión de usuarios
├── subscriptions-api/      # Planes y suscripciones (⭐ Ejemplo de referencia)
├── activities-api/         # Actividades e inscripciones
├── payments-api/           # Sistema de pagos con múltiples gateways
├── search-api/             # Búsqueda y caché
├── frontend/               # Aplicación React
│
├── docker-compose.new.yml  # Infraestructura completa
│
├── ARQUITECTURA_MICROSERVICIOS.md  # Patrones y decisiones arquitectónicas
├── DIAGRAMA_ENTIDADES.md           # Modelo de datos completo
├── GUIA_IMPLEMENTAR_MICROSERVICIO.md
├── GUIA_COMPLETA_MICROSERVICIOS.md
└── INSTRUCCIONES_DOCKER.md
```

---

## 📚 Documentación

### Documentación General

- **[ARQUITECTURA_MICROSERVICIOS.md](ARQUITECTURA_MICROSERVICIOS.md)** - Patrones de diseño y decisiones arquitectónicas
- **[DIAGRAMA_ENTIDADES.md](DIAGRAMA_ENTIDADES.md)** - Modelo de datos completo con relaciones
- **[GUIA_IMPLEMENTAR_MICROSERVICIO.md](GUIA_IMPLEMENTAR_MICROSERVICIO.md)** - Guía para crear nuevos microservicios
- **[GUIA_COMPLETA_MICROSERVICIOS.md](GUIA_COMPLETA_MICROSERVICIOS.md)** - Guía de uso del sistema completo
- **[INSTRUCCIONES_DOCKER.md](INSTRUCCIONES_DOCKER.md)** - Instrucciones para Docker

### Documentación por Microservicio

Cada microservicio tiene su propio README con detalles específicos:

- [users-api/README.md](users-api/README.md) - API de usuarios y autenticación
- [subscriptions-api/README.md](subscriptions-api/README.md) - ⭐ **Ejemplo de referencia con arquitectura limpia**
- [activities-api/README.md](activities-api/README.md) - API de actividades
- [payments-api/README.md](payments-api/README.md) - API de pagos con gateways
  - [ARQUITECTURA_GATEWAYS_PAGOS.md](payments-api/ARQUITECTURA_GATEWAYS_PAGOS.md) - Arquitectura de gateways
  - [GUIA_IMPLEMENTACION_GATEWAYS.md](payments-api/GUIA_IMPLEMENTACION_GATEWAYS.md) - Guía de implementación
- [search-api/README.md](search-api/README.md) - API de búsqueda

---

## 🎯 Características Destacadas

### Patrones Implementados

- **Arquitectura Limpia** (Clean Architecture)
  - Separación de capas: Domain, Repository, Services, Controllers
  - Dependency Injection manual
  - DTOs separados de Entities

- **Event-Driven Architecture**
  - RabbitMQ para comunicación asíncrona
  - Eventos: subscription.created, inscription.created, etc.

- **Cache-Aside Pattern**
  - Caché de dos niveles (CCache local + Memcached distribuido)
  - TTL configurables

- **Repository Pattern**
  - Abstracción de acceso a datos
  - Interfaces + implementaciones (MongoDB, MySQL)

- **Gateway Pattern** (en payments-api)
  - Integración con múltiples pasarelas de pago
  - Strategy Pattern para intercambiar gateways
  - Factory Pattern para creación de instancias

### Seguridad

- **JWT Authentication** (users-api)
- **Password Hashing** (SHA-256)
- **Validación de Contraseñas Fuertes**
- **CORS Configurado**

### Observabilidad

- **Health Checks** en todos los servicios
- **Logs Estructurados**
- **Headers de Caché** (`X-Cache: HIT/MISS`)

---

## 🔄 Flujos de Datos

### Flujo 1: Crear Suscripción

```
1. Usuario → POST /subscriptions → subscriptions-api
2. subscriptions-api valida usuario con users-api (HTTP)
3. subscriptions-api crea suscripción con estado "pendiente_pago"
4. Publica evento a RabbitMQ: subscription.created
5. search-api consume evento y indexa
```

### Flujo 2: Crear Inscripción

```
1. Usuario → POST /inscripciones → activities-api
2. activities-api valida usuario y suscripción activa
3. activities-api crea inscripción
4. Publica evento a RabbitMQ: inscription.created
5. search-api actualiza cupo disponible
```

### Flujo 3: Búsqueda con Caché

```
1. Usuario → GET /search?q=yoga → search-api
2. Busca en CCache local (30s TTL)
   ├─ HIT → Return + Header "X-Cache: HIT"
   └─ MISS → Busca en Memcached (60s TTL)
       ├─ HIT → Guarda en CCache → Return
       └─ MISS → Ejecuta búsqueda → Guarda en ambos → Return
```

---

## 🛠️ Tecnologías

### Backend
- **Go 1.23** - Todos los microservicios
- **Gin** - Framework web HTTP

### Bases de Datos
- **MySQL 8.0** - users-api, activities-api
- **MongoDB 7.0** - subscriptions-api, payments-api

### Mensajería y Caché
- **RabbitMQ 3.12** - Comunicación asíncrona
- **Memcached 1.6** - Caché distribuido
- **CCache** - Caché local in-memory

### Infraestructura
- **Docker & Docker Compose**
- **Apache Solr 9** (opcional para search-api)

---

## 📊 Arquitectura Limpia (subscriptions-api)

**subscriptions-api es el ejemplo de referencia** que implementa correctamente todos los patrones:

```
subscriptions-api/
├── cmd/api/main.go                    # ✅ DI manual completa
├── internal/
│   ├── domain/
│   │   ├── entities/                  # ✅ Entidades de BD
│   │   └── dtos/                      # ✅ DTOs Request/Response
│   ├── repository/                    # ✅ Interfaces + MongoDB
│   ├── services/                      # ✅ Lógica de negocio con DI
│   ├── infrastructure/                # ✅ Servicios externos
│   ├── controllers/                   # ✅ Capa HTTP
│   ├── middleware/
│   ├── database/
│   └── config/
```

**Ver [subscriptions-api/README.md](subscriptions-api/README.md) para detalles completos.**

---

## 🧪 Testing Rápido

### Registrar Usuario

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan",
    "apellido": "Pérez",
    "username": "juanp",
    "email": "juan@example.com",
    "password": "Password123"
  }'
```

### Crear Plan

```bash
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Premium",
    "descripcion": "Acceso completo",
    "precio_mensual": 100.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true
  }'
```

### Buscar Actividades

```bash
curl "http://localhost:8084/search?q=yoga&type=activity"
```

---

## 🚧 Próximos Pasos

### Corto Plazo
- [ ] Implementar frontend completo (React)
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

---

## 🆘 Soporte

Para preguntas o problemas:
1. Revisar la documentación del microservicio específico
2. Consultar [ARQUITECTURA_MICROSERVICIOS.md](ARQUITECTURA_MICROSERVICIOS.md)
3. Verificar logs: `docker-compose logs <servicio>`

---

## 👥 Equipo

Proyecto desarrollado como parte de **Arquitectura de Software II** - Universidad Católica de Córdoba

---

## 📄 Licencia

Proyecto académico - Universidad Católica de Córdoba

---

**Última actualización**: 2025-01-15
