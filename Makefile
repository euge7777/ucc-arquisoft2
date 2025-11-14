.PHONY: help docker-up docker-down docker-logs local-up local-down install-deps clean healthz

# Variables
COMPOSE_FILE := docker-compose.new.yml
DOCKER_COMPOSE := docker compose -f $(COMPOSE_FILE)

# Colores
GREEN := \033[0;32m
YELLOW := \033[0;33m
CYAN := \033[0;36m
RED := \033[0;31m
NC := \033[0m # No Color

help:
	@echo "$(CYAN)╔════════════════════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)║          COMANDOS DISPONIBLES - MICROSERVICIOS GYM             ║$(NC)"
	@echo "$(CYAN)╚════════════════════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(GREEN)MODO DOCKER (TODO en contenedores):$(NC)"
	@echo "  make docker-up          → Levantar todo (infraestructura + servicios)"
	@echo "  make docker-down        → Detener todo"
	@echo "  make docker-logs        → Ver logs en tiempo real"
	@echo "  make docker-restart     → Reiniciar todos los servicios"
	@echo ""
	@echo "$(GREEN)MODO LOCAL (Infraestructura en Docker, servicios en tu máquina):$(NC)"
	@echo "  make local-up           → Levantar infraestructura"
	@echo "  make local-down         → Detener infraestructura"
	@echo "  make local-start-users  → Iniciar users-api"
	@echo "  make local-start-subs   → Iniciar subscriptions-api"
	@echo "  make local-start-acts   → Iniciar activities-api"
	@echo "  make local-start-pays   → Iniciar payments-api"
	@echo "  make local-start-search → Iniciar search-api"
	@echo ""
	@echo "$(GREEN)HERRAMIENTAS:$(NC)"
	@echo "  make install-deps       → Descargar dependencias de Go en todos los servicios"
	@echo "  make healthz            → Verificar salud de todos los servicios"
	@echo "  make db-mysql           → Conectarse a MySQL"
	@echo "  make db-mongo           → Conectarse a MongoDB"
	@echo "  make clean              → Limpiar contenedores y volúmenes"
	@echo ""
	@echo "$(CYAN)PUERTOS:$(NC)"
	@echo "  • Users API:         http://localhost:8080"
	@echo "  • Subscriptions API: http://localhost:8081"
	@echo "  • Activities API:    http://localhost:8082"
	@echo "  • Payments API:      http://localhost:8083"
	@echo "  • Search API:        http://localhost:8084"
	@echo "  • Frontend:          http://localhost:5173"
	@echo "  • RabbitMQ Admin:    http://localhost:15672 (admin/admin)"
	@echo "  • Solr Admin:        http://localhost:8983/solr"
	@echo ""

# Docker Mode
docker-up:
	@echo "$(CYAN)Levantando infraestructura + microservicios con Docker...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✓ Servicios levantados$(NC)"
	@make healthz

docker-down:
	@echo "$(CYAN)Deteniendo servicios...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✓ Servicios detenidos$(NC)"

docker-restart:
	@echo "$(CYAN)Reiniciando servicios...$(NC)"
	$(DOCKER_COMPOSE) restart
	@echo "$(GREEN)✓ Servicios reiniciados$(NC)"

docker-logs:
	$(DOCKER_COMPOSE) logs -f

docker-logs-users:
	$(DOCKER_COMPOSE) logs -f users-api

docker-logs-subs:
	$(DOCKER_COMPOSE) logs -f subscriptions-api

docker-logs-acts:
	$(DOCKER_COMPOSE) logs -f activities-api

docker-logs-pays:
	$(DOCKER_COMPOSE) logs -f payments-api

docker-logs-search:
	$(DOCKER_COMPOSE) logs -f search-api

# Local Mode
local-up:
	@echo "$(CYAN)Levantando infraestructura (MySQL, MongoDB, RabbitMQ, etc)...$(NC)"
	$(DOCKER_COMPOSE) up -d mysql mongo rabbitmq memcached solr
	@echo "$(GREEN)✓ Infraestructura lista$(NC)"
	@echo "$(YELLOW)Los microservicios deben iniciarse en terminales separadas:$(NC)"
	@echo "  make local-start-users"
	@echo "  make local-start-subs"
	@echo "  make local-start-acts"
	@echo "  make local-start-pays"
	@echo "  make local-start-search"

local-down:
	@echo "$(CYAN)Deteniendo infraestructura...$(NC)"
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✓ Infraestructura detenida$(NC)"

local-start-users:
	@echo "$(CYAN)Iniciando users-api...$(NC)"
	@cd users-api && \
	[ -f .env ] || cp .env.example .env && \
	go mod download && \
	go run cmd/api/main.go

local-start-subs:
	@echo "$(CYAN)Iniciando subscriptions-api...$(NC)"
	@cd subscriptions-api && \
	[ -f .env ] || cp .env.example .env && \
	go mod download && \
	go run cmd/api/main.go

local-start-acts:
	@echo "$(CYAN)Iniciando activities-api...$(NC)"
	@cd activities-api && \
	[ -f .env ] || cp .env.example .env && \
	go mod download && \
	go run cmd/api/main.go

local-start-pays:
	@echo "$(CYAN)Iniciando payments-api...$(NC)"
	@cd payments-api && \
	[ -f .env ] || cp .env.example .env && \
	go mod download && \
	go run cmd/api/main.go

local-start-search:
	@echo "$(CYAN)Iniciando search-api...$(NC)"
	@cd search-api && \
	[ -f .env ] || cp .env.example .env && \
	go mod download && \
	go run cmd/api/main.go

# Tools
install-deps:
	@echo "$(CYAN)Descargando dependencias...$(NC)"
	@for dir in users-api subscriptions-api activities-api payments-api search-api; do \
		echo "  • $$dir..."; \
		cd $$dir && go mod download && cd ..; \
	done
	@echo "$(GREEN)✓ Dependencias descargadas$(NC)"

healthz:
	@echo "$(CYAN)Verificando salud de servicios...$(NC)"
	@echo ""
	@for port in 8080 8081 8082 8083 8084; do \
		service_name=$$([ $$port -eq 8080 ] && echo "Users API" || [ $$port -eq 8081 ] && echo "Subscriptions API" || [ $$port -eq 8082 ] && echo "Activities API" || [ $$port -eq 8083 ] && echo "Payments API" || echo "Search API"); \
		if curl -s http://localhost:$$port/healthz > /dev/null 2>&1; then \
			echo "$(GREEN)✓$$service_name$$NC (puerto $$port)"; \
		else \
			echo "$(RED)✗$$service_name$$NC (puerto $$port)"; \
		fi; \
	done
	@echo ""

db-mysql:
	@echo "$(CYAN)Conectando a MySQL...$(NC)"
	docker exec -it mysql mysql -uroot -proot123

db-mongo:
	@echo "$(CYAN)Conectando a MongoDB...$(NC)"
	docker exec -it mongo mongosh

db-mysql-logs:
	$(DOCKER_COMPOSE) logs -f mysql

db-mongo-logs:
	$(DOCKER_COMPOSE) logs -f mongo

# Cleaning
clean:
	@echo "$(CYAN)Limpiando contenedores y volúmenes...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)✓ Limpio$(NC)"

status:
	@echo "$(CYAN)Estado de contenedores:$(NC)"
	docker ps --filter "label=com.docker.compose.project=ucc-arquisoft2"

ps:
	@echo "$(CYAN)Procesos docker:$(NC)"
	$(DOCKER_COMPOSE) ps

# Scripts
setup-env:
	@echo "$(CYAN)Configurando archivos .env...$(NC)"
	@for dir in users-api subscriptions-api activities-api payments-api search-api; do \
		if [ ! -f $$dir/.env ] && [ -f $$dir/.env.example ]; then \
			cp $$dir/.env.example $$dir/.env; \
			echo "$(GREEN)✓ Creado $$dir/.env$(NC)"; \
		fi; \
	done

# Alias
up: docker-up
down: docker-down
logs: docker-logs
restart: docker-restart
infra: local-up
clean-all: clean

.DEFAULT_GOAL := help
