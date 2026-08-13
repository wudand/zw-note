# zw-note

个人知识笔记：文档列表、Markdown 编辑、全文预览。本仓库是前后端一体的单体仓库，push 到 `main` 后由 GitHub Actions 构建镜像并部署。

## 功能

- 文档列表：查看、定位、编辑标题 / 描述 / 标签
- 文档编辑：左侧目录大纲 + 右侧 Markdown 编辑器
- 全文预览：只读阅读，不含编辑操作

## 仓库结构

```
notes/
├── zw-note-frontend/    # Vue 3 笔记前端
├── zw-note-backend/     # Go API（笔记 / 管理后台 / 小程序）
├── docker/              # 前后端镜像与 compose
├── Makefile             # 本地一键启停
├── 数据库设计.md
└── 需求分析.md
```

## 技术栈

| 层 | 选型 |
|----|------|
| 前端 | Vue 3 + Vite + TypeScript + Element Plus |
| 后端 | Go + Gin + PostgreSQL |
| 部署 | Docker Compose；push `main` 后 GitHub Actions 发版 |

## 环境要求

- Go 1.25+
- Node.js >= 22.20
- PostgreSQL 16
- Docker（可选，部署用）

首次启动前端前需在 `zw-note-frontend/` 执行 `npm install`。

## 快速开始

命令均在**仓库根目录**执行。

### 本地开发

```bash
# 1. 建库 + 迁移 + 种子（需本机或容器里已有 Postgres）
make db-setup-docker   # 或 make db-setup-local

# 2. 启动前后端
make start

# 3. 查看状态
make status
```

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:5555 |
| 后端 | http://localhost:8004 |
| Swagger | http://localhost:8004/swagger/index.html |

### Docker

```bash
cp docker/.env.example docker/.env   # 按需填写数据库连接，勿提交 .env
make docker-up
```

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:8005 |
| 后端 | http://localhost:8006 |
| Swagger | http://localhost:8006/swagger/index.html |

compose、数据库依赖与生产部署见 [docker/README.md](docker/README.md)。

## 常用命令

| 命令 | 说明 |
|------|------|
| `make start` / `stop` / `restart` | 本地启停前后端 |
| `make logs` | 跟踪日志；`make logs S=backend` 只看一侧 |
| `make status` | 启动状态与访问地址 |
| `make docker-up` / `docker-down` | 构建并启动 / 停止 compose |
| `make docker-logs` / `docker-status` | Docker 日志与状态 |
| `make db-setup-docker` | 容器 `postgres-dev` 建库 + 迁移 + 种子 |
| `make db-setup-local` | 本机 `psql` 建库 + 迁移 + 种子 |

## 文档索引

| 文档 | 内容 |
|------|------|
| [zw-note-frontend/README.md](zw-note-frontend/README.md) | 前端功能与本地开发 |
| [zw-note-backend/README.md](zw-note-backend/README.md) | 后端 API、配置与开发 |
| [docker/README.md](docker/README.md) | Docker、生产部署、CI Secrets |
| [数据库设计.md](数据库设计.md) | 表结构 |
| [需求分析.md](需求分析.md) | 需求说明 |
