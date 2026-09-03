# ArgoCD Applications

This directory intentionally contains no ArgoCD Application manifests.

The ArgoCD app-of-apps is a **single application**: `ecommerce-production`
(defined in `../production/production.yaml`), which points to
`deployments/kubernetes/overlays/production`. That overlay renders the entire
base (`deployments/kubernetes/base`) with the production image overrides
(`newTag` pinned by CI) and the `GHCR_OWNER` templating.

Previously this directory held one Application per service/infra component,
each pointing at a `deployments/kubernetes/base/<dir>` (which now renders
un-pullable `__GHCR_OWNER__` placeholders until transformed by the production
overlay). Those per-service Applications were orphans (not referenced by any
kustomization) and have been removed — adding a new service requires no new
ArgoCD Application; it is picked up automatically via the base/overlay.
