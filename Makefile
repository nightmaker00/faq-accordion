.PHONY: help run build fmt vet lint lint-install swagger docker-up docker-down docker-logs migrate-up migrate-down


help:
	@echo "Targets:"
	@echo "  run          - run API locally"
	@echo "  build        - build binary to ./bin/app"
	@echo "  test         - run tests"
	@echo "  fmt          - gofmt all packages"
	@echo "  vet          - go vet all packages"
	@echo "  lint         - run golangci-lint"
	@echo "  lint-install - install golangci-lint"
	@echo "  swagger      - generate swagger docs"
	@echo "  docker-up    - start postgres + api in docker"
	@echo "  docker-down  - stop docker services"
	@echo "  docker-logs  - follow api logs"
	@echo "  migrate-up   - apply migrations (needs docker)"
	@echo "  migrate-down - rollback migrations (needs docker)"

run:
	go run ./cmd/app

build:
	mkdir -p bin
	go build -o ./bin/app ./cmd/app

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

swagger:
	swag init -g cmd/app/main.go

docker-up:
	docker compose -f deployments/docker-compose.yml up -d --build

docker-down:
	docker compose -f deployments/docker-compose.yml down -v

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f api

migrate-up:
	docker compose -f deployments/docker-compose.yml exec -T postgres \
		psql -U postgres -d faq -f /migrations/000001_faq_structure.up.sql

migrate-down:
	docker compose -f deployments/docker-compose.yml exec -T postgres \
		psql -U postgres -d faq -f /migrations/000001_faq_structure.down.sql

