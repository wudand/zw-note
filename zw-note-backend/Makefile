# ============================================================
# Variables
# ============================================================
APP_NAME   := go-web-api
BUILD_DIR  := bin
MAIN       := ./cmd/server/main.go
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    := -ldflags "-X main.Version=$(VERSION) -s -w"
GOFLAGS    := CGO_ENABLED=0

.PHONY: all build run test lint vet fmt clean \
        docker-build docker-run docker-stop \
        pg-up pg-down pg-logs pg-shell pg-wait migrate seed dev \
        mod-tidy mod-download generate swagger

# ============================================================
# Default
# ============================================================
all: lint test build

# ============================================================
# Build
# ============================================================
build:
	@echo ">> Building $(APP_NAME) $(VERSION)"
	@mkdir -p $(BUILD_DIR)
	$(GOFLAGS) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)

build-linux:
	@echo ">> Cross-compiling for linux/amd64"
	@mkdir -p $(BUILD_DIR)
	$(GOFLAGS) GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(MAIN)

# ============================================================
# Run
# ============================================================
run:
	go run $(MAIN) start

run-config:
	go run $(MAIN) start --config configs/config.yaml

# ============================================================
# Test
# ============================================================
test:
	go test -v -race -cover ./...

test-short:
	go test -short ./...

# ============================================================
# Code quality
# ============================================================
lint:
	@which golangci-lint > /dev/null 2>&1 || \
		(echo "golangci-lint not found, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1)
	golangci-lint run ./...

vet:
	go vet ./...

fmt:
	gofmt -w -s .
	goimports -w . 2>/dev/null || true

# ============================================================
# Docker（配置在仓库根 docker/ 目录）
# ============================================================
COMPOSE := docker compose -f ../docker/docker-compose.yml

docker-build:
	$(COMPOSE) build backend

docker-run:
	$(COMPOSE) up -d backend

docker-stop:
	$(COMPOSE) stop backend

# ============================================================
# PostgreSQL（默认连接本机已有实例，不启动新容器）
# ============================================================
PG_HOST     ?= 127.0.0.1
PG_PORT     ?= 5432
PG_USER     ?= hsop
PG_PASSWORD ?= P0stgre2022
PG_DB       ?= zw-note
export PGPASSWORD := $(PG_PASSWORD)

# 可选：临时起一个本项目独立库（与本机已有 postgres_data 分离）
pg-up:
	@echo ">> Starting local-db profile PostgreSQL..."
	$(COMPOSE) --profile local-db up -d postgres

pg-down:
	@echo ">> Stopping local-db PostgreSQL..."
	$(COMPOSE) --profile local-db down

pg-logs:
	$(COMPOSE) --profile local-db logs -f postgres

pg-shell:
	psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(PG_DB)

pg-wait:
	@echo ">> Waiting for PostgreSQL at $(PG_HOST):$(PG_PORT)/$(PG_DB)..."
	@until pg_isready -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(PG_DB) >/dev/null 2>&1; do \
		echo "  PostgreSQL not ready, retrying in 2s..."; sleep 2; \
	done
	@echo ">> PostgreSQL ready"

migrate:
	@echo ">> Running migrations..."
	@for f in scripts/migrations/*.sql; do \
		echo "  Running $$f..."; \
		psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(PG_DB) -v ON_ERROR_STOP=1 -f "$$f"; \
	done
	@echo ">> Migrate done"

# 仅 dev/test 环境：执行 scripts/seed/*.sql（幂等）
seed:
	@echo ">> Seeding dev/test data..."
	@for f in scripts/seed/*.sql; do \
		echo "  Running $$f..."; \
		psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(PG_DB) -v ON_ERROR_STOP=1 -f "$$f"; \
	done
	@echo ">> Seed done"

# 一键：等待本机 DB → migrate → seed → 启动 API（不启动新 Postgres）
dev: pg-wait
	@$(MAKE) migrate
	@$(MAKE) seed
	@echo ">> Starting API server..."
	@$(MAKE) run

# ============================================================
# Modules
# ============================================================
mod-tidy:
	go mod tidy

mod-download:
	go mod download

generate:
	go generate ./...

# ============================================================
# Swagger
# ============================================================
swagger:
	@echo ">> Generating Swagger docs..."
	@swag init -g cmd/server/main.go --parseDependency --parseInternal || \
		go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/server/main.go --parseDependency --parseInternal
	@echo ">> Swagger docs generated in docs/"

# ============================================================
# Clean
# ============================================================
clean:
	@rm -rf $(BUILD_DIR)
	@echo ">> Cleaned"
