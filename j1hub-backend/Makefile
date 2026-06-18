# J1 Hub Backend Makefile

.PHONY: run build test lint migrate-up migrate-down docker-up docker-down

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

lint:
	golangci-lint run

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
