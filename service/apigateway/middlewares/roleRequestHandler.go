package middlewares

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.uber.org/zap"
)

type roleResponseHandler struct {
	validator *RoleValidator
}

func (h *roleResponseHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *roleResponseHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *roleResponseHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.validator.logger.Debug("Received Kafka response message",
			zap.String("topic", msg.Topic),
			zap.String("key", string(msg.Key)))

		var roleResponse response.RoleResponsePayload
		if err := json.Unmarshal(msg.Value, &roleResponse); err != nil {
			h.validator.logger.Error("Failed to unmarshal role response", zap.Error(err))
			session.MarkMessage(msg, "")
			continue
		}

		correlationID := roleResponse.CorrelationID
		if correlationID == "" {
			h.validator.logger.Error("Missing correlation ID in response")
			session.MarkMessage(msg, "")
			continue
		}

		h.validator.mu.RLock()
		ch, exists := h.validator.responseChans[correlationID]
		h.validator.mu.RUnlock()

		if exists && ch != nil {
			select {
			case ch <- &roleResponse:
				h.validator.logger.Debug("Response delivered to channel",
					zap.String("correlation_id", correlationID))
			default:
				h.validator.logger.Debug("Response channel full or closed",
					zap.String("correlation_id", correlationID))
			}
		} else {
			h.validator.logger.Debug("No waiting channel for correlation ID",
				zap.String("correlation_id", correlationID))
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
