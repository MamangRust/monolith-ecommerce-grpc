#!/usr/bin/env bash
# endpoint-test.sh — jalankan SEMUA route API yang terdaftar di swagger.json.
#
# Kriteria:
#   - GET list (query)  : harus 200
#   - route lain        : tidak boleh 5xx (2xx/3xx/4xx diterima; 400/401/403/404
#                         adalah respon bisnis/auth yang sah)
#   - /swagger/doc.json : harus 200 (verifikasi spec ter-register)
#
# Pemakaian: bash endpoint-test.sh [base_url] [swagger.json]
set -uo pipefail

BASE="${1:-http://localhost:5000}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SWAGGER="${2:-$SCRIPT_DIR/../../../service/apigateway/docs/swagger.json}"

PASS=0
FAIL=0
FAILED=()

ok()   { PASS=$((PASS + 1)); }
bad()  { FAIL=$((FAIL + 1)); FAILED+=("$1"); echo "FAIL $1"; }

echo "== endpoint-test: $BASE =="

# ---------------------------------------------------------------------------
# 1. Auth setup: register + login -> token & user id
# ---------------------------------------------------------------------------
EMAIL="ept.$(date +%s)@test.com"
REG=$(curl -s --max-time 15 -X POST "$BASE/api/auth/register" \
  -H 'Content-Type: application/json' \
  -d "{\"firstname\":\"EPT\",\"lastname\":\"Test\",\"email\":\"$EMAIL\",\"password\":\"SmokePass123\",\"confirm_password\":\"SmokePass123\",\"is_verified\":true}")
LOGIN=$(curl -s --max-time 15 -X POST "$BASE/api/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\",\"password\":\"SmokePass123\"}")
TOKEN=$(printf '%s' "$LOGIN" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
if [ -z "$TOKEN" ]; then
  echo "AUTH FAILED: tidak bisa login (register=$(printf '%s' "$REG" | head -c 120))"
  exit 2
fi
USER_ID=$(printf '%s' "$LOGIN" | sed -n 's/.*"user_id":\([0-9]*\).*/\1/p')
[ -z "$USER_ID" ] && USER_ID=1
echo "auth ok: user_id=$USER_ID token_len=${#TOKEN}"

# ---------------------------------------------------------------------------
# 2. Enumerasi route dari swagger.json
# ---------------------------------------------------------------------------
if [ ! -f "$SWAGGER" ]; then
  echo "swagger.json tidak ditemukan: $SWAGGER"
  exit 2
fi
python3 - "$SWAGGER" > /tmp/ept_routes.tsv <<'PY'
import json, sys
d = json.load(open(sys.argv[1]))
for p, ms in sorted(d.get("paths", {}).items()):
    for m in ms:
        print(f"{m.upper()}\t{p}")
PY
TOTAL=$(wc -l < /tmp/ept_routes.tsv)
echo "routes terdaftar di swagger: $TOTAL"

# ---------------------------------------------------------------------------
# 3. Verifikasi swagger spec tersaji (regresi fix blank-import docs)
# ---------------------------------------------------------------------------
SC=$(curl -s -o /tmp/ept_swagger -w '%{http_code}' --max-time 10 "$BASE/swagger/doc.json")
if [ "$SC" = "200" ]; then ok; else bad "GET /swagger/doc.json -> $SC"; fi

# ---------------------------------------------------------------------------
# 4. Jalankan setiap route
# ---------------------------------------------------------------------------
is_strict_list() { # path list utama -> harus 200
  local p="$1"
  case "$p" in
    *-query|*-query/active|*-query/trashed) return 0 ;;
    /api/order-item|/api/order-item/active|/api/order-item/trashed) return 0 ;;
  esac
  return 1
}

while IFS=$'\t' read -r M P; do
  URL=$(printf '%s' "$P" | sed -e "s|{user_id}|$USER_ID|g" -e 's|{[^}]*}|1|g')
  # beri pagination default pada list query
  case "$URL" in
    *-query|*-query/active|*-query/trashed|/api/order-item|/api/order-item/active|/api/order-item/trashed)
      case "$URL" in *\?*) ;; *) URL="$URL?page=1&page_size=10" ;; esac ;;
  esac

  if [ "$M" = "GET" ]; then
    CODE=$(curl -s -o /tmp/ept_body -w '%{http_code}' --max-time 10 \
      -H "Authorization: Bearer $TOKEN" "$BASE$URL")
  else
    CODE=$(curl -s -o /tmp/ept_body -w '%{http_code}' --max-time 10 \
      -X "$M" -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
      -d '{}' "$BASE$URL")
  fi

  if is_strict_list "$P"; then
    if [ "$CODE" = "200" ]; then ok; else bad "$M $URL -> $CODE (strict 200)"; fi
  elif [[ "$CODE" =~ ^5 ]]; then
    bad "$M $URL -> $CODE  body=$(head -c 100 /tmp/ept_body | tr '\n' ' ')"
  else
    ok
  fi
done < /tmp/ept_routes.tsv

# ---------------------------------------------------------------------------
# 5. Ringkasan
# ---------------------------------------------------------------------------
echo
echo "== RESULT =="
echo "PASS: $PASS   FAIL: $FAIL   TOTAL: $((PASS + FAIL))"
if [ "$FAIL" -gt 0 ]; then
  echo
  echo "Daftar endpoint gagal (5xx / strict-200 gagal):"
  for f in "${FAILED[@]}"; do echo "  - $f"; done
  exit 1
fi
echo "Semua endpoint OK."
