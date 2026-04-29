package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "requestId"
const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			return "", err
		}
		result[i] = characters[num.Int64()]
	}
	return string(result), nil
}

func GenerateRequestID() string {
	rand, _ := GenerateRandomString(8)
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), rand)
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func RequestID(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// context.go

const startTimeKey ctxKey = "startTime"

func WithStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, startTimeKey, t)
}

func StartTime(ctx context.Context) time.Time {
	if v, ok := ctx.Value(startTimeKey).(time.Time); ok {
		return v
	}
	return time.Time{}
}

const (
	methodKey    ctxKey = "method"
	operationKey ctxKey = "operation"
)

func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey, method)
}

func Method(ctx context.Context) string {
	if v, ok := ctx.Value(methodKey).(string); ok {
		return v
	}
	return ""
}

func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

func Operation(ctx context.Context) string {
	if v, ok := ctx.Value(operationKey).(string); ok {
		return v
	}
	return ""
}
