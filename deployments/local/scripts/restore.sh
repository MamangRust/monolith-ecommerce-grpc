#!/usr/bin/env bash
#
# restore.sh — PostgreSQL restore for the local e-commerce stack.
#
# Usage:
#   ./scripts/restore.sh <backup_file>
#
# The backup file may be a raw .sql dump or a .sql.gz dump produced by
# backup.sh. Restoring drops and recreates the target database, so this is
# destructive — a confirmation prompt is shown unless --yes is passed.
#
#   ./scripts/restore.sh --yes backups/ecommerce_20260101_000000.sql.gz
#
# Procedure (Fase 6 checklist 11.3):
#   1. Validate the backup file exists.
#   2. Confirm intent (destructive).
#   3. Drop + recreate the database (best-effort, ignores active connections).
#   4. Restore from the dump.
#   5. Report success/failure.

set -euo pipefail

COMPOSE_FILE="$(cd "$(dirname "$0")/.." && pwd)/docker-compose.yml"
DB_USER="${POSTGRES_USER:-DRAGON}"
DB_NAME="${POSTGRES_DB:-ECOMMERCE}"

AUTO_YES=0
BACKUP_FILE=""

for arg in "$@"; do
  case "$arg" in
    --yes) AUTO_YES=1 ;;
    *)     BACKUP_FILE="$arg" ;;
  esac
done

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: $0 [--yes] <backup_file>" >&2
  echo "  backup_file: path to a .sql or .sql.gz dump" >&2
  exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
  echo "ERROR: backup file not found: ${BACKUP_FILE}" >&2
  exit 1
fi

# Validate the dump BEFORE the destructive drop: a truncated/empty/corrupt file
# would wipe the database and restore nothing.
if [ ! -s "$BACKUP_FILE" ] || [ "$(stat -c%s "$BACKUP_FILE")" -lt 1024 ]; then
  echo "ERROR: backup file is empty or too small (< 1 KB); refusing to restore." >&2
  exit 1
fi

if [[ "$BACKUP_FILE" == *.gz ]]; then
  if ! gzip -t "$BACKUP_FILE" 2>/dev/null; then
    echo "ERROR: backup file is not a valid gzip archive; refusing to restore." >&2
    exit 1
  fi
fi

if ! docker compose -f "${COMPOSE_FILE}" ps postgres >/dev/null 2>&1; then
  echo "ERROR: postgres container is not present in compose." >&2
  exit 1
fi

if [ "$AUTO_YES" -ne 1 ]; then
  read -r -p "Restore will DROP and recreate database '${DB_NAME}'. Continue? [y/N] " answer
  case "$answer" in
    y|Y) ;;
    *) echo "Aborted."; exit 1 ;;
  esac
fi

echo "Restoring ${BACKUP_FILE} -> ${DB_NAME}"

# Terminate existing connections and drop the database so the restore is clean.
docker compose -f "${COMPOSE_FILE}" exec -T postgres \
  psql -U "${DB_USER}" -d postgres -v ON_ERROR_STOP=1 \
  -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '${DB_NAME}' AND pid <> pg_backend_pid();" \
  -c "DROP DATABASE IF EXISTS ${DB_NAME};" \
  -c "CREATE DATABASE ${DB_NAME};"

case "$BACKUP_FILE" in
  *.gz)
    gunzip -c "$BACKUP_FILE" | docker compose -f "${COMPOSE_FILE}" exec -T postgres \
      psql -U "${DB_USER}" -d "${DB_NAME}" -v ON_ERROR_STOP=1
    ;;
  *)
    cat "$BACKUP_FILE" | docker compose -f "${COMPOSE_FILE}" exec -T postgres \
      psql -U "${DB_USER}" -d "${DB_NAME}" -v ON_ERROR_STOP=1
    ;;
esac

echo "Restore complete: ${DB_NAME} is back in service."
