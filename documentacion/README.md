# 📚 Documentación del Proyecto

Esta carpeta contiene la documentación general del sistema de microservicios.

## 📑 Documentos Disponibles

### [ARQUITECTURA_MICROSERVICIOS.md](ARQUITECTURA_MICROSERVICIOS.md)
Patrones de diseño y decisiones arquitectónicas del proyecto:
- Arquitectura de microservicios
- Patrones implementados (Repository, Strategy, Factory, Gateway)
- Comunicación entre servicios
- Principios SOLID aplicados

### [DIAGRAMA_ENTIDADES.md](DIAGRAMA_ENTIDADES.md)
Modelo de datos completo del sistema:
- Diagrama de entidades y relaciones
- Esquemas de bases de datos (MySQL y MongoDB)
- Descripción de cada entidad
- Relaciones entre microservicios

### [GUIA_IMPLEMENTAR_MICROSERVICIO.md](GUIA_IMPLEMENTAR_MICROSERVICIO.md)
Guía paso a paso para crear nuevos microservicios:
- Estructura de carpetas
- Implementación de capas (Domain, Repository, Services, Controllers)
- Dependency Injection manual
- Buenas prácticas

### [GUIA_COMPLETA_MICROSERVICIOS.md](GUIA_COMPLETA_MICROSERVICIOS.md)
Guía completa de uso del sistema:
- Instalación y configuración
- Endpoints de cada microservicio
- Ejemplos de uso con curl
- Flujos completos de usuario

### [INSTRUCCIONES_DOCKER.md](INSTRUCCIONES_DOCKER.md)
Instrucciones para trabajar con Docker:
- Configuración de contenedores
- Docker Compose
- Comandos útiles
- Troubleshooting

---

## 📖 Orden de Lectura Recomendado

### Para nuevos desarrolladores:
1. **ARQUITECTURA_MICROSERVICIOS.md** - Entender la arquitectura general
2. **DIAGRAMA_ENTIDADES.md** - Conocer el modelo de datos
3. **GUIA_COMPLETA_MICROSERVICIOS.md** - Aprender a usar el sistema
4. **INSTRUCCIONES_DOCKER.md** - Configurar el entorno

### Para implementar un nuevo microservicio:
1. **GUIA_IMPLEMENTAR_MICROSERVICIO.md** - Guía paso a paso
2. **ARQUITECTURA_MICROSERVICIOS.md** - Patrones a seguir
3. **subscriptions-api/README.md** - Ejemplo de referencia

---

## 🔗 Documentación Adicional

Cada microservicio tiene su propia documentación en su carpeta:
- [users-api/README.md](../users-api/README.md)
- [subscriptions-api/README.md](../subscriptions-api/README.md) ⭐ Ejemplo de referencia
- [activities-api/README.md](../activities-api/README.md)
- [payments-api/README.md](../payments-api/README.md)
  - [ARQUITECTURA_GATEWAYS_PAGOS.md](../payments-api/ARQUITECTURA_GATEWAYS_PAGOS.md)
  - [GUIA_IMPLEMENTACION_GATEWAYS.md](../payments-api/GUIA_IMPLEMENTACION_GATEWAYS.md)
- [search-api/README.md](../search-api/README.md)

---

**Última actualización**: 2025-01-15
