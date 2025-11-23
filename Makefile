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