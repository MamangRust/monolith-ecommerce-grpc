#!/usr/bin/env bash
#
# load-test.sh — Dependency-free load test for the local e-commerce gateway.
#
# Usage:
#   ./load-test.sh [base_url] [requests] [concurrency]
#
# Runs `requests` HTTP calls against the gateway with `concurrency` workers
# (defaults: 200 requests, 10 workers). Uses only curl + bash so it works on
# any machine that can run the smoke test. Reports:
#   - total/ok/failed counts
#   - requests per second
#   - p50/p95/p99 latency (computed from recorded timings)
#
# The target endpoint is /health by default (no auth needed). Pass a second
# endpoint via ENDPOINT env var if you want to load something else, e.g.:
#   ENDPOINT=/api/auth/login ./load-test.sh
#
# Note: this is a coarse baseline, not a substitute for a real load-testing
# tool (k6/hey/wrk). It exists to catch gross regressions (e.g. a dead
# endpoint, connection leaks, unbounded latency) without adding tooling deps.

set -euo pipefail

BASE_URL="${1:-http://localhost:5000}"
REQUESTS="${2:-200}"
CONCURRENCY="${3:-10}"
ENDPOINT="${ENDPOINT:-/health}"

URL="${BASE_URL}${ENDPOINT}"

log() { printf '[%s] %s\n' "$(date +%H:%M:%S)" "$*"; }

log "Load test: ${REQUESTS} requests x ${CONCURRENCY} workers -> ${URL}"

START=$(date +%s%N)
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# Simple barrier: each worker writes timings to its own file.
worker() {
  local id=$1
  local per_file="${TMPDIR}/w${id}"
  : > "${per_file}"
  for ((i = 0; i < REQUESTS / CONCURRENCY; i++)); do
    local t0 t1 code
    t0=$(date +%s%N)
    code=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "$URL" || echo 000)
    t1=$(date +%s%N)
    echo "$((t1 - t0)) $code" >> "${per_file}"
  done
}

PIDS=()
for ((w = 0; w < CONCURRENCY; w++)); do
  worker "$w" &
  PIDS+=("$!")
done

for pid in "${PIDS[@]}"; do
  wait "$pid" || true
done

END=$(date +%s%N)
TOTAL_MS=$(( (END - START) / 1000000 ))

TIMINGS=$(cat "${TMPDIR}"/w* 2>/dev/null || true)
OK=$(echo "$TIMINGS" | awk '$2 == 200 || $2 == 201 {c++} END {print c+0}')
FAIL=$(echo "$TIMINGS" | awk '$2 != 200 && $2 != 201 {c++} END {print c+0}')
DONE=$(echo "$TIMINGS" | wc -l | tr -d ' ')

# Percentile helper (single pass bucket-free: sort once).
sort_timings() { echo "$TIMINGS" | awk '{print $1}' | sort -n; }

P50=$(sort_timings | awk -v n="$DONE" 'NR==int(n*0.5) || NR==int(n*0.5)+1 {s+=$1; c++} END {if (c) printf "%.1f", s/c/1000000; else print "0"}')
P95=$(sort_timings | awk -v n="$DONE" 'NR==int(n*0.95) || NR==int(n*0.95)+1 {s+=$1; c++} END {if (c) printf "%.1f", s/c/1000000; else print "0"}')
P99=$(sort_timings | awk -v n="$DONE" 'NR==int(n*0.99) || NR==int(n*0.99)+1 {s+=$1; c++} END {if (c) printf "%.1f", s/c/1000000; else print "0"}')

RPS=$(awk -v n="$DONE" -v ms="$TOTAL_MS" 'BEGIN {if (ms > 0) printf "%.1f", n / (ms / 1000); else print "0"}')

log "----"
log "Total: ${DONE} | OK: ${OK} | Failed: ${FAIL}"
log "Duration: ${TOTAL_MS} ms"
log "Throughput: ${RPS} req/s"
log "Latency: p50=${P50}ms p95=${P95}ms p99=${P99}ms"

if [ "$FAIL" -gt 0 ]; then
  log "RESULT: FAIL (${FAIL} non-2xx responses)"
  exit 1
fi

log "RESULT: PASS"
exit 0
