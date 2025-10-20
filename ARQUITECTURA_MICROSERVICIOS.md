# 🏗️ Arquitectura de Microservicios - Explicación Visual

## 🤔 Monolito vs Microservicios

### ❌ ANTES - Monolito (backend/)

```
┌─────────────────────────────────────────────────┐
│              backend/main.go                     │
│  ┌───────────────────────────────────────────┐  │
│  │  Todas las rutas en un solo proceso:     │  │
│  │  - /register                              │  │
│  │  - /login                                 │  │
│  │  - /actividades                           │  │
│  │  - /inscripciones                         │  │
│  │  - /suscripciones                         │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌───────────────────────────────────────────┐  │
│  │  middleware/ (compartido)                 │  │
│  │  - jwt.go                                 │  │
│  │  - cors.go                                │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  ┌───────────────────────────────────────────┐  │
│  │  Toda la lógica en un solo lugar         │  │
│  └───────────────────────────────────────────┘  │
│                                                  │
│  Un solo proceso → Si falla, TODO falla         │
└─────────────────────────────────────────────────┘
```

**Problemas:**
- ❌ Si falla una parte, falla TODO
- ❌ No puedes escalar solo una parte
- ❌ Difícil trabajar en equipo (conflictos)
- ❌ Despliegues arriesgados (todo o nada)

---

### ✅ AHORA - Microservicios

```
┌──────────────────┐   ┌────────────────────┐   ┌──────────────────┐
│  users-api       │   │  activities-api    │   │ subscriptions-   │
│  :8080           │   │  :8082             │   │ api :8081        │
│  ┌────────────┐  │   │  ┌──────────────┐  │   │  ┌────────────┐  │
│  │ main.go    │  │   │  │ main.go      │  │   │  │ main.go    │  │
│  │ Proceso 1  │  │   │  │ Proceso 2    │  │   │  │ Proceso 3  │  │
│  └────────────┘  │   │  └──────────────┘  │   │  └────────────┘  │
│                  │   │                    │   │                  │
│  internal/       │   │  internal/         │   │  internal/       │
│  ├─middleware/   │   │  ├─middleware/     │   │  ├─middleware/   │
│  ├─repository/   │   │  ├─repository/     │   │  ├─repository/   │
│  ├─services/     │   │  ├─services/       │   │  ├─services/     │
│  └─controllers/  │   │  └─controllers/    │   │  └─controllers/  │
│                  │   │                    │   │                  │
│  Rutas:          │   │  Rutas:            │   │  Rutas:          │
│  /register       │   │  /actividades      │   │  /planes         │
│  /login          │   │  /inscripciones    │   │  /suscripciones  │
│  /users          │   │  /sucursales       │   │                  │
│                  │   │                    │   │                  │
│  MySQL           │   │  MySQL             │   │  MongoDB         │
└──────────────────┘   └────────────────────┘   └──────────────────┘
      ▲                        ▲                        ▲
      │                        │                        │
      └────────────────────────┴────────────────────────┘
                    Comunicación HTTP/RabbitMQ
```

**Ventajas:**
- ✅ Si falla users-api, activities-api sigue funcionando
- ✅ Puedes escalar solo activities-api si tiene mucho tráfico
- ✅ Equipos diferentes trabajan sin conflictos
- ✅ Despliegues independientes (menos riesgo)

---

## 🔍 ¿Por qué middleware duplicado?

### ❌ Opción MALA - Middleware compartido

```
shared/
└── middleware/
    └── jwt.go  ← Todos dependen de esto

users-api/        ← Importa shared/middleware
activities-api/   ← Importa shared/middleware
subscriptions-api ← Importa shared/middleware

Problema: Si cambias jwt.go, debes:
1. Actualizar shared/
2. Rebuild users-api
3. Rebuild activities-api
4. Rebuild subscriptions-api
5. Desplegar LOS 3 servicios al mismo tiempo

❌ Perdiste la independencia de microservicios
```

### ✅ Opción BUENA - Middleware independiente

```
users-api/internal/middleware/jwt.go
activities-api/internal/middleware/jwt.go
subscriptions-api/internal/middleware/jwt.go

Ventaja: Si cambias el JWT de users-api:
1. Modificas users-api/internal/middleware/jwt.go
2. Rebuild SOLO users-api
3. Desplegar SOLO users-api
4. activities-api y subscriptions-api NO se afectan

✅ Mantienes la independencia de microservicios
```

**¿Código duplicado?**

Sí, pero es **intencional**. Es el precio de la independencia.

**Comparación:**

| Aspecto | Middleware Compartido | Middleware Duplicado |
|---------|----------------------|---------------------|
| Código | ✅ No duplicado | ❌ Duplicado |
| Independencia | ❌ Acoplado | ✅ Independiente |
| Despliegues | ❌ Deben ser juntos | ✅ Individuales |
| Riesgo | ❌ Alto (afecta a todos) | ✅ Bajo (afecta a uno) |
| Recomendado | ❌ NO para microservicios | ✅ SÍ para microservicios |

---

## 🎯 ¿Por qué múltiples main.go?

### Ejemplo Real:

```bash
# Terminal 1 - Iniciar users-api
$ cd users-api
$ go run cmd/api/main.go
→ Inicia proceso en puerto 8080
→ Carga middleware/jwt.go de users-api
→ Conecta a MySQL

# Terminal 2 - Iniciar activities-api
$ cd activities-api
$ go run cmd/api/main.go
→ Inicia OTRO proceso en puerto 8082
→ Carga middleware/jwt.go de activities-api
→ Conecta a MySQL (misma BD, diferentes tablas)

# Ahora tienes 2 procesos independientes corriendo
```

**¿Cómo se comunican?**

```
Usuario → Frontend → users-api:8080/login
                     ↓
                     Devuelve JWT token
                     ↓
Usuario → Frontend → activities-api:8082/inscripciones
                     ↓
                     activities-api valida JWT
                     ↓
                     activities-api llama HTTP GET a:
                     users-api:8080/users/:id
                     ↓
                     ¿Usuario existe? → Sí → Crea inscripción
```

---

## 📦 Estructura de Archivos - Explicación

### ✅ Estructura CORRECTA (actual)

```
ucc-arquisoft2/
│
├── docker-compose.new.yml    ← Orquesta TODOS los servicios
│
├── users-api/                ← Microservicio 1 (100% independiente)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go       ← Ejecutable 1 (puerto 8080)
│   │
│   ├── internal/             ← Código PRIVADO de users-api
│   │   ├── config/
│   │   ├── domain/
│   │   ├── dao/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── controllers/
│   │   └── middleware/       ← JWT propio
│   │       ├── cors.go
│   │       └── jwt.go
│   │
│   ├── go.mod                ← Dependencias propias
│   ├── Dockerfile            ← Build propio
│   └── .env.example
│
├── activities-api/           ← Microservicio 2 (100% independiente)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go       ← Ejecutable 2 (puerto 8082)
│   │
│   ├── internal/             ← Código PRIVADO de activities-api
│   │   ├── config/
│   │   ├── domain/
│   │   ├── dao/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── controllers/
│   │   └── middleware/       ← JWT propio (COPIA de users-api)
│   │       ├── cors.go       ← Código duplicado (intencional)
│   │       └── jwt.go        ← Código duplicado (intencional)
│   │
│   ├── go.mod                ← Dependencias propias
│   ├── Dockerfile            ← Build propio
│   └── .env.example
│
└── backend/                  ← LEGACY (monolito viejo)
    └── main.go               ← NO USAR (vieja arquitectura)
```

**Preguntas Frecuentes:**

**P: ¿Por qué `internal/`?**

R: Es una convención de Go. Todo lo que está en `internal/` NO puede ser importado por otros módulos. Esto FUERZA la independencia.

```go
// ❌ Esto NO compilará:
import "users-api/internal/middleware"  // Desde activities-api

// ✅ Esto SÍ funciona:
import "activities-api/internal/middleware"  // Desde activities-api
```

**P: ¿Por qué `cmd/api/`?**

R: Porque un proyecto puede tener múltiples ejecutables:

```
users-api/
├── cmd/
│   ├── api/main.go        ← Servidor HTTP
│   ├── migrate/main.go    ← Migración de BD (futuro)
│   └── seed/main.go       ← Seed de datos (futuro)
```

**P: ¿Por qué cada servicio tiene su `go.mod`?**

R: Porque son **módulos Go independientes**. Cada uno puede usar versiones diferentes de librerías.

---

## 🚀 Cómo Correr Todo

### Opción 1: Docker Compose (Recomendado)

```bash
# Un solo comando inicia TODO
docker-compose -f docker-compose.new.yml up --build

# Esto inicia:
# - MySQL (puerto 3306)
# - users-api (puerto 8080) ← main.go de users-api
# - activities-api (puerto 8082) ← main.go de activities-api
# - RabbitMQ, Solr, etc.
```

**¿Cómo funciona?**

```yaml
# docker-compose.new.yml

users-api:
  build: ./users-api      ← Ejecuta users-api/Dockerfile
  ports: ["8080:8080"]    ← Expone puerto 8080

activities-api:
  build: ./activities-api ← Ejecuta activities-api/Dockerfile
  ports: ["8082:8082"]    ← Expone puerto 8082
```

Cada `Dockerfile` hace:

```dockerfile
# users-api/Dockerfile
RUN go build -o users-api ./cmd/api  ← Compila cmd/api/main.go
CMD ["./users-api"]                  ← Ejecuta el binario
```

### Opción 2: Manual (Desarrollo)

```bash
# Terminal 1
cd users-api
go run cmd/api/main.go  ← Inicia proceso 1

# Terminal 2
cd activities-api
go run cmd/api/main.go  ← Inicia proceso 2
```

---

## 📊 Comunicación entre Microservicios

### Ejemplo: Inscribirse a una actividad

```
1. Usuario hace POST /inscripciones
   ↓
2. activities-api recibe la petición
   ↓
3. activities-api extrae user_id del JWT
   ↓
4. activities-api valida que usuario existe:
   HTTP GET → users-api:8080/users/:user_id
   ↓
5. activities-api valida que tiene suscripción activa:
   HTTP GET → subscriptions-api:8081/subscriptions/active/:user_id
   ↓
6. activities-api crea la inscripción en BD
   ↓
7. activities-api publica evento en RabbitMQ:
   "inscription.created"
   ↓
8. search-api escucha RabbitMQ y indexa en Solr
```

**Código real (en activities-api/internal/services/inscripciones.go):**

```go
// TODO: Validar que usuario existe
resp, err := http.Get(fmt.Sprintf("http://users-api:8080/users/%d", userID))
if resp.StatusCode == 404 {
    return errors.New("user not found")
}

// TODO: Validar suscripción activa
resp, err := http.Get(fmt.Sprintf("http://subscriptions-api:8081/subscriptions/active/%d", userID))
if resp.StatusCode == 404 {
    return errors.New("no active subscription")
}

// Crear inscripción
inscripcion, err := s.repository.Create(ctx, inscripcion)

// TODO: Publicar evento RabbitMQ
rabbitmq.Publish("inscription.created", inscripcion.ID)
```

---

## ✅ Verificación Final

### ¿Está bien la estructura?

✅ **SÍ** - Sigue las mejores prácticas de:
- Go Standard Project Layout
- Arquitectura de Microservicios
- Domain-Driven Design
- Dependency Injection

### ¿Está bien el middleware duplicado?

✅ **SÍ** - Es intencional para mantener independencia

### ¿Está bien tener múltiples main.go?

✅ **SÍ** - Cada microservicio es un ejecutable independiente

### ¿Cómo corro todo?

✅ **Docker Compose:**

```bash
docker-compose -f docker-compose.new.yml up --build
```

---

## 📚 Recursos de Aprendizaje

**Arquitectura de Microservicios:**
- [Martin Fowler - Microservices](https://martinfowler.com/articles/microservices.html)
- [12 Factor App](https://12factor.net/)

**Go Project Layout:**
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

**Comparaciones:**

| Concepto | Monolito | Microservicios |
|----------|----------|----------------|
| **Despliegue** | Todo junto | Independiente |
| **Escalamiento** | Escala todo | Escala por servicio |
| **Base de datos** | Una sola | Una por servicio (o compartida) |
| **Código compartido** | Compartido | Duplicado (intencional) |
| **Comunicación** | Interna (funciones) | HTTP/RabbitMQ |
| **Fallos** | Todo falla | Fallo aislado |
| **Complejidad** | Simple | Más complejo |
| **Equipo** | Pequeño | Múltiples equipos |

---

🎉 **Tu arquitectura está PERFECTAMENTE estructurada para microservicios!**
