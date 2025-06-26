package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type RoleValidator struct {
	kafka         *kafka.Kafka
	logger        logger.LoggerInterface
	requestTopic  string
	responseTopic string
	timeout       time.Duration
	responseChans map[string]chan *response.RoleResponsePayload
	mu            sync.RWMutex
}

func NewRoleValidator(k *kafka.Kafka, requestTopic, responseTopic string, timeout time.Duration, logger logger.LoggerInterface) *RoleValidator {
	v := &RoleValidator{
		kafka:         k,
		requestTopic:  requestTopic,
		responseTopic: responseTopic,
		timeout:       timeout,
		responseChans: make(map[string]chan *response.RoleResponsePayload),
		logger:        logger,
	}

	handler := &roleResponseHandler{validator: v}
	go func() {
		err := k.StartConsumers([]string{responseTopic}, "role-validator-group", handler)
		if err != nil {
			v.logger.Fatal("Failed to start kafka consumer", zap.Error(err))
			panic("failed to start kafka consumer: " + err.Error())
		}
	}()

	return v
}

func (v *RoleValidator) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDVal := c.Get("user_id")
			if userIDVal == nil {
				v.logger.Error("User ID not found in context")
				return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found in context")
			}

			userID, err := v.extractUserID(userIDVal)
			if err != nil {
				v.logger.Error("Invalid User ID format", zap.Any("value", userIDVal), zap.Error(err))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid User ID format")
			}

			correlationID := uuid.NewString()
			v.logger.Info("Validating user role", zap.Int("user_id", userID), zap.String("correlation_id", correlationID))

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
