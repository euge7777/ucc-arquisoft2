param(
    [switch]$Docker = $false,
    [switch]$Help = $false
)

function Show-Help {
    Write-Host "=========================================="
    Write-Host "SCRIPT DE INICIO - MICROSERVICIOS GYM"
    Write-Host "=========================================="
    Write-Host ""
    Write-Host "OPCIONES:"
    Write-Host "  .\start-dev.ps1              -> Levanta servicios en modo LOCAL"
    Write-Host "  .\start-dev.ps1 -Docker      -> Levanta TODO con Docker Compose"
    Write-Host "  .\start-dev.ps1 -Help        -> Muestra esta ayuda"
    Write-Host ""
    Write-Host "MODO LOCAL (defecto):"
    Write-Host "  - Infraestructura en Docker (MySQL, MongoDB, RabbitMQ, etc)"
    Write-Host "  - Microservicios en tu maquina (puedes debuggear)"
    Write-Host ""
    Write-Host "MODO DOCKER:"
    Write-Host "  - TODO en contenedores (infraestructura + microservicios)"
    Write-Host "  - Mas parecido a produccion"
    Write-Host ""
    Write-Host "REQUISITOS:"
    Write-Host "  - Docker Desktop instalado"
    Write-Host "  - Go 1.23+ instalado (para modo local)"
    Write-Host "  - PowerShell 5.0+"
    Write-Host ""
    Write-Host "PUERTOS UTILIZADOS:"
    Write-Host "  - 8080: Users API"
    Write-Host "  - 8081: Subscriptions API"
    Write-Host "  - 8082: Activities API"
    Write-Host "  - 8083: Payments API"
    Write-Host "  - 8084: Search API"
    Write-Host "  - 5173: Frontend React"
    Write-Host "  - 3306: MySQL"
    Write-Host "  - 27017: MongoDB"
    Write-Host "  - 5672: RabbitMQ"
    Write-Host "  - 15672: RabbitMQ Admin"
    Write-Host "  - 11211: Memcached"
    Write-Host "  - 8983: Solr"
    Write-Host ""
}

function Check-Prerequisites {
    Write-Host "Verificando prerequisitos..." -ForegroundColor Cyan

    if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
        Write-Host "ERROR: Docker no esta instalado o no esta en PATH" -ForegroundColor Red
        exit 1
    }
    Write-Host "  OK: Docker encontrado" -ForegroundColor Green

    if (!$Docker) {
        if (!(Get-Command go -ErrorAction SilentlyContinue)) {
            Write-Host "ERROR: Go no esta instalado o no esta en PATH (necesario para modo LOCAL)" -ForegroundColor Red
            Write-Host "  Usa: .\start-dev.ps1 -Docker" -ForegroundColor Yellow
            exit 1
        }
        Write-Host "  OK: Go encontrado" -ForegroundColor Green
    }
}

function Start-DockerMode {
    Write-Host ""
    Write-Host "=========================================="
    Write-Host "INICIANDO MODO DOCKER"
    Write-Host "=========================================="
    Write-Host ""

    Write-Host "Levantando infraestructura + microservicios..." -ForegroundColor Yellow
    docker compose -f docker-compose.new.yml up -d

    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Error al levantar los servicios" -ForegroundColor Red
        exit 1
    }

    Write-Host "OK: Esperando a que los servicios esten listos..." -ForegroundColor Green
    Start-Sleep -Seconds 5

    Show-ServiceStatus
}

function Start-LocalMode {
    Write-Host ""
    Write-Host "=========================================="
    Write-Host "INICIANDO MODO LOCAL"
    Write-Host "=========================================="
    Write-Host ""

    Write-Host "Levantando infraestructura..." -ForegroundColor Yellow
    docker compose -f docker-compose.new.yml up -d mysql mongo rabbitmq memcached solr

    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Error al levantar la infraestructura" -ForegroundColor Red
        exit 1
    }

    Write-Host "   OK: Infraestructura lista" -ForegroundColor Green
    Start-Sleep -Seconds 3

    $services = @(
        @{name="users-api"; port="8080"; path="users-api"},
        @{name="subscriptions-api"; port="8081"; path="subscriptions-api"},
        @{name="activities-api"; port="8082"; path="activities-api"},
        @{name="payments-api"; port="8083"; path="payments-api"},
        @{name="search-api"; port="8084"; path="search-api"}
    )

    foreach ($service in $services) {
        Write-Host "Iniciando $($service.name) (puerto $($service.port))..." -ForegroundColor Yellow

        $command = "cd '$($service.path)'; go mod download 2>`$null; go run cmd/api/main.go"
        Start-Process PowerShell -ArgumentList "-NoExit", "-Command", $command -WindowStyle Normal

        Start-Sleep -Seconds 1
    }

    Write-Host ""
    Write-Host "OK: Todos los servicios iniciados en nuevas ventanas PowerShell" -ForegroundColor Green
    Write-Host ""
    Write-Host "Microservicios disponibles:" -ForegroundColor Cyan
    foreach ($service in $services) {
        Write-Host "   - $($service.name): http://localhost:$($service.port)" -ForegroundColor White
    }
}

function Show-ServiceStatus {
    Write-Host ""
    Write-Host "ESTADO DE SERVICIOS:" -ForegroundColor Cyan

    $services = @(
        @{name="Users API"; url="http://localhost:8080/healthz"},
        @{name="Subscriptions API"; url="http://localhost:8081/healthz"},
        @{name="Activities API"; url="http://localhost:8082/healthz"},
        @{name="Payments API"; url="http://localhost:8083/healthz"},
        @{name="Search API"; url="http://localhost:8084/healthz"}
    )

    foreach ($service in $services) {
        try {
            $response = Invoke-WebRequest -Uri $service.url -TimeoutSec 2 -ErrorAction SilentlyContinue
            if ($response.StatusCode -eq 200) {
                Write-Host "OK: $($service.name): ACTIVO" -ForegroundColor Green
            }
        }
        catch {
            Write-Host "ESPERANDO: $($service.name): Iniciando..." -ForegroundColor Yellow
        }
    }

    Write-Host ""
    Write-Host "ACCESO A SERVICIOS:" -ForegroundColor Cyan
    Write-Host "   - Frontend:     http://localhost:5173" -ForegroundColor White
    Write-Host "   - RabbitMQ:     http://localhost:15672 (admin/admin)" -ForegroundColor White
    Write-Host "   - Solr:         http://localhost:8983/solr" -ForegroundColor White
}

function Show-Info {
    Write-Host ""
    Write-Host "=========================================="
    Write-Host "INFORMACION UTIL"
    Write-Host "=========================================="
    Write-Host ""
    Write-Host "Ver logs de Docker:" -ForegroundColor Yellow
    Write-Host "   docker compose -f docker-compose.new.yml logs -f" -ForegroundColor White

    Write-Host ""
    Write-Host "Detener todo:" -ForegroundColor Yellow
    Write-Host "   docker compose -f docker-compose.new.yml down" -ForegroundColor White

    Write-Host ""
    Write-Host "Reiniciar servicios:" -ForegroundColor Yellow
    Write-Host "   docker compose -f docker-compose.new.yml restart" -ForegroundColor White

    Write-Host ""
    Write-Host "Entrar a MySQL:" -ForegroundColor Yellow
    Write-Host "   docker exec -it mysql mysql -uroot -proot123" -ForegroundColor White

    Write-Host ""
    Write-Host "Entrar a MongoDB:" -ForegroundColor Yellow
    Write-Host "   docker exec -it mongo mongosh" -ForegroundColor White
    Write-Host ""
}

if ($Help) {
    Show-Help
    exit 0
}

Check-Prerequisites

if ($Docker) {
    Start-DockerMode
}
else {
    Start-LocalMode
}

Show-ServiceStatus
Show-Info
