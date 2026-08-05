# Docker 部署

本目录集中管理前后端部署相关配置。

```
docker/
├── docker-compose.yml      # 统一编排
├── backend/
│   ├── Dockerfile
│   └── .dockerignore       # 参考；实际生效的是仓库根 .dockerignore
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── .dockerignore
└── README.md
```

构建上下文为**仓库根目录**（`notes/`），便于同时拷贝代码与本目录下的 nginx 配置。

## 常用命令

在仓库根目录执行：

```bash
# 构建并启动前后端
docker compose -f docker/docker-compose.yml up -d --build

# 仅构建后端镜像
docker compose -f docker/docker-compose.yml build backend

# 仅构建前端镜像
docker compose -f docker/docker-compose.yml build frontend

# 可选：临时启动独立 PostgreSQL
docker compose -f docker/docker-compose.yml --profile local-db up -d postgres

# 停止
docker compose -f docker/docker-compose.yml down
```

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost |
| 后端 | http://localhost:8080 |
| Postgres（local-db） | localhost:5432 |

## 后端 Makefile 快捷方式

在 `zw-note-backend/` 下仍可使用：

```bash
make docker-build
make docker-run
make pg-up
make pg-down
```
