ifneq (,$(wildcard .env))
    include .env
    export
endif

APP_NAME=avito-app

run:
	go run main.go

build:
	go build -o bin/$(APP_NAME) main.go

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	docker compose run migrator up

migrate-down:
	docker compose run migrator down 1

lint:
	golangci-lint run --timeout 10000s

test:
	go test ./...

swagger-generate:
	swag init -g main.go -o ./docs/swagger

help:
	@echo "Доступные команды:"
	@echo "  run           - Запустить приложение"
	@echo "  build         - Собрать приложение"
	@echo "  docker-up     - Запустить Docker контейнеры"
	@echo "  docker-down   - Остановить Docker контейнеры"
	@echo "  migrate-up    - Применить миграции базы данных"
	@echo "  migrate-down  - Откатить миграции базы данных"
	@echo "  lint          - Запустить линтер"
	@echo "  test          - Запустить тесты"
	@echo "  swagger-generate - Сгенерировать Swagger документацию"
	@echo "  help          - Показать эту справку"