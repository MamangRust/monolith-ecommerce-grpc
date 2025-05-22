FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN mkdir -p ./pkg/database/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app


COPY --from=builder /app/migrate /usr/local/bin/migrate
COPY --from=builder /app/docker.env* ./
COPY --from=builder /app/pkg/database/migrations ./pkg/database/migrations

ENTRYPOINT ["/usr/local/bin/migrate", "up"]