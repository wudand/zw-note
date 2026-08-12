# Docker 部署

本目录集中管理前后端部署相关配置。

```
docker/
├── docker-compose.yml      # 统一编排
├── .env.example            # 数据库连接参数模板，复制为 .env 后生效
├── backend/
│   ├── Dockerfile
│   └── .dockerignore       # 参考；实际生效的是仓库根 .dockerignore
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   └── .dockerignore       # 参考；实际生效的是仓库根 .dockerignore
└── README.md
```

构建上下文为**仓库根目录**（`notes/`），便于同时拷贝代码与本目录下的 nginx 配置。

## 常用命令

首次使用需先准备好数据库连接参数（`.env` 已加入 `.gitignore`，不会被提交）：

```bash
cp docker/.env.example docker/.env
# 按需修改 docker/.env 里的 DB_USER / DB_PASSWORD / DB_HOST / DB_PORT / DB_NAME
```

在仓库根目录执行：

```bash
# 构建并启动前后端
docker compose -f docker/docker-compose.yml up -d --build

# 仅构建后端镜像
docker compose -f docker/docker-compose.yml build backend

# 仅构建前端镜像
docker compose -f docker/docker-compose.yml build frontend

# 可选：临时启动独立 PostgreSQL（会占用宿主机 5432 端口，
# 若本机已有 postgres-dev 容器占用该端口需先停掉，避免端口冲突）
docker compose -f docker/docker-compose.yml --profile local-db up -d postgres

# 停止
docker compose -f docker/docker-compose.yml down
```

| 服务 | 地址 |
|------|------|
| 前端 | http://localhost:8005 |
| 后端 | http://localhost:8006 |
| Postgres（local-db） | localhost:5432 |

## 数据库依赖

`backend` 服务默认不在 compose 里内置数据库，而是通过 `docker/.env` 里的
`DB_HOST:DB_PORT`（默认 `host.docker.internal:5432`）连接宿主机上已经在跑的
`postgres-dev` 容器（该容器需要用 `-p 5432:5432` 把端口发布到宿主机）。
部署前请确认：

- 已执行 `cp docker/.env.example docker/.env` 并按需修改；
- `postgres-dev`（或同等的 Postgres 实例）已启动，且端口已发布到宿主机 `5432`；
- 用户名 / 密码 / 库名与 `docker/.env` 一致
  （默认 `hsop` / `P0stgre2022` / `zw-note`，与 `zw-note-backend/scripts/db.sh` 默认值保持一致）；
- 首次部署需要建库 + 迁移 + 种子数据，可用 `make db-setup-docker`
  （默认操作名为 `postgres-dev` 的容器，可用 `CONTAINER=xxx` 覆盖）。

**生产环境注意**：`host.docker.internal` 只在 Docker Desktop（Mac/Windows）及新版
Docker Engine（配合 `extra_hosts: host-gateway`）下可用，指向的是运行 Docker 的
那台宿主机本身。若生产环境的 Postgres 是：

- 部署在**同一台服务器**上、以普通进程（非容器）运行：`host.docker.internal` 仍可用；
- 部署在**另一台服务器/云数据库**：应将 `docker/.env` 里的 `DB_HOST` 改成该数据库的
  实际 IP 或域名，`extra_hosts` 那行可以保留（不影响，只是不会被用到）；
- 也用 compose 起在同一网络里（即使用本文件的 `local-db` profile）：
  应将 `DB_HOST` 改成服务名 `postgres`，让 backend 通过 Docker 内部 DNS 访问。

## CI/CD 自动部署（`.github/workflows/deploy.yml`）

push 到 `main` 分支时会自动：

1. 构建 backend / frontend 镜像并推送到 Docker Hub
   （`wudande/zw-note-backend`、`wudande/zw-note-frontend`）；
2. 把仓库里的 `docker/docker-compose.yml` 直接 SCP 传到服务器（每次覆盖更新）；
3. SSH 登录服务器，`docker compose pull && docker compose up -d` 完成滚动更新。

**服务器上不需要 `git clone` 整个仓库**，只需要一个普通目录，里面放两个文件：
`docker-compose.yml`（由 CI 每次自动传输更新）和 `.env`（手动维护，CI 不会碰它）。

服务器端首次部署前需提前准备好：

- 建好一个目标目录（如 `/root/wudan/zw-note-deploy`），在该目录下执行
  一次 `.env` 初始化（可以先从本地 `docker/.env.example` 复制内容手动创建，
  或临时用 `scp` 传一份上去）并按需修改数据库连接参数；
- `wudande/zw-note-backend`、`wudande/zw-note-frontend` 保持 **public**，
  部署脚本里没有 `docker login`，直接 `docker compose pull` 拉公开镜像；
  若改成 private，需要在部署脚本里加回 `docker login`（用
  `DOCKERHUB_USERNAME` / `DOCKERHUB_TOKEN`）。

需要在仓库 **Settings → Secrets and variables → Actions** 里配置：

| Secret               | 说明                                        |
|----------------------|---------------------------------------------|
| `DOCKERHUB_USERNAME` | Docker Hub 用户名（`wudande`），构建推送镜像时登录用 |
| `DOCKERHUB_TOKEN`    | Docker Hub Access Token（Account Settings → Security → New Access Token，权限选 Read & Write） |
| `SSH_HOST`           | 服务器地址                                  |
| `SSH_USER`           | SSH 登录用户名                              |
| `SSH_PRIVATE_KEY`    | SSH 私钥（PEM 全文）                        |
| `SSH_PORT`           | SSH 端口，可不填（默认 22）                 |
| `DEPLOY_PATH`        | 服务器上一个普通目录的绝对路径（如 `/root/wudan/zw-note-deploy`），不需要是 git 仓库 |

`docker-compose.yml` 里 backend / frontend 服务同时保留了 `image:` 和
`build:`：本机 `up -d --build` 走本地构建，服务器上不带 `--build` 的
`pull && up -d` 则直接使用 CI 推送的镜像。

## 后端 Makefile 快捷方式

在 `zw-note-backend/` 下仍可使用：

```bash
make docker-build
make docker-run
make pg-up
make pg-down
```
