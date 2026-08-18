#!/usr/bin/env bash
set -Eeuo pipefail

# 真实启动 Compose backend；探活通过后由 runtime_smoke 终止本脚本，trap 负责清理。
export COMPOSE_PROJECT_NAME="${COMPOSE_PROJECT_NAME:-runtime_a_flower_shop_ecommerce}"
export DB_PORT="${DB_PORT:-0}" REDIS_PORT="${REDIS_PORT:-0}" MINIO_PORT="${MINIO_PORT:-0}" MINIO_CONSOLE_PORT="${MINIO_CONSOLE_PORT:-0}" FRONTEND_PORT="${FRONTEND_PORT:-0}"
compose=(docker compose --env-file .env.example)
cleanup() {
  if [[ -n "${logs_pid:-}" ]]; then kill "$logs_pid" >/dev/null 2>&1 || true; fi
  "${compose[@]}" down --volumes --remove-orphans >/dev/null 2>&1 || true
}
trap 'cleanup; exit 143' TERM INT
trap cleanup EXIT
cleanup
"${compose[@]}" up --build --no-color --detach backend
"${compose[@]}" logs --follow --no-color backend &
logs_pid=$!
wait "$logs_pid"
