package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-apigateway/internal/redis"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RoleValidator is responsible for validating user roles via Kafka-based request-response messaging.
//
// This struct sends role validation requests through Kafka and listens for asynchronous responses.
// It includes logic to manage concurrent access to response channels and timeout control.
type RoleValidator struct {
	// kafka is the Kafka client used to send role validation requests and listen for responses.
	kafka *kafka.Kafka
	// logger provides structured logging for tracking validator behavior and errors.
	logger logger.LoggerInterface
	// requestTopic is the Kafka topic where role validation requests are published.
	requestTopic string
	// responseTopic is the Kafka topic where role validation responses are received.
	responseTopic string
	// timeout defines how long the validator waits for a response before timing out.
	timeout time.Duration
	// responseChans is a map of correlation IDs to channels used for receiving responses.
	// It enables concurrent request-response tracking for each role validation request.
	responseChans map[string]chan *response.RoleResponsePayload
	// mu is a read-write mutex used to safely access the responseChans map concurrently.
	mu sync.RWMutex

	cache mencache.RoleCache
}

// NewRoleValidator creates a new RoleValidator instance.
//
// It starts a Kafka consumer that listens to the given response topic and waits
// for a response to the validation request published to the given request topic.
// The validator is used as a middleware to validate the role ID in the request
// header and store the role data in the Echo context.
//
// It panics if the Kafka consumer cannot be started.
func NewRoleValidator(k *kafka.Kafka, requestTopic, responseTopic string, timeout time.Duration, logger logger.LoggerInterface, cache mencache.RoleCache) *RoleValidator {
	v := &RoleValidator{
		kafka:         k,
		requestTopic:  requestTopic,
		responseTopic: responseTopic,
		timeout:       timeout,
		cache:         cache,
		responseChans: make(map[string]chan *response.RoleResponsePayload),
		logger:        logger,
	}
	handler := &roleResponseHandler{validator: v}
	go func() {
		err := k.StartConsumers([]string{responseTopic}, "role-validator-gateway", handler)
		if err != nil {
			v.logger.Fatal("Failed to start kafka consumer", zap.Error(err))
			panic("failed to start kafka consumer: " + err.Error())
		}
	}()
	return v
}

// Middleware returns an Echo middleware that validates the user role by publishing
// a message to a Kafka topic and waiting for a response. If the validation is
// successful, the role names are stored in the Echo context and the request is
// passed to the next handler. If the validation fails or times out, an HTTP error
// is returned.
func (v *RoleValidator) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDVal := c.Get("user_id")
			v.logger.Debug("Validating user role", zap.Any("user_id", userIDVal))
			if userIDVal == nil {
				v.logger.Error("User ID not found in context")
				return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
			}
			userID, err := v.extractUserID(userIDVal)
			if err != nil {
				v.logger.Error("Invalid User ID format", zap.Any("value", userIDVal), zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID format")
			}

			if roles, found := v.cache.GetRoleCache(c.Request().Context(), strconv.Itoa(userID)); found {
				v.logger.Debug("Role found in cache", zap.Int("user_id", userID), zap.Strings("roles", roles))
				c.Set("role_names", roles)
				return next(c)
			}

			correlationID := uuid.NewString()
			v.logger.Info("Validating user role via Kafka", zap.Int("user_id", userID), zap.String("correlation_id", correlationID))

			respChan := make(chan *response.RoleResponsePayload, 1)

			v.mu.Lock()
			v.responseChans[correlationID] = respChan
			v.mu.Unlock()

			defer func() {
				v.mu.Lock()
				delete(v.responseChans, correlationID)
				close(respChan)
				v.mu.Unlock()
			}()

			if err := v.sendValidationRequest(userID, correlationID); err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(c.Request().Context(), v.timeout)
			defer cancel()
			select {
			case roleResponse := <-respChan:
				if roleResponse == nil {
					v.logger.Error("Received nil response", zap.String("correlation_id", correlationID))
					return echo.NewHTTPError(http.StatusInternalServerError, "Invalid response received")
				}
				if !roleResponse.Valid || len(roleResponse.RoleNames) == 0 {
					v.logger.Debug("Role validation failed",
						zap.Int("user_id", userID),
						zap.String("correlation_id", correlationID),
						zap.Bool("valid", roleResponse.Valid),
						zap.Int("role_count", len(roleResponse.RoleNames)))
					return echo.NewHTTPError(http.StatusUnauthorized, "Role validation failed")
				}

				v.logger.Info("Role validation success",
					zap.Int("user_id", userID),
					zap.Strings("roles", roleResponse.RoleNames),
					zap.String("correlation_id", correlationID))

				v.cache.SetRoleCache(ctx, strconv.Itoa(userID), roleResponse.RoleNames)

				c.Set("role_names", roleResponse.RoleNames)
				return next(c)

			case <-ctx.Done():
				v.logger.Error("Timeout waiting for Kafka response",
					zap.String("correlation_id", correlationID),
					zap.Duration("timeout", v.timeout))
				return echo.NewHTTPError(http.StatusRequestTimeout, "Timeout waiting for role validation")
			}
		}
	}
}

func (v *RoleValidator) extractUserID(userIDVal interface{}) (int, error) {
	switch val := userIDVal.(type) {
	case float64:
		return int(val), nil
	case int:
		return val, nil
	case string:
		return strconv.Atoi(val)
	default:
		return 0, echo.NewHTTPError(http.StatusUnauthorized, "Unknown user ID type")
	}
}

func (v *RoleValidator) sendValidationRequest(userID int, correlationID string) error {
	payload := requests.RoleRequestPayload{
		UserID:        userID,
		CorrelationID: correlationID,
		ReplyTopic:    v.responseTopic,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		v.logger.Error("Failed to encode payload", zap.Error(err), zap.String("correlation_id", correlationID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to encode payload")
	}

	err = v.kafka.SendMessage(v.requestTopic, correlationID, data)
	if err != nil {
		v.logger.Error("Failed to send Kafka message", zap.Error(err), zap.String("correlation_id", correlationID))
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send Kafka message")
	}
	v.logger.Info("Kafka message sent for role validation",
		zap.String("topic", v.requestTopic),
		zap.String("correlation_id", correlationID))
	return nil
}
