#!/usr/bin/env bash
#
# services-local-start.sh — start every Go service locally (go run) in the
# background, with logs written to deployments/local/logs/<service>.log so
# Promtail can ship them to Loki.
#
# Services are detached with setsid so they survive the parent shell.
#
# Usage:
#   bash deployments/local/scripts/services-local-start.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
LOG_DIR="$PROJECT_ROOT/deployments/local/logs"

mkdir -p "$LOG_DIR"

# Order matters for readability only; the API gateway is started last so all
# gRPC backends are (usually) already listening when it boots.
SERVICES=(
  auth role user category merchant merchant_award merchant_business merchant_detail merchant_policy
  order order_item product transaction cart review review_detail slider shipping_address banner email apigateway
)

for svc in "${SERVICES[@]}"; do
  setsid nohup bash -c "cd '$PROJECT_ROOT' && go run service/$svc/cmd/main.go > '$LOG_DIR/$svc.log' 2>&1" > /dev/null 2>&1 &
  echo "started $svc"
done

echo
echo "All $((${#SERVICES[@]})) services launched. Logs: $LOG_DIR"
