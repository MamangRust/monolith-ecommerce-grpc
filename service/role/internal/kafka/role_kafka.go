package myhandlerkafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-grpc-role/internal/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/kafka"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/response"
	"go.uber.org/zap"
)

type roleKafkaHandler struct {
	logger      logger.LoggerInterface
	roleService service.RoleQueryService
	kafka       *kafka.Kafka
	ctx         context.Context
}

func NewRoleKafkaHandler(ctx context.Context, roleService service.RoleQueryService, kafka *kafka.Kafka, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &roleKafkaHandler{
		ctx:         ctx,
		roleService: roleService,
		kafka:       kafka,
		logger:      logger,
	}
}

func (h *roleKafkaHandler) Setup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Role Kafka handler setup completed")
	return nil
}

func (h *roleKafkaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	h.logger.Info("Role Kafka handler cleanup completed")
	return nil
}

func (h *roleKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		msgCtx, cancel := context.WithTimeout(h.ctx, 20*time.Second)
		defer cancel()

		h.logger.Info("Received role validation request",
			zap.String("topic", msg.Topic),
			zap.String("key", string(msg.Key)))

		var payload requests.RoleRequestPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			h.logger.Error("Invalid role request payload", zap.Error(err))
			session.MarkMessage(msg, "")
			continue
		}

		if payload.CorrelationID == "" || payload.ReplyTopic == "" {
			h.logger.Error("Missing required fields in role request",
				zap.String("correlation_id", payload.CorrelationID),
				zap.String("reply_topic", payload.ReplyTopic))
			session.MarkMessage(msg, "")
			continue
		}

		// Log pemrosesan
		h.logger.Info("Processing role validation request",
			zap.Int("user_id", payload.UserID),
			zap.String("correlation_id", payload.CorrelationID))

		// Panggil service untuk mendapatkan role
		roles, errResp := h.roleService.FindByUserId(msgCtx, payload.UserID)

		// Siapkan payload respons
		resp := response.RoleResponsePayload{
			CorrelationID: payload.CorrelationID,
			Valid:         errResp == nil && len(roles) > 0,
			RoleNames:     make([]string, 0),
		}

		// Isi nama role jika valid
		if errResp == nil && len(roles) > 0 {
			for _, r := range roles {
				resp.RoleNames = append(resp.RoleNames, r.Name)
			}
			h.logger.Info("Role validation successful",
				zap.Int("user_id", payload.UserID),
				zap.Strings("roles", resp.RoleNames),
				zap.String("correlation_id", payload.CorrelationID))
		} else {
			h.logger.Debug("Role validation failed",
				zap.Int("user_id", payload.UserID),
				zap.Any("error", errResp),
				zap.String("correlation_id", payload.CorrelationID))
		}

		respBytes, err := json.Marshal(resp)
		if err != nil {
			h.logger.Error("Failed to marshal role response",
				zap.Any("error", err),
				zap.String("correlation_id", payload.CorrelationID))
			session.MarkMessage(msg, "")
			continue
		}

		err = h.kafka.SendMessage(payload.ReplyTopic, payload.CorrelationID, respBytes)
		if err != nil {
			h.logger.Error("Failed to send Kafka role response",
				zap.Any("error", err),
				zap.String("reply_topic", payload.ReplyTopic),
				zap.String("correlation_id", payload.CorrelationID))
		} else {
			h.logger.Info("Role response sent successfully",
				zap.String("reply_topic", payload.ReplyTopic),
				zap.String("correlation_id", payload.CorrelationID))
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
