# GymPro - Sistema de Gestión de Gimnasio

GymPro es una aplicación web para la gestión de actividades y clases en un gimnasio. Permite a los usuarios ver, inscribirse y gestionar actividades deportivas, mientras que los administradores pueden gestionar las clases y usuarios.

## Características Principales

- 🏋️‍♂️ Catálogo de actividades deportivas
- 👤 Sistema de autenticación de usuarios
- 📝 Inscripción a actividades
- 👨‍💼 Panel de administración
- 📊 Control de cupos por actividad

## Tecnologías Utilizadas

### Frontend
- React
- React Router DOM
- CSS para estilos
- LocalStorage para gestión de sesión

### Backend
- Go (Golang)
- Gin (Framework web)
- GORM (ORM)
- MySQL
- JWT para autenticación

## Estructura del Proyecto

```
proyecto2025-morini-heredia/
├── frontend/           # Aplicación React
├── backend/           # Servidor Go
│   ├── controllers/   # Controladores de la API
│   ├── services/      # Lógica de negocio
│   ├── model/        # Modelos de datos
│   └── db/           # Configuración de base de datos
```

## Requisitos Previos

- Node.js y npm para el frontend
- Go 1.24.2 o superior
- MySQL

## Configuración

### Backend
1. Crear un archivo `.env` en la carpeta `backend` con las siguientes variables:
```
DB_USER=tu_usuario
DB_PASS=tu_contraseña
DB_HOST=localhost
DB_PORT=3306
DB_SCHEMA=proyecto_integrador
JWT_SECRET=tu_secreto
```

2. Instalar dependencias:
```bash
cd backend
go mod download
```

### Frontend
1. Instalar dependencias:
```bash
cd frontend
npm install
```

## Ejecución

### Backend
```bash
cd backend
go run main.go
```

### Frontend
```bash
cd frontend
npm run dev
```

## Características de Seguridad

- Contraseñas hasheadas con SHA-256
- Autenticación mediante JWT
- Middleware de protección para rutas
- Control de roles (admin/usuario)
- CORS configurado

## Funcionalidades por Rol

### Usuarios
- Ver catálogo de actividades
- Inscribirse/desinscribirse de actividades
- Ver sus inscripciones

### Administradores
- Gestionar actividades (crear, editar, eliminar)
- Ver todas las inscripciones
- Gestionar cupos

## Contribuidores
- Morini
- Heredia 