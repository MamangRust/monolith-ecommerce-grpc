FROM golang:1.23.9-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o order_item ./cmd/main.go

FROM alpine:latest

RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

RUN apk --no-cache add ca-certificates

WORKDIR /app

RUN mkdir -p /var/log/app && chown -R 1000:1000 /var/log/app

COPY --from=builder /app/order_item .

RUN chown 1000:1000 ./order_item

USER appuser

EXPOSE 50056 8086

CMD ["./order_item"]
