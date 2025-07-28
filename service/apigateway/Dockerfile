FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o apigateway ./cmd/main.go

FROM alpine:latest

RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /var/log/app && chown -R 1000:1000 /var/log/app

COPY --from=builder /app/apigateway .

RUN chown 1000:1000 ./apigateway

USER appuser

EXPOSE 5000 8100

CMD ["./apigateway"]
