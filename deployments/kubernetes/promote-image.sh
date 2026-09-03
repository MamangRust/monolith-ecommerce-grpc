#!/usr/bin/env bash
set -euo pipefail

sha="${1:-${GITHUB_SHA:-}}"
if [[ ! "$sha" =~ ^[0-9a-f]{40}$ ]]; then
  echo "usage: $0 <40-character git SHA>" >&2
  exit 2
fi

file="deployments/kubernetes/overlays/production/kustomization.yaml"
sed -i -E "s/newTag: [0-9a-f]{40}/newTag: ${sha}/g" "$file"

echo "Pinned production images to ${sha} in ${file}"
