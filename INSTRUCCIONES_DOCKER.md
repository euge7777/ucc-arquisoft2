# 🐳 Cómo Correr Todo el Proyecto con Docker Compose

## 📋 Respuestas a tus Preguntas

### 1. **¿Por qué cada microservicio tiene su propio middleware?**

✅ **SÍ, está CORRECTO.**

En arquitectura de microservicios, cada servicio debe ser **independiente y autónomo**:

```
users-api/internal/middleware/     ✅ CORRECTO - Independiente
activities-api/internal/middleware/ ✅ CORRECTO - Independiente
```

**¿Por qué NO compartir middleware?**

❌ **Acoplamiento:**
```
shared/middleware/jwt.go  ← Si cambias esto, afectas a TODOS los servicios
```

✅ **Independencia:**
```
users-api/internal/middleware/jwt.go      ← Servicio 1
activities-api/internal/middleware/jwt.go ← Servicio 2
```

**Ventajas:**
- Cada servicio se despliega independientemente
- Puedes usar versiones diferentes de Go en cada servicio
- Si un servicio falla, NO afecta a los demás
- Equipos diferentes pueden trabajar sin interferir

**Comparación:**

| Monolito (Anterior) | Microservicios (Actual) |
|---------------------|-------------------------|
| `backend/middleware/` (compartido) | `users-api/internal/middleware/` |
| Un solo proceso | Múltiples procesos independientes |
| Un solo main.go | Un main.go por servicio |
| Si falla, todo falla | Si falla uno, los demás siguen |

---

### 2. **¿Por qué tantos main.go?**

✅ **SÍ, cada microservicio tiene su propio main.go.**

Cada `main.go` inicia un **proceso independiente**:

```
users-api/cmd/api/main.go         → Proceso 1 (puerto 8080)
activities-api/cmd/api/main.go    → Proceso 2 (puerto 8082)
subscriptions-api/cmd/api/main.go → Proceso 3 (puerto 8081)
```

**¿Cómo se comunican?**

```
┌─────────────┐      HTTP/RabbitMQ     ┌──────────────────┐
│ users-api   │ ←──────────────────────→ │ activities-api   │
│ :8080       │                          │ :8082            │
│ (Proceso 1) │                          │ (Proceso 2)      │
└─────────────┘                          └──────────────────┘
```

**Ejemplo:**
1. Usuario se registra → `users-api:8080` crea usuario y devuelve JWT
2. Usuario se inscribe a clase → `activities-api:8082` valida JWT y crea inscripción
3. Activities-api valida que usuario existe → HTTP GET a `users-api:8080/users/:id`

---

### 3. **¿Cómo correr TODO de una vez?**

✅ **Con Docker Compose.**

---

## 🚀 Instrucciones para Correr Todo

### Opción 1: Con Docker Compose (Recomendado)

#### Paso 1: Verificar que Docker esté corriendo

```bash
docker --version
docker-compose --version
```

#### Paso 2: Correr todos los servicios

```bash
# Desde la raíz del proyecto
cd C:\Users\eli_v\ucc-arquisoft2

# Iniciar TODOS los servicios
docker-compose -f docker-compose.new.yml up --build
```

**Esto iniciará:**
- ✅ MySQL (puerto 3306)
- ✅ MongoDB (puerto 27017)
- ✅ RabbitMQ (puerto 5672, UI en 15672)
- ✅ Memcached (puerto 11211)
- ✅ Solr (puerto 8983)
- ✅ **users-api** (puerto 8080)
- ✅ **activities-api** (puerto 8082)

#### Paso 3: Verificar que todo funciona

**Abrir otra terminal y probar:**

```bash
# Health checks
curl http://localhost:8080/healthz
curl http://localhost:8082/healthz

# Users API
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Juan",
    "apellido": "Pérez",
    "username": "jperez",
    "email": "juan@example.com",
    "password": "Password123"
  }'

# Activities API
curl http://localhost:8082/actividades
```

#### Paso 4: Ver logs

```bash
# Ver logs de todos los servicios
docker-compose -f docker-compose.new.yml logs -f

# Ver logs de un servicio específico
docker-compose -f docker-compose.new.yml logs -f users-api
docker-compose -f docker-compose.new.yml logs -f activities-api
```

#### Paso 5: Detener todo

```bash
# Detener sin borrar datos
docker-compose -f docker-compose.new.yml down

# Detener Y borrar datos (⚠️ borra la BD)
docker-compose -f docker-compose.new.yml down -v
```

---

### Opción 2: Correr servicios individuales (Desarrollo)

#### Terminal 1: MySQL

```bash
docker run --name gym-mysql \
  -e MYSQL_ROOT_PASSWORD=root123 \
  -e MYSQL_DATABASE=proyecto_integrador \
  -p 3306:3306 \
  -d mysql:8.0
```

#### Terminal 2: users-api

```bash
cd users-api
go run cmd/api/main.go
```

#### Terminal 3: activities-api

```bash
cd activities-api
go run cmd/api/main.go
```

**Pros:** Más rápido para desarrollo (no rebuilds)
**Contras:** Debes iniciar cada servicio manualmente

---

## 🌐 Puertos y URLs

Una vez que todo esté corriendo con Docker Compose:

| Servicio | Puerto | URL |
|----------|--------|-----|
| **users-api** | 8080 | http://localhost:8080 |
| **activities-api** | 8082 | http://localhost:8082 |
| **MySQL** | 3306 | localhost:3306 |
| **MongoDB** | 27017 | localhost:27017 |
| **RabbitMQ UI** | 15672 | http://localhost:15672 (admin/admin) |
| **Solr UI** | 8983 | http://localhost:8983 |

---

## 📊 Arquitectura Visual

```
┌──────────────────────────────────────────────────────────┐
│                     Docker Network                        │
│                                                           │
│  ┌──────────┐   ┌──────────┐   ┌──────────┐            │
│  │ MySQL    │   │ MongoDB  │   │ RabbitMQ │            │
│  │ :3306    │   │ :27017   │   │ :5672    │            │
│  └────┬─────┘   └────┬─────┘   └────┬─────┘            │
│       │              │              │                    │
│  ┌────▼─────────────────────────────▼─────┐             │
│  │                                         │             │
│  │  ┌─────────────┐  ┌─────────────────┐  │             │
│  │  │ users-api   │  │ activities-api  │  │             │
│  │  │ :8080       │←→│ :8082           │  │             │
│  │  └──────┬──────┘  └──────┬──────────┘  │             │
│  │         │                │              │             │
│  └─────────┼────────────────┼──────────────┘             │
│            │                │                            │
└────────────┼────────────────┼────────────────────────────┘
             │                │
             ▼                ▼
        localhost:8080   localhost:8082
```

---

## ✅ Checklist de Verificación

Después de correr `docker-compose up`, verifica:

- [ ] MySQL está corriendo: `docker ps | grep gym-mysql`
- [ ] users-api está corriendo: `curl http://localhost:8080/healthz`
- [ ] activities-api está corriendo: `curl http://localhost:8082/healthz`
- [ ] Puedes registrar un usuario
- [ ] Puedes hacer login
- [ ] Puedes listar actividades
- [ ] Logs no muestran errores: `docker-compose -f docker-compose.new.yml logs`

---

## 🐛 Troubleshooting

### Problema: "Error connecting to MySQL"

**Solución:**

```bash
# Ver logs de MySQL
docker-compose -f docker-compose.new.yml logs mysql

# Esperar a que MySQL esté ready
# users-api y activities-api esperan automáticamente (healthcheck)
```

### Problema: "Port already in use"

**Solución:**

```bash
# Ver qué está usando el puerto
netstat -ano | findstr :8080
netstat -ano | findstr :3306

# Opción 1: Matar el proceso
taskkill /PID <PID_NUMBER> /F

# Opción 2: Cambiar puerto en docker-compose.new.yml
# Por ejemplo, cambiar "8080:8080" a "8081:8080"
```

### Problema: "Cannot compile Go code"

**Solución:**

```bash
# Rebuild forzando
docker-compose -f docker-compose.new.yml up --build --force-recreate

# Ver logs del build
docker-compose -f docker-compose.new.yml build users-api
docker-compose -f docker-compose.new.yml build activities-api
```

---

## 📂 Estructura de Carpetas (Microservicios)

```
ucc-arquisoft2/
│
├── docker-compose.new.yml    ← Orquesta TODOS los servicios
│
├── users-api/                ← Microservicio 1 (INDEPENDIENTE)
│   ├── cmd/api/main.go       ← Proceso independiente
│   ├── internal/
│   │   ├── middleware/       ← Middleware propio
│   │   ├── repository/
│   │   ├── services/
│   │   └── controllers/
│   ├── Dockerfile            ← Build propio
│   └── go.mod                ← Dependencias propias
│
├── activities-api/           ← Microservicio 2 (INDEPENDIENTE)
│   ├── cmd/api/main.go       ← Proceso independiente
│   ├── internal/
│   │   ├── middleware/       ← Middleware propio (copiado, no compartido)
│   │   ├── repository/
│   │   ├── services/
│   │   └── controllers/
│   ├── Dockerfile            ← Build propio
│   └── go.mod                ← Dependencias propias
│
└── subscriptions-api/        ← Microservicio 3 (TODO: crear)
    └── ...
```

**IMPORTANTE:**

✅ **Cada servicio tiene:**
- Su propio `main.go`
- Su propio `go.mod`
- Su propio `Dockerfile`
- Su propia carpeta `internal/`
- Su propio middleware (aunque sea código duplicado)

❌ **NO hay carpeta `shared/`:**
- No hay código compartido entre servicios
- Cada servicio es 100% independiente

---

## 🎯 Comandos Útiles

```bash
# Ver servicios corriendo
docker-compose -f docker-compose.new.yml ps

# Ver logs de todos
docker-compose -f docker-compose.new.yml logs -f

# Ver logs de uno solo
docker-compose -f docker-compose.new.yml logs -f users-api

# Rebuild solo un servicio
docker-compose -f docker-compose.new.yml up --build users-api

# Entrar a un contenedor
docker exec -it gym-users-api sh
docker exec -it gym-mysql mysql -uroot -proot123

# Reiniciar un servicio
docker-compose -f docker-compose.new.yml restart users-api

# Ver uso de recursos
docker stats

# Limpiar todo
docker-compose -f docker-compose.new.yml down -v
docker system prune -a
```

---

## 📝 Resumen Final

### ¿Es correcto tener middleware duplicado?
✅ **SÍ** - Cada microservicio es independiente

### ¿Es correcto tener múltiples main.go?
✅ **SÍ** - Cada microservicio es un ejecutable independiente

### ¿Cómo corro todo de una vez?
✅ **Docker Compose:**
```bash
docker-compose -f docker-compose.new.yml up --build
```

### ¿La distribución de carpetas está bien?
✅ **SÍ** - Sigue las mejores prácticas de microservicios

---

🎉 **¡Todo está correctamente estructurado para microservicios!**
