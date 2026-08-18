#!/usr/bin/env bash
set -Eeuo pipefail

# 真实启动 Compose backend，并在退出时清理依赖容器和匿名卷。
export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-runtime_a_flower_shop_ecommerce}"
compose=(docker compose --env-file .env.example)
cleanup() {
  "${compose[@]}" down --volumes --remove-orphans >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM
cleanup
"${compose[@]}" up --build --no-color backend
