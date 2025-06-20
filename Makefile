COMPOSE_FILE=deployments/local/docker-compose.yml
SERVICES := apigateway migrate auth user role banner cart category email merchant merchant_award merchant_business merchant_detail merchant_policy order order_item product review review_detail shipping_address slider transaction
DOCKER_COMPOSE=docker compose

migrate:
	go run service/migrate/main.go up

migrate-down:
	go run service/migrate/main.go down


generate-proto:
	protoc --proto_path=pkg/proto --go_out=shared/pb --go_opt=paths=source_relative --go-grpc_out=shared/pb --go-grpc_opt=paths=source_relative pkg/proto/*.proto


generate-swagger:
	swag init -g service/apigateway/cmd/main.go -o service/apigateway/docs


seeder:
	go run service/seeder/main.go


build-image:
	@for service in $(SERVICES); do \
		echo "🔨 Building $$service-ecommerce-service..."; \
		docker build -t $$service-ecommerce-service:1.0 -f service/$$service/Dockerfile service/$$service || exit 1; \
	done
	@echo "✅ All services built successfully."


image-load:
	@for service in $(SERVICES); do \
		echo "🚚 Loading $$service-service..."; \
		minikube image load $$service-service:1.0 || exit 1; \
	done
	@echo "✅ All services loaded successfully."


go-mod-tidy:
	@for service in $(SERVICES); do \
		echo "🔧 Running go mod tidy for $$service..."; \
		cd service/$$service && go mod tidy || exit 1; \
	done
	@echo "✅ All go.mod tidy completed."


ps:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) ps

up:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) up -d

down:
	${DOCKER_COMPOSE} -f $(COMPOSE_FILE) down

build-up:
	make build-image && make up