set shell := ["bash", "-c"]

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
    protoc --proto_path=proto --go_out=shared/pb --go_opt=module=github.com/MamangRust/monolith-ecommerce-shared/pb --go-grpc_out=shared/pb --go-grpc_opt=module=github.com/MamangRust/monolith-ecommerce-shared/pb $(find proto -name "*.proto")

# Generate Swagger
generate-swagger:
    swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs

# Run Seeder
seeder:
    go run service/seeder/cmd/main.go

# Build Docker images for all services
build-image:
    @for service in apigateway migrate auth user role banner cart category email merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product review review_detail shipping_address slider transaction; do \
        echo "🔨 Building $service-ecommerce-service..."; \
        docker build -t $service-ecommerce-service:1.0 -f service/$service/Dockerfile service/$service || exit 1; \
    done
    @echo "✅ All service images built successfully."

# Docker Compose up
up:
    docker compose -f deployments/local/docker-compose.yml up -d

# Docker Compose down
down:
    docker compose -f deployments/local/docker-compose.yml down

# Build images and start compose
build-up: build-image up
