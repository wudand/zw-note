# go-web-api

A production-ready RESTful API framework built with Go, following Go idioms: simplicity, explicitness, composition, and concurrency.

同时支撑 **PC 管理后台**（`/api/admin/v1`）与 **微信小程序**（`/api/mp/v1`）两套 API。

## Tech Stack

| Concern        | Choice                           |
|----------------|----------------------------------|
| Web framework  | [Gin](https://github.com/gin-gonic/gin) |
| CLI            | [Cobra](https://github.com/spf13/cobra) |
| Config         | [Viper](https://github.com/spf13/viper) |
| Logger         | [Uber Zap](https://github.com/uber-go/zap) |
| Database       | PostgreSQL 16 + [sqlx](https://github.com/jmoiron/sqlx) + [pgx](https://github.com/jackc/pgx) |
| Auth           | JWT ([golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)) |
| Metrics        | [Prometheus](https://github.com/prometheus/client_golang) |

## Project Layout

```
go-web-api/
├── cmd/server/main.go          # Entry point (Cobra CLI)
├── bootstrap/                  # App wiring & router setup
├── internal/
│   ├── api/admin/              # Admin HTTP handlers
│   ├── api/mp/                 # Mini-program HTTP handlers
│   ├── service/                # Business logic (interface-first)
│   ├── repository/             # Data access (interface-first)
│   ├── model/                  # DB structs (sqlx tags)
│   ├── dto/                    # Request / Response types
│   ├── middleware/             # auth, rbac, cors, logger, metrics
│   └── config/                 # Config struct + Viper loader
├── pkg/
│   ├── database/               # PostgreSQL connection pool helper
│   ├── logger/                 # Zap initialisation
│   └── utils/                  # JWT, unified response, app errors
├── configs/config.yaml         # Default configuration
└── scripts/migrations/         # SQL migration files
└── scripts/seed/               # Dev/test seed data
```

## Quick Start

### 1. Prerequisites

- Go 1.22+
- PostgreSQL 16（或 Docker）

### 2. Create the database and run migrations

```bash
createdb go_web_api
psql -d go_web_api -f scripts/migrations/001_create_users_table.sql
# …按序执行 scripts/migrations/*.sql
```

使用 Docker 时，首次启动容器会自动执行 `scripts/migrations/`。

### 3. Configure

Edit `configs/config.yaml` or override via environment variables:

```bash
export APP_DATABASE_DSN="postgres://postgres:password@127.0.0.1:5432/go_web_api?sslmode=disable"
export APP_SERVER_MODE="release"
```

### 4. Run

```bash
# With make
make run

# Or directly
go run ./cmd/server/main.go start

# With a custom config file
go run ./cmd/server/main.go start --config /path/to/config.yaml
```

### 5. 使用 Docker 快速启动测试环境（推荐）

```bash
make dev                  # 一键启动 PostgreSQL + seed + API
```

或分步执行：

```bash
make pg-up                # 启动 PostgreSQL 16 容器
make pg-wait              # 等待 PostgreSQL 就绪
make seed                 # 写入测试数据
make run                  # 启动 API
```

其他命令：

```bash
make pg-down              # 停止并移除 PostgreSQL 容器
make pg-logs              # 查看 PostgreSQL 日志
make pg-shell             # 进入 psql
```

### 6. Build & Docker

Docker 相关文件统一放在仓库根目录 `docker/`（见 `docker/README.md`）。

```bash
make build                # native binary → bin/go-web-api
make docker-build         # 构建后端镜像
make docker-run           # 启动后端容器 :8080
```

或在仓库根目录：

```bash
docker compose -f docker/docker-compose.yml up -d --build
```

## API Reference

### Admin（`/api/admin/v1`）

| Method | Path                    | Description      |
|--------|-------------------------|------------------|
| POST   | `/auth/login`           | Admin login → JWT |
| GET    | `/users`                | List admin users |
| CRUD   | `/categories`           | Product categories |
| CRUD   | `/products`             | Products |
| CRUD   | `/coupons`              | Coupons |
| CRUD   | `/carousels`            | Home carousels |
| CRUD   | `/redemption-codes`     | Redemption codes |

### Mini-Program（`/api/mp/v1`）

| Method | Path                              | Description        |
|--------|-----------------------------------|--------------------|
| POST   | `/auth/wx-login`                  | WeChat login → JWT |
| GET    | `/categories` / `/products`       | Public catalog     |
| GET    | `/carousels` / `/coupons`         | Public content     |
| CRUD   | `/addresses`                      | JWT required       |
| GET    | `/coupons/my`                     | My coupons         |
| POST   | `/coupons/:id/claim`              | Claim coupon       |
| POST   | `/redemption-codes/validate`      | Validate code      |

### Observability

| Path       | Description                   |
|------------|-------------------------------|
| `GET /health`  | Liveness + DB ping        |
| `GET /metrics` | Prometheus metrics        |
| `GET /swagger/*any` | Swagger UI           |

### Response Envelope

Every response is wrapped in:

```json
{
  "code":    0,
  "message": "success",
  "data":    { ... }
}
```

Error codes: `0` = success, `1xxxx` = generic, `2xxxx` = admin, `3xxxx` = mini-program, `4xxxx` = domain.

### Example: Admin login

```bash
curl -X POST http://localhost:8888/api/admin/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"superadmin","password":"Admin@123456"}'
```

## Configuration Reference

All keys can be overridden with `APP_<KEY>` env vars (dots → underscores), e.g. `APP_SERVER_PORT=9090`.

| Key                        | Default  | Description                      |
|----------------------------|----------|----------------------------------|
| `server.port`              | `8888`   | Listening port                   |
| `server.mode`              | `debug`  | `debug` / `release` / `test`     |
| `server.read_timeout`      | `30`     | HTTP read timeout (seconds)      |
| `server.write_timeout`     | `30`     | HTTP write timeout (seconds)     |
| `database.dsn`             | —        | PostgreSQL DSN string            |
| `database.max_open_conns`  | `20`     | Connection pool max open         |
| `jwt.admin_secret`         | —        | **Change in production!**        |
| `jwt.mp_secret`            | —        | **Change in production!**        |
| `jwt.admin_expire_hours`   | `8`      | Admin token TTL                  |
| `jwt.mp_expire_hours`      | `720`    | MP token TTL                     |
| `log.level`                | `info`   | `debug` / `info` / `warn` / `error` |
| `log.format`               | `json`   | `json` / `console`               |

## Development

```bash
make test       # run all tests with race detector
make vet        # go vet
make lint       # golangci-lint (install separately)
make fmt        # gofmt + goimports
```
