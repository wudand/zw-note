#!/usr/bin/env bash
# 数据库初始化脚本：建库 / 迁移 / 种子数据
#
# 用法:
#   ./scripts/db.sh local  <create|migrate|seed|setup|shell>
#   ./scripts/db.sh docker <create|migrate|seed|setup|shell>
#
# 环境变量可覆盖默认值:
#   PG_HOST PG_PORT PG_USER PG_PASSWORD PG_DB CONTAINER ADMIN_DB
#
# 示例:
#   ./scripts/db.sh docker setup
#   CONTAINER=postgres-dev ./scripts/db.sh docker migrate
#   ./scripts/db.sh local setup

set -euo pipefail

MODE="${1:-}"
ACTION="${2:-}"

PG_HOST="${PG_HOST:-127.0.0.1}"
PG_PORT="${PG_PORT:-5432}"
PG_USER="${PG_USER:-hsop}"
PG_PASSWORD="${PG_PASSWORD:-P0stgre2022}"
PG_DB="${PG_DB:-zw-note}"
CONTAINER="${CONTAINER:-postgres-dev}"
ADMIN_DB="${ADMIN_DB:-postgres}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="${SCRIPT_DIR}/migrations"
SEED_DIR="${SCRIPT_DIR}/seed"

usage() {
  cat <<EOF
Usage:
  $0 local  <create|migrate|seed|setup|shell>
  $0 docker <create|migrate|seed|setup|shell>

Commands:
  create   Create database (skip if already exists)
  migrate  Run scripts/migrations/*.sql
  seed     Run scripts/seed/*.sql
  setup    create + migrate + seed
  shell    Open interactive psql

Environment:
  PG_HOST       default: 127.0.0.1
  PG_PORT       default: 5432
  PG_USER       default: hsop
  PG_PASSWORD   default: P0stgre2022
  PG_DB         default: zw-note
  CONTAINER     default: postgres-dev  (docker mode only)
  ADMIN_DB      default: postgres      (used when creating DB)

Examples:
  $0 docker setup
  $0 local migrate
  CONTAINER=postgres-dev $0 docker shell
EOF
  exit 1
}

[[ -n "$MODE" && -n "$ACTION" ]] || usage
[[ "$MODE" == "local" || "$MODE" == "docker" ]] || usage

export PGPASSWORD="$PG_PASSWORD"

check_deps() {
  if [[ "$MODE" == "local" ]]; then
    command -v psql >/dev/null 2>&1 || {
      echo "error: psql not found; install PostgreSQL client or use docker mode" >&2
      exit 1
    }
  else
    command -v docker >/dev/null 2>&1 || {
      echo "error: docker not found" >&2
      exit 1
    }
    docker inspect "$CONTAINER" >/dev/null 2>&1 || {
      echo "error: container \"$CONTAINER\" not found or not running" >&2
      exit 1
    }
  fi
}

# 在指定库上执行一条 SQL（-c）
run_sql_c() {
  local db="$1"
  local sql="$2"
  if [[ "$MODE" == "local" ]]; then
    psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$db" \
      -v ON_ERROR_STOP=1 -c "$sql"
  else
    docker exec -i "$CONTAINER" \
      psql -U "$PG_USER" -d "$db" -v ON_ERROR_STOP=1 -c "$sql"
  fi
}

# 在指定库上执行 SQL 文件
run_sql_file() {
  local db="$1"
  local file="$2"
  echo ">> $file"
  if [[ "$MODE" == "local" ]]; then
    psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$db" \
      -v ON_ERROR_STOP=1 -f "$file"
  else
    docker exec -i "$CONTAINER" \
      psql -U "$PG_USER" -d "$db" -v ON_ERROR_STOP=1 < "$file"
  fi
}

# 查询单列结果（用于存在性判断）
query_scalar() {
  local db="$1"
  local sql="$2"
  if [[ "$MODE" == "local" ]]; then
    psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$db" \
      -tAc "$sql"
  else
    docker exec -i "$CONTAINER" \
      psql -U "$PG_USER" -d "$db" -tAc "$sql"
  fi
}

cmd_create() {
  echo ">> Creating database \"$PG_DB\" (skip if exists)..."
  local exists
  exists="$(query_scalar "$ADMIN_DB" "SELECT 1 FROM pg_database WHERE datname = '$PG_DB'" || true)"
  if [[ "$exists" == "1" ]]; then
    echo "   already exists"
  else
    run_sql_c "$ADMIN_DB" "CREATE DATABASE \"$PG_DB\";"
    echo "   created"
  fi
  echo ">> Database ready"
}

cmd_migrate() {
  echo ">> Running migrations against \"$PG_DB\"..."
  local f
  shopt -s nullglob
  local files=("$MIGRATIONS_DIR"/*.sql)
  shopt -u nullglob
  if [[ ${#files[@]} -eq 0 ]]; then
    echo "error: no migration files in $MIGRATIONS_DIR" >&2
    exit 1
  fi
  for f in "${files[@]}"; do
    run_sql_file "$PG_DB" "$f"
  done
  echo ">> Migrate done"
}

cmd_seed() {
  echo ">> Seeding \"$PG_DB\"..."
  local f
  shopt -s nullglob
  local files=("$SEED_DIR"/*.sql)
  shopt -u nullglob
  if [[ ${#files[@]} -eq 0 ]]; then
    echo "error: no seed files in $SEED_DIR" >&2
    exit 1
  fi
  for f in "${files[@]}"; do
    run_sql_file "$PG_DB" "$f"
  done
  echo ">> Seed done"
}

cmd_setup() {
  cmd_create
  cmd_migrate
  cmd_seed
}

cmd_shell() {
  if [[ "$MODE" == "local" ]]; then
    psql -h "$PG_HOST" -p "$PG_PORT" -U "$PG_USER" -d "$PG_DB"
  else
    docker exec -it "$CONTAINER" psql -U "$PG_USER" -d "$PG_DB"
  fi
}

check_deps

case "$ACTION" in
  create)  cmd_create ;;
  migrate) cmd_migrate ;;
  seed)    cmd_seed ;;
  setup)   cmd_setup ;;
  shell)   cmd_shell ;;
  *)       usage ;;
esac
