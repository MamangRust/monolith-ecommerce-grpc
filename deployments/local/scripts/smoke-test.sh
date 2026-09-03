#!/usr/bin/env bash
#
# smoke-test.sh — Representative smoke test for the local e-commerce stack.
#
# Usage:
#   ./scripts/smoke-test.sh [base_url] [email]
#
# Checks, in order:
#   1. Liveness  : GET /health                     -> 200 "healthy"
#   2. Readiness : GET /ready                      -> 200 "ready" (Redis ping)
#   3. Register  : POST /api/auth/register         -> 201/200
#   4. Login     : POST /api/auth/login            -> 200 + access token
#   5. Me        : GET /api/auth/me   (JWT)        -> 200
#   6. Products  : GET /api/product-query (JWT)    -> 200 (representative read)
#
# Exits non-zero on the first failure. Safe to re-run after a restart —
# register/login use a unique email per run when none is provided.
#
# Design notes (Fase 6 checklist):
#   - Does NOT treat a "running" container as success; it probes real HTTP.
#   - Representative endpoints cover auth (command+query) and one read path.
#   - Idempotent: a leftover user from a previous run does not fail the script.

set -euo pipefail

BASE_URL="${1:-http://localhost:5000}"
EMAIL="${2:-smoke-$(date +%s)@example.com}"
PASSWORD="SmokePass123"

PASS=0
FAIL=0

log()  { printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*"; }
ok()   { log "PASS: $*"; PASS=$((PASS + 1)); }
fail() { log "FAIL: $*"; FAIL=$((FAIL + 1)); }

http_code() { # url [data]
  if [ -n "${2:-}" ]; then
    curl -s -o /dev/null -w '%{http_code}' -X POST -H 'Content-Type: application/json' -d "$2" "$1"
  else
    curl -s -o /dev/null -w '%{http_code}' "$1"
  fi
}

# Body check ensures the endpoint really is our gateway (returns JSON with
# status "healthy"/"ready") and not some other server squatting on the port.
# Grep is whitespace-tolerant ("status":"healthy" or "status": "healthy").
body_contains() { # url regex
  curl -s --max-time 5 "$1" | grep -Eq "$2"
}

log "Smoke test against ${BASE_URL} (email: ${EMAIL})"

HEALTHY_RE='"status"[[:space:]]*:[[:space:]]*"healthy"'
READY_RE='"status"[[:space:]]*:[[:space:]]*"ready"'

# --- 1. Liveness ----------------------------------------------------------
code=$(http_code "${BASE_URL}/health")
if [ "$code" = "200" ] && body_contains "${BASE_URL}/health" "$HEALTHY_RE"; then
  ok "/health -> $code (healthy)"
else
  fail "/health -> $code or body not 'healthy' (is this the e-commerce gateway?)"
fi

# --- 2. Readiness ---------------------------------------------------------
code=$(http_code "${BASE_URL}/ready")
if [ "$code" = "200" ] && body_contains "${BASE_URL}/ready" "$READY_RE"; then
  ok "/ready -> $code (ready)"
else
  fail "/ready -> $code or body not 'ready' (is this the e-commerce gateway?)"
fi

# --- 3. Register ----------------------------------------------------------
REGISTER_JSON=$(printf '{"firstname":"Smoke","lastname":"Test","email":"%s","password":"%s","confirm_password":"%s","is_verified":true}' "$EMAIL" "$PASSWORD" "$PASSWORD")
REGISTER_RESP=$(curl -s -X POST -H 'Content-Type: application/json' -d "$REGISTER_JSON" "${BASE_URL}/api/auth/register")
# Register returns 201 with {"status":"success",...}. Check the success
# marker in the body (not the HTTP code, which the gateway may map differently).
if echo "$REGISTER_RESP" | grep -q '"status"[[:space:]]*:[[:space:]]*"success"'; then
  ok "/api/auth/register -> success"
else
  fail "/api/auth/register -> body: $(echo "$REGISTER_RESP" | head -c 120)"
fi

# --- 4. Login -------------------------------------------------------------
LOGIN_JSON=$(printf '{"email":"%s","password":"%s"}' "$EMAIL" "$PASSWORD")
LOGIN_RESP=$(curl -s -X POST -H 'Content-Type: application/json' -d "$LOGIN_JSON" "${BASE_URL}/api/auth/login")
TOKEN=$(echo "$LOGIN_RESP" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p' | head -1)

# Login can be rate-limited/locked if the same email is reused after failures;
# retry with a fresh email once to keep the smoke test idempotent.
if [ -z "$TOKEN" ]; then
  EMAIL="smoke-$(date +%s)-retry@example.com"
  REGISTER_JSON=$(printf '{"firstname":"Smoke","lastname":"Test","email":"%s","password":"%s","confirm_password":"%s","is_verified":true}' "$EMAIL" "$PASSWORD" "$PASSWORD")
  curl -s -X POST -H 'Content-Type: application/json' -d "$REGISTER_JSON" "${BASE_URL}/api/auth/register" >/dev/null 2>&1 || true
  LOGIN_JSON=$(printf '{"email":"%s","password":"%s"}' "$EMAIL" "$PASSWORD")
  LOGIN_RESP=$(curl -s -X POST -H 'Content-Type: application/json' -d "$LOGIN_JSON" "${BASE_URL}/api/auth/login")
  TOKEN=$(echo "$LOGIN_RESP" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p' | head -1)
fi

if [ -n "$TOKEN" ]; then ok "/api/auth/login -> token acquired"; else fail "/api/auth/login -> no access_token in response"; fi

# --- 5. Me (JWT protected) ------------------------------------------------
if [ -n "$TOKEN" ]; then
  code=$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/auth/me")
  if [ "$code" = "200" ]; then ok "/api/auth/me -> $code"; else fail "/api/auth/me -> $code (want 200)"; fi
else
  fail "/api/auth/me skipped (no token)"
fi

# --- 6. Representative read path ------------------------------------------
if [ -n "$TOKEN" ]; then
  code=$(curl -s -o /dev/null -w '%{http_code}' -H "Authorization: Bearer ${TOKEN}" "${BASE_URL}/api/product-query")
  if [ "$code" = "200" ]; then ok "/api/product-query -> $code"; else fail "/api/product-query -> $code (want 200)"; fi
else
  fail "/api/product-query skipped (no token)"
fi

log "----"
log "Result: $PASS passed, $FAIL failed"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi

exit 0
