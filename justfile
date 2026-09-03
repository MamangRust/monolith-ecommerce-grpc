set shell := ["bash", "-c"]
PROTOBUF_INCLUDE := env_var_or_default("PROTOBUF_INCLUDE", "third_party/protobuf")

# Tidy all go.mod files
tidy-all:
    @for mod in service/*/go.mod; do \
        dir=$(dirname $mod); \
        service=$(basename $dir); \
        echo "🧹 Tidying $service..."; \
        (cd $dir && go mod tidy) || exit 1; \
    done
    @echo "✅ All services tidied successfully."

# Build all services that contain a go.mod file
build:
    @mkdir -p bin
    @for mod in service/*/go.mod; do \
        dir=$(dirname $mod); \
        service=$(basename $dir); \
        echo "🔨 Building $service..."; \
        (cd $dir && go build -o ../../bin/$service ./cmd/main.go) || exit 1; \
    done
    @echo "✅ All services built successfully in bin/ folder."

# Clean build artifacts
clean:
    rm -rf bin
    @echo "🧹 Cleaned bin/ folder."

# Migrate up
migrate:
    go run service/migrate/cmd/main.go up

# Migrate down
migrate-down:
    go run service/migrate/cmd/main.go down

# Generate Proto
generate-proto:
    @if [ -z "{{PROTOBUF_INCLUDE}}" ] || [ ! -f "{{PROTOBUF_INCLUDE}}/google/protobuf/wrappers.proto" ]; then \
        echo "Missing protobuf well-known types. Set PROTOBUF_INCLUDE to a directory containing google/protobuf/wrappers.proto." >&2; \
        exit 1; \
    fi
    protoc --proto_path=proto --proto_path="{{PROTOBUF_INCLUDE}}" --go_out=shared/pb --go_opt=module=github.com/MamangRust/monolith-ecommerce-shared/pb --go_opt=Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb --go_opt=Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb --go-grpc_out=shared/pb --go-grpc_opt=module=github.com/MamangRust/monolith-ecommerce-shared/pb --go-grpc_opt=Mgoogle/protobuf/wrappers.proto=google.golang.org/protobuf/types/known/wrapperspb --go-grpc_opt=Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb $(find proto -name "*.proto" -print)

# Generate SQLC output
generate-sql:
    sqlc generate

# Generate Swagger
generate-swagger:
    swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs

# Run Seeder
seeder:
    go run service/seeder/cmd/main.go

# Build Docker images for all services
build-image:
    @for service in apigateway migrate auth user role banner cart category email merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product review review_detail shipping_address slider transaction seeder; do \
        echo "🔨 Building $service-ecommerce-service..."; \
        docker build -t $service-ecommerce-service:1.0 -f service/$service/Dockerfile . || exit 1; \
    done
    @echo "✅ All service images built successfully."

# Docker Compose up
up:
    docker compose -f deployments/local/docker-compose.yml up -d

# Docker Compose down
down:
    docker compose -f deployments/local/docker-compose.yml down

# Infra-only compose up (postgres, redis, kafka, observability) — for running Go services locally
infra-up:
    docker compose -f deployments/local/docker-compose.infra.yml up -d

# Infra-only compose down
infra-down:
    docker compose -f deployments/local/docker-compose.infra.yml down

# Run migrations against local PostgreSQL
db-migrate:
    go run service/migrate/cmd/main.go -dir service/migrate/migrations up

# Seed local PostgreSQL with sample data
db-seeder:
    go run service/seeder/cmd/main.go

# Start all Go services locally (background, logs under deployments/local/logs)
services-local-start:
    @bash deployments/local/scripts/services-local-start.sh

# Stop all locally running Go services
services-local-stop:
    @bash deployments/local/scripts/services-local-stop.sh

# Run every E2E hurl suite against the running gateway (optional BASE_URL)
e2e-hurl base_url = "http://localhost:5000":
    @bash deployments/local/scripts/e2e-hurl.sh "{{base_url}}"

# Build images and start compose
build-up: build-image up

# Validate docker-compose config
compose-config:
    docker compose -f deployments/local/docker-compose.yml config --quiet && echo "✅ Compose config valid"

# Run smoke test against the local gateway (optional BASE_URL + EMAIL args)
smoke-test base_url = "http://localhost:5000":
    @bash deployments/local/scripts/smoke-test.sh "{{base_url}}"

# Run a dependency-free load test (optional BASE_URL, REQUESTS, CONCURRENCY)
load-test base_url = "http://localhost:5000" requests = "100" concurrency = "10":
    @bash deployments/local/scripts/load-test.sh "{{base_url}}" "{{requests}}" "{{concurrency}}"

# Run every route in swagger.json against the running gateway (no-5xx + strict-200 for lists)
endpoint-test base_url = "http://localhost:5000":
    @bash deployments/local/scripts/endpoint-test.sh "{{base_url}}"

# Backup PostgreSQL to deployments/local/backups (optional retention days)
backup retention = "7":
    @bash deployments/local/scripts/backup.sh "{{retention}}"

# Restore PostgreSQL from a backup file (destructive; use --yes to skip prompt)
restore file:
    @bash deployments/local/scripts/restore.sh "{{file}}"

# Show migration status
migrate-status:
    go run service/migrate/cmd/main.go status

# Rollback one migration version
migrate-rollback:
    go run service/migrate/cmd/main.go down

# Tail local service logs (optional service name glob)
logs svc = "*":
    @ls deployments/local/logs/{{svc}}.log 2>/dev/null | while read f; do echo "--- $$(basename $$f)"; tail -n 20 "$$f"; done

# Kubernetes: render all manifests via kustomize
k8s-render:
    @if [ ! -f deployments/kubernetes/secrets.yaml ]; then cp deployments/kubernetes/secrets.example.yaml deployments/kubernetes/secrets.yaml; echo "ℹ️ Created deployments/kubernetes/secrets.yaml from example (fill real values before real apply)"; fi
    @kubectl kustomize deployments/kubernetes > /tmp/ecommerce-k8s.yaml
    @echo "✅ Rendered $(grep -c '^kind:' /tmp/ecommerce-k8s.yaml) resources to /tmp/ecommerce-k8s.yaml"

# Kubernetes: validate manifests (client dry-run, no cluster needed)
k8s-validate:
    @if [ ! -f deployments/kubernetes/secrets.yaml ]; then cp deployments/kubernetes/secrets.example.yaml deployments/kubernetes/secrets.yaml; echo "ℹ️ Created deployments/kubernetes/secrets.yaml from example (fill real values before real apply)"; fi
    @kubectl kustomize deployments/kubernetes > /tmp/ecommerce-k8s.yaml
    @kubectl apply --dry-run=client --validate=false -f /tmp/ecommerce-k8s.yaml > /tmp/k8s-validate.log 2>&1; \
    status=$$?; \
    if [ $$status -ne 0 ]; then \
        echo "❌ kubectl dry-run failed (exit $$status)"; \
        head -20 /tmp/k8s-validate.log; \
        exit 1; \
    fi; \
    errors=$$(grep -icE 'error|invalid' /tmp/k8s-validate.log || true); \
    echo "✅ Client-side dry-run OK ($$errors validation lines)"

# Kubernetes: apply to current cluster (staging)
k8s-apply:
    kubectl apply -k deployments/kubernetes

# Kubernetes: wait for migration job then rollout status
k8s-rollout:
    kubectl wait --for=condition=complete job/migrate -n ecommerce --timeout=300s
    kubectl rollout status deployment/apigateway -n ecommerce --timeout=300s

# Kubernetes: rollback a deployment to previous revision
k8s-rollback svc:
    kubectl rollout undo deployment/{{svc}} -n ecommerce

# Run unit tests in pkg/
test-unit:
    @echo "🧪 Running unit tests in pkg/..."
    @cd pkg && go test ./... -v

# Run integration tests in tests/
test-integration:
    @echo "🧪 Running integration tests in tests/..."
    @cd tests && GOWORK=off APP_ENV=development go test ./... -v

# Run all tests (unit and integration)
test-all: test-unit test-integration
