# ============================================================
# zw-note 根目录 Makefile — 本地前后端启停 / 日志 / 状态
# ============================================================
ROOT        := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
BACKEND_DIR := $(ROOT)/zw-note-backend
FRONTEND_DIR:= $(ROOT)/zw-note-frontend
RUN_DIR     := $(ROOT)/.run
COMPOSE     := docker compose -f $(ROOT)/docker/docker-compose.yml

BACKEND_PORT  ?= 8004
FRONTEND_PORT ?= 5555
SWAGGER_URL   := http://localhost:$(BACKEND_PORT)/swagger/index.html
FRONTEND_URL  := http://localhost:$(FRONTEND_PORT)
BACKEND_URL   := http://localhost:$(BACKEND_PORT)
HEALTH_URL    := $(BACKEND_URL)/health

.PHONY: help start stop restart logs status \
        start-backend start-frontend stop-backend stop-frontend \
        docker-up docker-down docker-logs docker-status \
        db-setup-docker db-setup-local

# ------------------------------------------------------------
# 默认：帮助
# ------------------------------------------------------------
help:
	@echo ""
	@echo "  zw-note 常用命令（在仓库根目录执行）"
	@echo ""
	@echo "  本地开发"
	@echo "    make start      启动后端 + 前端（后台）"
	@echo "    make stop       停止后端 + 前端"
	@echo "    make restart    重启"
	@echo "    make logs       跟踪前后端日志（Ctrl+C 退出）"
	@echo "    make logs S=backend|frontend  只看一侧日志"
	@echo "    make status     查看启动状态与 Swagger 地址"
	@echo ""
	@echo "  Docker 部署（docker/docker-compose.yml）"
	@echo "    make docker-up / docker-down / docker-logs / docker-status"
	@echo ""
	@echo "  数据库（调用 zw-note-backend/scripts/db.sh）"
	@echo "    make db-setup-docker   容器 postgres-dev 建库+迁移+种子"
	@echo "    make db-setup-local    本机 psql 建库+迁移+种子"
	@echo ""

# ------------------------------------------------------------
# 本地启停
# ------------------------------------------------------------
start: start-backend start-frontend
	@echo ""
	@$(MAKE) --no-print-directory status

stop: stop-frontend stop-backend
	@echo ">> stopped"

restart: stop start

start-backend:
	@mkdir -p "$(RUN_DIR)"
	@if $(call port_listening,$(BACKEND_PORT)); then \
		echo ">> backend already running on :$(BACKEND_PORT)"; \
	else \
		echo ">> starting backend on :$(BACKEND_PORT)..."; \
		cd "$(BACKEND_DIR)" && \
		  nohup go run ./cmd/server/main.go start --config configs/config.yaml \
		  > "$(RUN_DIR)/backend.log" 2>&1 & echo $$! > "$(RUN_DIR)/backend.pid"; \
		$(call wait_http,$(HEALTH_URL),30) && echo ">> backend ready" || \
		  (echo ">> backend failed to become healthy; see $(RUN_DIR)/backend.log" && exit 1); \
	fi

start-frontend:
	@mkdir -p "$(RUN_DIR)"
	@if $(call port_listening,$(FRONTEND_PORT)); then \
		echo ">> frontend already running on :$(FRONTEND_PORT)"; \
	else \
		echo ">> starting frontend on :$(FRONTEND_PORT)..."; \
		cd "$(FRONTEND_DIR)" && \
		  nohup npm run dev \
		  > "$(RUN_DIR)/frontend.log" 2>&1 & echo $$! > "$(RUN_DIR)/frontend.pid"; \
		$(call wait_port,$(FRONTEND_PORT),60) && echo ">> frontend ready" || \
		  (echo ">> frontend failed to start; see $(RUN_DIR)/frontend.log" && exit 1); \
	fi

stop-backend:
	@echo ">> stopping backend..."
	@$(call kill_port,$(BACKEND_PORT))
	@rm -f "$(RUN_DIR)/backend.pid"

stop-frontend:
	@echo ">> stopping frontend..."
	@$(call kill_port,$(FRONTEND_PORT))
	@rm -f "$(RUN_DIR)/frontend.pid"

# ------------------------------------------------------------
# 日志
# ------------------------------------------------------------
# S=backend|frontend|all（默认 all）
S ?= all

logs:
	@mkdir -p "$(RUN_DIR)"
	@touch "$(RUN_DIR)/backend.log" "$(RUN_DIR)/frontend.log"
	@case "$(S)" in \
	  backend)  echo ">> tail $(RUN_DIR)/backend.log";  tail -n 100 -f "$(RUN_DIR)/backend.log" ;; \
	  frontend) echo ">> tail $(RUN_DIR)/frontend.log"; tail -n 100 -f "$(RUN_DIR)/frontend.log" ;; \
	  all|*)    echo ">> tail $(RUN_DIR)/backend.log + frontend.log"; \
	            tail -n 50 -f "$(RUN_DIR)/backend.log" "$(RUN_DIR)/frontend.log" ;; \
	esac

# ------------------------------------------------------------
# 状态
# ------------------------------------------------------------
status:
	@echo ""
	@echo "========== zw-note status =========="
	@printf "  后端  (:$(BACKEND_PORT))  "; \
	  if $(call port_listening,$(BACKEND_PORT)); then echo "已启动"; else echo "未启动"; fi
	@printf "  前端  (:$(FRONTEND_PORT))  "; \
	  if $(call port_listening,$(FRONTEND_PORT)); then echo "已启动"; else echo "未启动"; fi
	@echo ""
	@echo "  前端地址 : $(FRONTEND_URL)"
	@echo "  后端地址 : $(BACKEND_URL)"
	@echo "  Swagger  : $(SWAGGER_URL)"
	@echo "===================================="
	@echo ""

# ------------------------------------------------------------
# Docker 部署快捷方式
# ------------------------------------------------------------
docker-up:
	$(COMPOSE) up -d --build

docker-down:
	$(COMPOSE) down

docker-logs:
	$(COMPOSE) logs -f

docker-status:
	@echo ""
	@echo "========== docker status =========="
	@$(COMPOSE) ps || true
	@echo ""
	@echo "  Frontend : http://localhost:8005"
	@echo "  Backend  : http://localhost:8006"
	@echo "  Swagger  : http://localhost:8006/swagger/index.html"
	@echo "==================================="
	@echo ""

# ------------------------------------------------------------
# 数据库
# ------------------------------------------------------------
db-setup-docker:
	@$(BACKEND_DIR)/scripts/db.sh docker setup

db-setup-local:
	@$(BACKEND_DIR)/scripts/db.sh local setup

# ------------------------------------------------------------
# 内部工具（shell snippets）
# ------------------------------------------------------------
# 端口是否在 LISTEN
define port_listening
lsof -nP -iTCP:$(1) -sTCP:LISTEN >/dev/null 2>&1
endef

# 杀掉占用端口的进程
define kill_port
(lsof -nP -tiTCP:$(1) -sTCP:LISTEN | xargs kill 2>/dev/null) || true; \
sleep 0.3; \
(lsof -nP -tiTCP:$(1) -sTCP:LISTEN | xargs kill -9 2>/dev/null) || true
endef

# 等待端口就绪：wait_port,port,seconds
define wait_port
i=0; \
while [ $$i -lt $(2) ]; do \
  if lsof -nP -iTCP:$(1) -sTCP:LISTEN >/dev/null 2>&1; then exit 0; fi; \
  i=$$((i+1)); sleep 1; \
done; exit 1
endef

# 等待 HTTP 200：wait_http,url,seconds
define wait_http
i=0; \
while [ $$i -lt $(2) ]; do \
  if curl -sf "$(1)" >/dev/null 2>&1; then exit 0; fi; \
  i=$$((i+1)); sleep 1; \
done; exit 1
endef
