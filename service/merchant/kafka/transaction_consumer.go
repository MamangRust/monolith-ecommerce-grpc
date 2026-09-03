package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	merchantCache "github.com/MamangRust/monolith-ecommerce-grpc-merchant/cache"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

type transactionEvent struct {
	MerchantID int32 `json:"merchantId"`
}

type transactionConsumer struct {
	ctx    context.Context
	cache  merchantCache.MerchantCommandCache
	logger logger.LoggerInterface
}

// NewTransactionConsumer handles transaction events for merchant cache coherence.
func NewTransactionConsumer(ctx context.Context, cache merchantCache.MerchantCommandCache, logger logger.LoggerInterface) sarama.ConsumerGroupHandler {
	return &transactionConsumer{ctx: ctx, cache: cache, logger: logger}
}

func (h *transactionConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *transactionConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *transactionConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var event transactionEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			h.logger.Error("invalid transaction event payload", zap.Error(err), zap.String("topic", message.Topic))
			session.MarkMessage(message, "")
			continue
		}
		if h.cache != nil && event.MerchantID > 0 {
			h.cache.DeleteCachedMerchant(h.ctx, int(event.MerchantID))
		}
		session.MarkMessage(message, "")
	}
	return nil
}
