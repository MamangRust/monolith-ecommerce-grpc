# Production Deployment Runbook

## Scope

Production is deployed only through:

```text
deployments/gitops/argocd/production-root.yaml
  -> deployments/gitops/argocd/production/kustomization.yaml
  -> deployments/gitops/argocd/production/production.yaml
  -> deployments/kubernetes/overlays/production
```

Do not enable the legacy per-service ArgoCD Applications at the same time as
`ecommerce-production`; they target overlapping resources and can fight over
ownership.

## Prerequisites

- Kubernetes context points to the intended production cluster.
- External Secrets Operator is installed.
- `ClusterSecretStore/ecommerce-secrets` is Ready.
- `Secret/app-secrets` and `Secret/ghcr-pull-secret` can be materialized.
- Metrics Server is installed for HPA.
- GHCR package access is valid for all 22 application images.
- PostgreSQL, Kafka, Redis, and storage classes are healthy.

Verify before sync:

```bash
kubectl config current-context
kubectl -n external-secrets get pods
kubectl get clustersecretstore ecommerce-secrets
kubectl -n ecommerce get secret app-secrets ghcr-pull-secret
kubectl api-resources | grep -E 'externalsecrets|clustersecretstores'
```

## Promote a release

Use the exact SHA produced by CI. Never use `latest` in production:

```bash
./deployments/kubernetes/promote-image.sh "$GITHUB_SHA"
./deployments/kubernetes/validate-production.sh
kubectl kustomize deployments/kubernetes/overlays/production >/tmp/ecommerce-production.yaml
```

Review the rendered diff, then commit the promoted SHA. CI must publish all
22 images for that SHA before ArgoCD is synced.

## Server-side validation and sync

```bash
kubectl apply --server-side --dry-run=server \
  -k deployments/kubernetes/overlays/production
kubectl diff -k deployments/kubernetes/overlays/production
kubectl apply -f deployments/gitops/argocd/production-root.yaml
```

ArgoCD ordering:

1. Infrastructure/common resources.
2. `migrate` PreSync hook, recreated for the promoted image.
3. Application Deployments at sync wave 2.
4. Observability and edge resources.

The migration hook must succeed before the application rollout is considered
healthy. Do not manually delete migration Jobs during a sync.

## Post-deploy verification

```bash
kubectl -n ecommerce get applications -o wide
kubectl -n ecommerce get pods -o wide
kubectl -n ecommerce get events --sort-by=.lastTimestamp | tail -50
kubectl -n ecommerce get hpa,pdb
kubectl -n ecommerce rollout status deployment/apigateway --timeout=10m
kubectl -n ecommerce get endpoints apigateway
kubectl -n ecommerce port-forward svc/nginx 8080:80
curl -fsS http://127.0.0.1:8080/health
```

Check metrics-server, target health, Kafka consumer lag, PostgreSQL exporter,
OTel/Jaeger traces, and alert rules before declaring the release healthy.

## Rollback

Rollback is a Git revert or a new promotion to the previous immutable SHA:

```bash
./deployments/kubernetes/promote-image.sh "$PREVIOUS_SHA"
git commit -am "rollback production to $PREVIOUS_SHA"
```

Then allow ArgoCD to reconcile and verify:

```bash
kubectl -n ecommerce rollout status deployment/apigateway --timeout=10m
kubectl -n ecommerce get pods
```

If the migration is not backward-compatible, stop and restore the database
from the approved backup before rolling application code back.

## Backup and restore gate

Before production go-live, test a PostgreSQL logical backup and restore into a
separate database/cluster. Record the backup timestamp, migration version,
restore duration, and verification query. A successful manifest sync alone is
not proof of disaster-recovery readiness.
