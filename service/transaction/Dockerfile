FROM golang:1.25.7-bookworm AS builder

WORKDIR /app

# Copy workspace and dependency manifests
COPY go.work go.work.sum* ./
COPY pb/go.mod pb/go.sum* ./pb/
COPY pkg/go.mod pkg/go.sum* ./pkg/
COPY shared/go.mod shared/go.sum* ./shared/
COPY tests/go.mod tests/go.sum* ./tests/

# Copy all service go.mod files
COPY service/ ./service/
RUN find service -type f ! -name "go.mod" ! -name "go.sum" -delete

WORKDIR /app/service/transaction

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

WORKDIR /app

# Copy source code
COPY pb/ ./pb/
COPY pkg/ ./pkg/
COPY shared/ ./shared/
COPY service/transaction/ ./service/transaction/

WORKDIR /app/service/transaction

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/transaction_bin ./cmd/main.go


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && \
    adduser -D -s /bin/sh -u 1000 -G appuser appuser

WORKDIR /app

RUN mkdir -p /var/log/app && chown -R 1000:1000 /var/log/app

COPY --from=builder /app/transaction_bin ./transaction

RUN chown -R 1000:1000 ./transaction 

USER appuser

EXPOSE 50058

CMD ["./transaction"]
