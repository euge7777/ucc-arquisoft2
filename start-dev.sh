#!/bin/bash

# Script para levantar todos los microservicios en desarrollo (Linux/Mac)
# Uso: chmod +x start-dev.sh && ./start-dev.sh [OPTIONS]

set -e

DOCKER_MODE=false
HELP=false

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Funciones
show_help() {
    cat << EOF
${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}
${CYAN}║          SCRIPT DE INICIO - MICROSERVICIOS GYM SYSTEM          ║${NC}
${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}

OPCIONES:
  ./start-dev.sh              → Levanta servicios en modo LOCAL
  ./start-dev.sh --docker     → Levanta TODO con Docker Compose
  ./start-dev.sh --help       → Muestra esta ayuda

MODO LOCAL (defecto):
  - Infraestructura en Docker (MySQL, MongoDB, RabbitMQ, etc)
  - Microservicios en tu máquina (puedes debuggear)
  - Cada servicio en un terminal separado (tmux/screen)

MODO DOCKER:
  - TODO en contenedores (infraestructura + microservicios)
  - Más parecido a producción
  - Ver logs: docker compose -f docker-compose.new.yml logs -f

REQUISITOS:
  - Docker y Docker Compose instalados
  - Go 1.23+ instalado (para modo local)
  - tmux o screen (para modo local, opcional)

PUERTOS UTILIZADOS:
  - 8080: Users API
  - 8081: Subscriptions API
  - 8082: Activities API
  - 8083: Payments API
  - 8084: Search API
  - 5173: Frontend React
  - 3306: MySQL
  - 27017: MongoDB
  - 5672: RabbitMQ
  - 15672: RabbitMQ Admin
  - 11211: Memcached
  - 8983: Solr

EOF
}

check_prerequisites() {
    echo -e "${CYAN}✓ Verificando prerequisitos...${NC}"

    if ! command -v docker &> /dev/null; then
        echo -e "${RED}✗ Docker no está instalado${NC}"
        exit 1
    fi
    echo -e "${GREEN}  ✓ Docker encontrado${NC}"

    if ! docker compose version &> /dev/null; then
        echo -e "${RED}✗ Docker Compose no está disponible${NC}"
        exit 1
    fi
    echo -e "${GREEN}  ✓ Docker Compose encontrado${NC}"

    if [ "$DOCKER_MODE" = false ]; then
        if ! command -v go &> /dev/null; then
            echo -e "${RED}✗ Go no está instalado (necesario para modo LOCAL)${NC}"
            echo -e "${YELLOW}  Usa: ./start-dev.sh --docker${NC}"
            exit 1
        fi
        echo -e "${GREEN}  ✓ Go encontrado${NC}"
    fi
}

start_docker_mode() {
    echo -e "${CYAN}\n╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║            INICIANDO MODO DOCKER (TODO EN CONTENEDORES)        ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

    echo -e "${YELLOW}Levantando infraestructura + microservicios...${NC}"
    docker compose -f docker-compose.new.yml up -d

    echo -e "${GREEN}✓ Esperando a que los servicios estén listos...${NC}"
    sleep 5

    show_service_status
}

start_local_mode() {
    echo -e "${CYAN}\n╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║        INICIANDO MODO LOCAL (Infraestructura en Docker)         ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

    # Levantar infraestructura
    echo -e "${YELLOW}1️⃣  Levantando infraestructura...${NC}"
    docker compose -f docker-compose.new.yml up -d mysql mongo rabbitmq memcached solr

    echo -e "${GREEN}   ✓ Infraestructura lista${NC}\n"
    sleep 3

    # Iniciar servicios
    local services=(
        "users-api:8080"
        "subscriptions-api:8081"
        "activities-api:8082"
        "payments-api:8083"
        "search-api:8084"
    )

    for service in "${services[@]}"; do
        local name="${service%:*}"
        local port="${service#*:}"

        echo -e "${YELLOW}2️⃣  Iniciando $name (puerto $port)...${NC}"

        # Crear .env si no existe
        if [ ! -f "$name/.env" ] && [ -f "$name/.env.example" ]; then
            cp "$name/.env.example" "$name/.env"
            echo -e "${GREEN}   ✓ Archivo .env creado${NC}"
        fi

        # Si tmux está disponible, usar tmux
        if command -v tmux &> /dev/null; then
            tmux new-session -d -s "$name" -c "$name" "go mod download && go run cmd/api/main.go"
            echo -e "${GREEN}   ✓ $name iniciado en sesión tmux${NC}"
        else
            # Si no, abrir en background
            cd "$name"
            go mod download > /dev/null 2>&1
            go run cmd/api/main.go > "../logs-$name.txt" 2>&1 &
            echo -e "${GREEN}   ✓ $name iniciado (logs en logs-$name.txt)${NC}"
            cd ..
        fi

        sleep 1
    done

    if command -v tmux &> /dev/null; then
        echo -e "${GREEN}\n✓ Todos los servicios iniciados en sesiones tmux${NC}"
        echo -e "${CYAN}Ver sesiones: tmux list-sessions${NC}"
        echo -e "${CYAN}Conectarse: tmux attach-session -t <nombre>${NC}"
    else
        echo -e "${GREEN}\n✓ Todos los servicios iniciados en background${NC}"
        echo -e "${CYAN}Ver logs: tail -f logs-<servicio>.txt${NC}"
    fi

    echo -e "\n${CYAN}📋 Microservicios disponibles:${NC}"
    echo -e "   - Users API:         http://localhost:8080"
    echo -e "   - Subscriptions API: http://localhost:8081"
    echo -e "   - Activities API:    http://localhost:8082"
    echo -e "   - Payments API:      http://localhost:8083"
    echo -e "   - Search API:        http://localhost:8084"
}

show_service_status() {
    echo -e "\n${CYAN}📋 ESTADO DE SERVICIOS:${NC}"

    local services=(
        "Users API:http://localhost:8080/healthz"
        "Subscriptions API:http://localhost:8081/healthz"
        "Activities API:http://localhost:8082/healthz"
        "Payments API:http://localhost:8083/healthz"
        "Search API:http://localhost:8084/healthz"
    )

    for service in "${services[@]}"; do
        local name="${service%:*}"
        local url="${service#*:}"

        if curl -s "$url" > /dev/null 2>&1; then
            echo -e "${GREEN}✓ $name: ACTIVO${NC}"
        else
            echo -e "${YELLOW}⏳ $name: Iniciando...${NC}"
        fi
    done

    echo -e "\n${CYAN}🌐 ACCESO A SERVICIOS:${NC}"
    echo -e "   - Frontend:     http://localhost:5173"
    echo -e "   - RabbitMQ:     http://localhost:15672 (admin/admin)"
    echo -e "   - Solr:         http://localhost:8983/solr"
}

show_info() {
    echo -e "\n${CYAN}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║                     INFORMACIÓN ÚTIL                            ║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════════╝${NC}\n"

    echo -e "${YELLOW}📖 Ver logs de Docker:${NC}"
    echo -e "   docker compose -f docker-compose.new.yml logs -f\n"

    echo -e "${YELLOW}🛑 Detener todo:${NC}"
    echo -e "   docker compose -f docker-compose.new.yml down\n"

    echo -e "${YELLOW}🔄 Reiniciar servicios:${NC}"
    echo -e "   docker compose -f docker-compose.new.yml restart\n"

    echo -e "${YELLOW}📦 Entrar a MySQL:${NC}"
    echo -e "   docker exec -it mysql mysql -uroot -proot123\n"

    echo -e "${YELLOW}📦 Entrar a MongoDB:${NC}"
    echo -e "   docker exec -it mongo mongosh\n"

    if command -v tmux &> /dev/null; then
        echo -e "${YELLOW}📋 Comandos tmux útiles:${NC}"
        echo -e "   tmux list-sessions        # Ver todas las sesiones"
        echo -e "   tmux attach -t <nombre>   # Conectarse a sesión"
        echo -e "   tmux kill-session -t <nombre>  # Matar sesión"
        echo -e "   Ctrl+B c                  # Nueva ventana en tmux"
        echo -e "   Ctrl+B [número]           # Cambiar ventana en tmux\n"
    fi
}

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        --docker)
            DOCKER_MODE=true
            shift
            ;;
        --help)
            HELP=true
            shift
            ;;
        *)
            echo "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Main
if [ "$HELP" = true ]; then
    show_help
    exit 0
fi

check_prerequisites

if [ "$DOCKER_MODE" = true ]; then
    start_docker_mode
else
    start_local_mode
fi

show_service_status
show_info
