#!/usr/bin/env bash
#
# services-local-stop.sh — stop every locally running Go service.
#
# Usage:
#   bash deployments/local/scripts/services-local-stop.sh

set -uo pipefail

# Kill the go run wrappers and the compiled service binaries listening on the
# service gRPC ports (the binaries live in /tmp/go-build*/b001/exe/main).
pkill -f 'go run service/' 2>/dev/null || true

for p in 50051 50052 50053 50054 50055 50056 50057 50058 50059 50060 50061 50062 50063 50064 50065 50066 50067 50068 50069 5000; do
  pid=$(ss -ltnp 2>/dev/null | grep ":$p " | sed -n 's/.*pid=\([0-9]*\).*/\1/p' | head -1)
  if [ -n "$pid" ]; then
    kill "$pid" 2>/dev/null || true
    echo "stopped service on :$p (pid $pid)"
  fi
done

echo "All local services stopped."
