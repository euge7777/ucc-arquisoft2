# Comandos de Testing Rápido

## 🚀 Inicio del Sistema

```bash
# Levantar todo
docker-compose -f docker-compose.new.yml up -d

# Esperar a que todo esté listo (30 segundos)
sleep 30

# Verificar que todo esté corriendo
docker-compose -f docker-compose.new.yml ps
```

## ✅ Health Checks

```bash
# Script para verificar todos los servicios
curl -s http://localhost:8080/healthz && echo " ✅ users-api"
curl -s http://localhost:8081/healthz && echo " ✅ subscriptions-api"
curl -s http://localhost:8082/healthz && echo " ✅ activities-api"
curl -s http://localhost:8083/healthz && echo " ✅ payments-api"
curl -s http://localhost:8084/healthz && echo " ✅ search-api"
```

## 🧪 Testing Completo - Flujo de Usuario

### 1. Registrar Usuario

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

**Salida esperada**:
```json
{
  "id_usuario": 1,
  "nombre": "Juan",
  "apellido": "Pérez",
  "username": "juanp",
  "email": "juan@example.com"
}
```

### 2. Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username_or_email": "juanp",
    "password": "Password123"
  }'
```

**Guardar el token JWT de la respuesta**

### 3. Crear Plan de Suscripción

```bash
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Premium",
    "descripcion": "Acceso completo a todas las actividades",
    "precio_mensual": 100.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true,
    "actividades_permitidas": []
  }'
```

**Guardar el ID del plan (ej: "507f1f77bcf86cd799439011")**

### 4. Crear Suscripción

```bash
# Reemplazar PLAN_ID con el ID del paso anterior
curl -X POST http://localhost:8081/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": "1",
    "plan_id": "PLAN_ID",
    "metodo_pago": "credit_card"
  }'
```

**Guardar el ID de la suscripción**

### 5. Crear Pago

```bash
# Reemplazar SUBSCRIPTION_ID
curl -X POST http://localhost:8083/payments \
  -H "Content-Type: application/json" \
  -d '{
    "entity_type": "subscription",
    "entity_id": "SUBSCRIPTION_ID",
    "user_id": "1",
    "amount": 100.00,
    "currency": "USD",
    "payment_method": "credit_card",
    "metadata": {
      "plan_nombre": "Plan Premium",
      "duracion_dias": 30
    }
  }'
```

**Guardar el ID del pago**

### 6. Procesar Pago

```bash
# Reemplazar PAYMENT_ID
curl -X POST http://localhost:8083/payments/PAYMENT_ID/process
```

### 7. Actualizar Estado de Suscripción

```bash
# Reemplazar SUBSCRIPTION_ID y PAYMENT_ID
curl -X PATCH http://localhost:8081/subscriptions/SUBSCRIPTION_ID/status \
  -H "Content-Type: application/json" \
  -d '{
    "estado": "activa",
    "pago_id": "PAYMENT_ID"
  }'
```

### 8. Crear Sucursal

```bash
curl -X POST http://localhost:8082/sucursales \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Sede Centro",
    "direccion": "Av. Principal 123",
    "telefono": "555-1234"
  }'
```

### 9. Crear Actividad

```bash
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
    "categoria": "Fitness",
    "requiere_plan_premium": false
  }'
```

**Guardar el ID de la actividad**

### 10. Inscribirse a Actividad

```bash
# Reemplazar ACTIVIDAD_ID
curl -X POST http://localhost:8082/inscripciones \
  -H "Content-Type: application/json" \
  -d '{
    "usuario_id": 1,
    "actividad_id": ACTIVIDAD_ID
  }'
```

### 11. Buscar Actividades

```bash
# Búsqueda simple
curl "http://localhost:8084/search?q=yoga&type=activity"

# Búsqueda avanzada
curl -X POST http://localhost:8084/search \
  -H "Content-Type: application/json" \
  -d '{
    "query": "yoga",
    "type": "activity",
    "filters": {
      "categoria": "Fitness",
      "dia": "Lunes"
    },
    "page": 1,
    "page_size": 10
  }'
```

### 12. Ver Estadísticas de Búsqueda

```bash
curl http://localhost:8084/search/stats
```

## 📊 Consultas Útiles

### Ver Usuarios

```bash
curl http://localhost:8080/users
```

### Ver Planes

```bash
curl http://localhost:8081/plans
```

### Ver Actividades

```bash
curl http://localhost:8082/actividades
```

### Ver Sucursales

```bash
curl http://localhost:8082/sucursales
```

### Ver Inscripciones de un Usuario

```bash
# Reemplazar USER_ID
curl http://localhost:8082/inscripciones/usuario/USER_ID
```

### Ver Suscripción Activa de Usuario

```bash
# Reemplazar USER_ID
curl http://localhost:8081/subscriptions/active/USER_ID
```

### Ver Pagos de un Usuario

```bash
# Reemplazar USER_ID
curl http://localhost:8083/payments/user/USER_ID
```

### Ver Pagos Pendientes

```bash
curl "http://localhost:8083/payments/status?status=pending"
```

## 🧹 Limpieza

```bash
# Detener todo
docker-compose -f docker-compose.new.yml down

# Detener y eliminar volúmenes (CUIDADO: borra todos los datos)
docker-compose -f docker-compose.new.yml down -v
```

## 🐛 Debugging

### Ver logs de un servicio específico

```bash
docker-compose -f docker-compose.new.yml logs -f users-api
docker-compose -f docker-compose.new.yml logs -f subscriptions-api
docker-compose -f docker-compose.new.yml logs -f activities-api
docker-compose -f docker-compose.new.yml logs -f payments-api
docker-compose -f docker-compose.new.yml logs -f search-api
```

### Ver logs de RabbitMQ

```bash
docker-compose -f docker-compose.new.yml logs -f rabbitmq
```

### Acceder a RabbitMQ Management UI

```
URL: http://localhost:15672
Usuario: admin
Password: admin
```

### Conectar a MySQL

```bash
docker exec -it gym-mysql mysql -uroot -proot123

# Dentro de MySQL
USE proyecto_integrador;
SHOW TABLES;
SELECT * FROM usuarios;
```

### Conectar a MongoDB

```bash
docker exec -it gym-mongo mongosh

# Dentro de MongoDB
use gym_subscriptions
db.planes.find()
db.suscripciones.find()

use payments
db.payments.find()
```

## 🔄 Testing de Eventos RabbitMQ

### Verificar que los eventos se publiquen

1. Acceder a RabbitMQ UI: http://localhost:15672
2. Ir a "Exchanges" → `gym_events`
3. Crear un plan o suscripción
4. Verificar en "Queues" → `search_indexer_queue` que llegó el mensaje

### Verificar que search-api consuma eventos

```bash
# Ver logs de search-api
docker-compose -f docker-compose.new.yml logs -f search-api

# Crear un plan
curl -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Test",
    "precio_mensual": 50.00,
    "tipo_acceso": "limitado",
    "duracion_dias": 30,
    "activo": true
  }'

# En los logs de search-api deberías ver:
# 📥 Evento recibido: plan.create (ID: ...)
# ✅ Documento indexado: plan_...
```

## 📈 Testing de Caché

### Primera búsqueda (MISS)

```bash
curl -i "http://localhost:8084/search?q=yoga"
# Ver header: X-Cache: MISS
```

### Segunda búsqueda (HIT)

```bash
curl -i "http://localhost:8084/search?q=yoga"
# Ver header: X-Cache: HIT
```

### Invalidar caché creando nueva actividad

```bash
# Crear actividad
curl -X POST http://localhost:8082/actividades \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Pilates",
    "descripcion": "Clase de pilates",
    "cupo": 15,
    "dia": "Martes",
    "horario_inicio": "10:00:00",
    "horario_final": "11:00:00",
    "sucursal_id": 1,
    "instructor": "Ana García",
    "categoria": "Fitness"
  }'

# La próxima búsqueda será MISS porque se invalidó el caché
curl -i "http://localhost:8084/search?q=yoga"
```

## 🎯 Script Completo de Testing

```bash
#!/bin/bash

echo "🚀 Iniciando sistema..."
docker-compose -f docker-compose.new.yml up -d

echo "⏳ Esperando 30 segundos para que todo esté listo..."
sleep 30

echo "✅ Verificando health checks..."
curl -s http://localhost:8080/healthz && echo " ✅ users-api OK"
curl -s http://localhost:8081/healthz && echo " ✅ subscriptions-api OK"
curl -s http://localhost:8082/healthz && echo " ✅ activities-api OK"
curl -s http://localhost:8083/healthz && echo " ✅ payments-api OK"
curl -s http://localhost:8084/healthz && echo " ✅ search-api OK"

echo "👤 Registrando usuario..."
curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Test",
    "apellido": "User",
    "username": "testuser",
    "email": "test@example.com",
    "password": "Test1234"
  }' | jq

echo "📋 Creando plan..."
curl -s -X POST http://localhost:8081/plans \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Plan Test",
    "precio_mensual": 100.00,
    "tipo_acceso": "completo",
    "duracion_dias": 30,
    "activo": true
  }' | jq

echo "🏢 Creando sucursal..."
curl -s -X POST http://localhost:8082/sucursales \
  -H "Content-Type: application/json" \
  -d '{
    "nombre": "Sede Test",
    "direccion": "Calle Test 123",
    "telefono": "555-0000"
  }' | jq

echo "🎯 Creando actividad..."
curl -s -X POST http://localhost:8082/actividades \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Yoga Test",
    "descripcion": "Clase de prueba",
    "cupo": 10,
    "dia": "Lunes",
    "horario_inicio": "08:00:00",
    "horario_final": "09:00:00",
    "sucursal_id": 1,
    "instructor": "Test Instructor",
    "categoria": "Fitness"
  }' | jq

echo "🔍 Buscando actividades..."
curl -s "http://localhost:8084/search?q=yoga" | jq

echo "✅ Testing completado!"
```

Guarda esto como `test.sh` y ejecútalo con:
```bash
chmod +x test.sh
./test.sh
```
