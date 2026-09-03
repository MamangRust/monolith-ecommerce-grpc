#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
overlay="$root/deployments/kubernetes/overlays/production"
rendered="$(mktemp)"
trap 'rm -f "$rendered"' EXIT

kubectl kustomize "$overlay" > "$rendered"

grep -q 'kind: ExternalSecret' "$rendered"
grep -q 'argocd.argoproj.io/hook: PreSync' "$rendered"
! grep -Eq 'image: .*:(latest|1\.0)([[:space:]]|$)' "$rendered"
! grep -Eq 'port: 0([[:space:]]|$)' "$rendered"
! grep -Eq '^[[:space:]]+(DB_PASSWORD|POSTGRES_PASSWORD|SECRET_KEY|REDIS_PASSWORD|SMTP_PASS):[[:space:]]+"[^$]' "$root/deployments/kubernetes/base/common/secrets.yaml"

echo "Production manifest validation passed"
