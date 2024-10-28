# Cargar las variables del archivo .env
# Cargar las variables del archivo .env usando export
include .env
export $(shell sed 's/=.*//' .env)

# Comando para levantar Postgres en Docker
start-db:
	@echo "Iniciando Postgres en Docker..."
	@docker run --name $(DB_CONTAINER_NAME) \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASS) \
	-e POSTGRES_DB=$(DB_NAME) \
	-p $(DB_PORT):5432 \
	-d $(DB_IMAGE) 

# Comando para detener y eliminar el contenedor de postgres
stop-db:
	@echo "Deteniendo y eliminando contenedor PostgreSQL..."
	@docker stop $(DB_CONTAINER_NAME) || true
	@docker rm $(DB_CONTAINER_NAME) || true

# Comando para compilar la aplicacion Go en la carpeta bin
build:
	@echo "Compilando la aplicacion en $(BIN_DIR)"
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(GO_APP)

# Comando para ejecutar la aplicacion compilada
run:
	@echo "Ejecutando la aplicacion"
	@$(BIN_DIR)/$(APP_NAME)

# Comando para correr las pruebas 
test:
	@echo "Ejecutando las pruebas..."
	@go test ./... -v

# Comando para levantar postgres y ejecutar la App
run-w-db: start-db build
	@echo "Esperando que PostgreSQL esté listo..."
	@sleep 10 # Espera unos segundos para que PostgreSQL esté listo
	@make run

# Comando para deter la aplicacion y la base de datos 
stop:
	@make stop-db

# Limpia contenedores, imágenes y volúmenes de Docker
clean:
	@echo "Limpiando contenedores, imágenes y volúmenes no utilizados..."
	@docker system prune -f

# Limpia los binarios compilados 
clean-bin:
	@echo "Limpiando binarios compilados..."
	@rm -rf $(BIN_DIR)

# Comando para correr la aplicación en modo producción (compilada en bin/)
run-prod: start-db build
	@echo "Ejecutando la aplicación en modo producción..."
	@./$(BIN_DIR)/$(APP_NAME)