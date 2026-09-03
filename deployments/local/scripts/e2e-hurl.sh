#!/usr/bin/env bash
#
# e2e-hurl.sh — run every E2E hurl suite in tests/hurl/ against the gateway.
#
# Prerequisites:
#   - infra stack running:  docker compose -f deployments/local/docker-compose.infra.yml up -d
#   - Go services running locally (via `just services-local-start`)
#   - PostgreSQL has seeded data (run `just seeder-local` at least once) so
#     rules_strict.hurl can pick a seeded role-less user.
#
# Usage:
#   bash deployments/local/scripts/e2e-hurl.sh [base_url]
#
# Exits non-zero if any suite fails.

set -uo pipefail

BASE_URL="${1:-http://localhost:5000}"
RESET_DB="${2:-yes}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
HURL_DIR="$PROJECT_ROOT/tests/hurl"
TEST_IMAGE="$HURL_DIR/assets/test.png"

TS=$(date +%s)
TEST_PASSWORD="HurlPass123"

PASS=0
FAIL=0
FAILED=()

echo "== e2e-hurl: $BASE_URL =="
echo "hurl: $(hurl --version 2>/dev/null | head -1)"

# ---------------------------------------------------------------------------
# Reset the database so every run starts from a deterministic state. The
# all-endpoints sweep permanently deletes/trashes seeded rows (e.g. role 1),
# so repeated runs would otherwise 5xx on already-deleted records.
# ---------------------------------------------------------------------------
if [ "$RESET_DB" = "yes" ]; then
  echo "Resetting database (drop schema + migrate up + seed)..."
  if ! docker exec postgres_ecommerce psql -U DRAGON -d ECOMMERCE -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;" > /tmp/hurl_db_drop.log 2>&1; then
    echo "WARN: drop schema failed (see /tmp/hurl_db_drop.log); continuing with existing data"
  fi
  if ! (cd "$PROJECT_ROOT" && go run service/migrate/cmd/main.go -dir service/migrate/migrations up > /tmp/hurl_migrate_up.log 2>&1); then
    echo "ERROR: migrate up failed (see /tmp/hurl_migrate_up.log)"
    exit 1
  fi
  # register auto-assigns ROLE_ADMIN; the seeder only creates Cashier/Manager/Admin/Supplier.
  docker exec postgres_ecommerce psql -U DRAGON -d ECOMMERCE -c \
    "INSERT INTO roles (role_name) VALUES ('ROLE_ADMIN'), ('ROLE_USER') ON CONFLICT DO NOTHING;" >/dev/null 2>&1 || true
  (cd "$PROJECT_ROOT" && go run service/seeder/cmd/main.go > /tmp/hurl_seeder.log 2>&1) || \
    echo "WARN: seeder reported errors (see /tmp/hurl_seeder.log); continuing"
fi

# ---------------------------------------------------------------------------
# Discover a seeded role-less user for rules_strict.hurl (firstname User1,
# seeded with password "password1"). Register assigns ROLE_ADMIN to every new
# user, so a role-less seeded user is required to prove the 403 path.
# ---------------------------------------------------------------------------
USER_EMAIL=""
if docker exec postgres_ecommerce psql -U DRAGON -d ECOMMERCE -t -A -c \
  "SELECT email FROM users WHERE firstname='User1' AND email LIKE 'user_%' ORDER BY user_id LIMIT 1;" \
  > /tmp/hurl_seed_user.txt 2>/dev/null; then
  USER_EMAIL=$(tr -d ' \r' < /tmp/hurl_seed_user.txt)
  # Guarantee the chosen user is role-less so the 401 denial path is deterministic
  # (register auto-assigns ROLE_ADMIN, and the role seeder may assign random roles).
  docker exec postgres_ecommerce psql -U DRAGON -d ECOMMERCE -c \
    "DELETE FROM user_roles WHERE user_id = (SELECT user_id FROM users WHERE email = '$USER_EMAIL');" >/dev/null 2>&1 || true
fi

# ---------------------------------------------------------------------------
# Run every suite. Each file gets its own unique testEmail so register is
# idempotent across re-runs.
# ---------------------------------------------------------------------------
for f in "$HURL_DIR"/*.hurl; do
  name=$(basename "$f")
  testEmail="hurl.${name%.hurl}.$TS@example.com"

  common_vars="--test --jobs 1 --variable baseUrl=$BASE_URL --variable testEmail=$testEmail --variable testPassword=$TEST_PASSWORD --variable testImage=$TEST_IMAGE"

  if [ "$name" = "rules_strict.hurl" ]; then
    if [ -z "$USER_EMAIL" ]; then
      echo "SKIP $name (no seeded role-less user found; run seeder first)"
      continue
    fi
    extra_vars="--variable adminEmail=admin.$TS@example.com --variable adminPassword=$TEST_PASSWORD --variable userEmail=$USER_EMAIL --variable userPassword=password1"
    hurl $common_vars $extra_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  else
    hurl $common_vars "$f" > "/tmp/hurl_$name.log" 2>&1
  fi

  if [ $? -eq 0 ]; then
    PASS=$((PASS + 1))
    echo "PASS $name"
  else
    FAIL=$((FAIL + 1))
    FAILED+=("$name")
    echo "FAIL $name (see /tmp/hurl_$name.log)"
  fi

  # Pace the auth endpoints: the gateway rate-limits register/login with a
  # 10 req/s token bucket shared across all suites, so pause between files.
  sleep 2
done

echo
echo "== RESULT =="
echo "PASS: $PASS   FAIL: $FAIL   TOTAL: $((PASS + FAIL))"

if [ "$FAIL" -gt 0 ]; then
  echo
  echo "Failed suites:"
  for f in "${FAILED[@]}"; do echo "  - $f"; done
  exit 1
fi

echo "All E2E hurl suites passed."
exit 0
