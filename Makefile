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
