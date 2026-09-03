package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	transactionCache "github.com/MamangRust/monolith-ecommerce-grpc-transaction/cache"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

type merchantStatusEvent struct {
	MerchantID int32  `json:"merchantId"`
	Status     string `json:"status"`
}

type merchantStatusConsumer struct {
	ctx    context.Context
	cache  transactionCache.TransactionCommandCache
	logger logger.LoggerInterface
}

// NewMerchantStatusConsumer handles merchant status events for transaction cache coherence.
func NewMerchantStatusConsumer(ctx context.Context, cache transactionCache.TransactionCommandCache, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &merchantStatusConsumer{ctx: ctx, cache: cache, logger: logger}
}

func (h *merchantStatusConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *merchantStatusConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *merchantStatusConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event merchantStatusEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			h.logger.Error("invalid merchant status event payload", zap.Error(err), zap.String("topic", message.Topic))
			session.MarkMessage(message, "")
			continue
		}
		if h.cache != nil {
			h.cache.InvalidateTransactionCache(h.ctx)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
