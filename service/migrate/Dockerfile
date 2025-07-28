FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o migrate ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/migrate ./migrate
RUN chown appuser:appuser ./migrate

COPY --from=builder /app/migrations ./migrations
RUN chown -R appuser:appuser ./migrations

USER appuser

ENTRYPOINT ["./migrate", "up"]