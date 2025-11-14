@echo off
REM Script para levantar todos los microservicios en desarrollo (Windows - CMD)
REM Uso: start-dev.bat [OPTIONS]

setlocal enabledelayedexpansion

set DOCKER_MODE=false
set HELP=false

REM Parsear argumentos
:parse_args
if "%1"=="" goto check_prereqs
if "%1"=="--docker" (
    set DOCKER_MODE=true
    shift
    goto parse_args
)
if "%1"=="-docker" (
    set DOCKER_MODE=true
    shift
    goto parse_args
)
if "%1"=="--help" (
    set HELP=true
    shift
    goto parse_args
)
if "%1"=="-help" (
    set HELP=true
    shift
    goto parse_args
)
shift
goto parse_args

:check_help
if "%HELP%"=="true" (
    goto show_help
)

:check_prereqs
cls
echo.
echo Verificando prerequisitos...

REM Verificar Docker
docker --version >nul 2>&1
if errorlevel 1 (
    color 04
    echo.
    echo [ERROR] Docker no esta instalado o no esta en PATH
    echo.
    pause
    exit /b 1
)
echo [OK] Docker encontrado

REM Verificar Docker Compose
docker compose version >nul 2>&1
if errorlevel 1 (
    color 04
    echo.
    echo [ERROR] Docker Compose no esta disponible
    echo.
    pause
    exit /b 1
)
echo [OK] Docker Compose encontrado

REM Para modo local, verificar Go
if "%DOCKER_MODE%"=="false" (
    go version >nul 2>&1
    if errorlevel 1 (
        color 04
        echo.
        echo [ERROR] Go no esta instalado (necesario para modo LOCAL)
        echo Usa: start-dev.bat --docker
        echo.
        pause
        exit /b 1
    )
    echo [OK] Go encontrado
)

if "%DOCKER_MODE%"=="true" (
    goto docker_mode
) else (
    goto local_mode
)

:docker_mode
cls
echo.
echo ========================================================================
echo              INICIANDO MODO DOCKER (TODO EN CONTENEDORES)
echo ========================================================================
echo.
echo Levantando infraestructura + microservicios...
docker compose -f docker-compose.new.yml up -d

if errorlevel 1 (
    color 04
    echo [ERROR] Error al levantar los servicios
    pause
    exit /b 1
)

echo.
echo Esperando a que los servicios esten listos...
timeout /t 5 /nobreak

goto show_status

:local_mode
cls
echo.
echo ========================================================================
echo         INICIANDO MODO LOCAL (Infraestructura en Docker)
echo ========================================================================
echo.
echo Levantando infraestructura...
docker compose -f docker-compose.new.yml up -d mysql mongo rabbitmq memcached solr

if errorlevel 1 (
    color 04
    echo [ERROR] Error al levantar la infraestructura
    pause
    exit /b 1
)

echo [OK] Infraestructura lista
echo.
echo Iniciando microservicios...
echo.

REM Iniciar servicios
echo.
echo IMPORTANTE: Los servicios se abriran en nuevas ventanas CMD
echo.
pause

REM Users API
echo Iniciando users-api...
start "Users API (8080)" cmd /k "cd users-api && (if not exist .env copy .env.example .env) && go mod download && go run cmd/api/main.go"
timeout /t 2 /nobreak

REM Subscriptions API
echo Iniciando subscriptions-api...
start "Subscriptions API (8081)" cmd /k "cd subscriptions-api && (if not exist .env copy .env.example .env) && go mod download && go run cmd/api/main.go"
timeout /t 2 /nobreak

REM Activities API
echo Iniciando activities-api...
start "Activities API (8082)" cmd /k "cd activities-api && (if not exist .env copy .env.example .env) && go mod download && go run cmd/api/main.go"
timeout /t 2 /nobreak

REM Payments API
echo Iniciando payments-api...
start "Payments API (8083)" cmd /k "cd payments-api && (if not exist .env copy .env.example .env) && go mod download && go run cmd/api/main.go"
timeout /t 2 /nobreak

REM Search API
echo Iniciando search-api...
start "Search API (8084)" cmd /k "cd search-api && (if not exist .env copy .env.example .env) && go mod download && go run cmd/api/main.go"

echo.
echo [OK] Todos los servicios iniciados en nuevas ventanas CMD
echo.

goto show_status

:show_status
cls
echo.
echo ========================================================================
echo                      ESTADO DE SERVICIOS
echo ========================================================================
echo.

echo Verificando salud de servicios (puede tomar unos segundos)...
echo.

REM Verificar cada servicio
set services=(8080 8081 8082 8083 8084)
set names=("Users API" "Subscriptions API" "Activities API" "Payments API" "Search API")

setlocal enabledelayedexpansion
for /L %%i in (0,1,4) do (
    set port=!services[%%i]!
    set name=!names[%%i]!

    timeout /t 1 /nobreak >nul
    curl -s http://localhost:!port!/healthz >nul 2>&1
    if errorlevel 1 (
        echo [LOADING] !name! (puerto !port!)
    ) else (
        echo [OK] !name! (puerto !port!)
    )
)

echo.
echo ========================================================================
echo                     ACCESO A SERVICIOS
echo ========================================================================
echo.
echo Frontend:     http://localhost:5173
echo RabbitMQ:     http://localhost:15672 (usuario: admin, password: admin)
echo Solr:         http://localhost:8983/solr
echo.

echo ========================================================================
echo                      INFORMACION UTIL
echo ========================================================================
echo.
echo Ver logs de Docker:
echo   docker compose -f docker-compose.new.yml logs -f
echo.
echo Detener todo:
echo   docker compose -f docker-compose.new.yml down
echo.
echo Reiniciar servicios:
echo   docker compose -f docker-compose.new.yml restart
echo.
echo Entrar a MySQL:
echo   docker exec -it mysql mysql -uroot -proot123
echo.
echo Entrar a MongoDB:
echo   docker exec -it mongo mongosh
echo.
pause
exit /b 0

:show_help
cls
echo.
echo ========================================================================
echo          SCRIPT DE INICIO - MICROSERVICIOS GYM SYSTEM
echo ========================================================================
echo.
echo OPCIONES:
echo   start-dev.bat              - Levanta servicios en modo LOCAL
echo   start-dev.bat --docker     - Levanta TODO con Docker Compose
echo   start-dev.bat --help       - Muestra esta ayuda
echo.
echo MODO LOCAL (defecto):
echo   - Infraestructura en Docker (MySQL, MongoDB, RabbitMQ, etc)
echo   - Microservicios en tu maquina (puedes debuggear)
echo   - Cada servicio en una ventana CMD separada
echo.
echo MODO DOCKER:
echo   - TODO en contenedores (infraestructura + microservicios)
echo   - Mas parecido a produccion
echo   - Ver logs: docker compose -f docker-compose.new.yml logs -f
echo.
echo REQUISITOS:
echo   - Docker Desktop instalado
echo   - Go 1.23+ instalado (para modo local)
echo.
echo PUERTOS UTILIZADOS:
echo   - 8080: Users API
echo   - 8081: Subscriptions API
echo   - 8082: Activities API
echo   - 8083: Payments API
echo   - 8084: Search API
echo   - 5173: Frontend React
echo   - 3306: MySQL
echo   - 27017: MongoDB
echo   - 5672: RabbitMQ (AMQP)
echo   - 15672: RabbitMQ Admin
echo   - 11211: Memcached
echo   - 8983: Solr
echo.
pause
exit /b 0
