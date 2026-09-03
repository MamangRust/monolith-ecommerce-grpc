package service

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

// Outbox constants control the durable retry behavior of the outbox relay.
const (
	OutboxMaxAttempts    = 5
	OutboxBackoff        = 30 * time.Second
	OutboxRelayInterval  = 5 * time.Second
	OutboxRelayBatchSize = 100
	// OutboxClaimLease is how long a relay worker owns a claimed event. If the
	// worker dies after claiming but before marking the event delivered, the
	// lease expires and another relay instance re-claims and retries it.
	OutboxClaimLease = 1 * time.Minute
	// OutboxRetention is how long delivered/dead events are kept before the
	// relay purges them as part of the retention policy.
	OutboxRetention = 7 * 24 * time.Hour

	// OutboxRetentionEveryTicks runs the retention purge every N relay ticks so
	// the DELETE scan does not run on every relay cycle.
	OutboxRetentionEveryTicks = 60
)

// OutboxPublisher is the minimal Kafka producer surface the relay needs.
// *kafka.Kafka satisfies this interface via SendMessage.
type OutboxPublisher interface {
	SendMessage(topic string, key string, value []byte) error
}

// OutboxService persists events durably after commit and relays them to Kafka
// with retry and dead-letter semantics.
type OutboxService interface {
	// Enqueue persists a pending event after the business transaction has already
	// committed. It is the NON-ATOMIC fallback path: a crash between the commit
	// and this insert silently loses the event, so it must not be used in
	// production. Production callers write the event inside the business
	// transaction via repository.OutboxRepository.CreateInTx and let the relay
	// publish it durably.
	Enqueue(ctx context.Context, topic, key string, payload []byte) error

	// PublishPending claims up to limit pending events whose retry window has
	// elapsed, publishes each to Kafka, and marks it delivered. Claiming is
	// atomic (FOR UPDATE SKIP LOCKED + lease), so concurrent relay instances
	// never publish the same event twice. It returns the number of events
	// successfully delivered.
	PublishPending(ctx context.Context, limit int) (int, error)

	// Start runs the relay loop until ctx is cancelled.
	Start(ctx context.Context, interval time.Duration, limit int)
}

type outboxService struct {
	repository repository.OutboxRepository
	publisher  OutboxPublisher
	logger     logger.LoggerInterface
}

func NewOutboxService(repo repository.OutboxRepository, publisher OutboxPublisher, log logger.LoggerInterface) OutboxService {
	return &outboxService{
		repository: repo,
		publisher:  publisher,
		logger:     log,
	}
}

func (s *outboxService) Enqueue(ctx context.Context, topic, key string, payload []byte) error {
	// Fallback path only: the event is inserted AFTER the business commit, so
	// this is best-effort and non-atomic. Production flows use CreateInTx.
	_, err := s.repository.Create(ctx, topic, key, payload)
	if err != nil {
		return err
	}
	s.logger.Info("outbox event enqueued", zap.String("topic", topic), zap.String("key", key))
	return nil
}

func (s *outboxService) PublishPending(ctx context.Context, limit int) (int, error) {
	if s.publisher == nil {
		return 0, nil
	}
	events, err := s.repository.Claim(ctx, limit, time.Now().Add(OutboxClaimLease))
	if err != nil {
		return 0, err
	}

	delivered := 0
	for _, event := range events {
		if err := s.publisher.SendMessage(event.Topic, event.EventKey, event.Payload); err != nil {
			s.logger.Error("failed to publish outbox event, scheduling retry",
				zap.Error(err),
				zap.Int64("outbox_id", event.OutboxID),
				zap.String("topic", event.Topic),
				zap.Int32("attempts", event.Attempts),
			)
			if int(event.Attempts)+1 >= OutboxMaxAttempts {
				if _, deadErr := s.repository.MarkDead(ctx, event.OutboxID); deadErr != nil {
					s.logger.Error("failed to dead-letter outbox event", zap.Error(deadErr), zap.Int64("outbox_id", event.OutboxID))
				}
				continue
			}
			nextAttempt := time.Now().Add(OutboxBackoff * time.Duration(event.Attempts+1))
			if _, failErr := s.repository.MarkFailed(ctx, event.OutboxID, nextAttempt); failErr != nil {
				s.logger.Error("failed to record outbox retry", zap.Error(failErr), zap.Int64("outbox_id", event.OutboxID))
			}
			continue
		}
		if _, err := s.repository.MarkDelivered(ctx, event.OutboxID); err != nil {
			s.logger.Error("failed to mark outbox event delivered", zap.Error(err), zap.Int64("outbox_id", event.OutboxID))
			continue
		}
		delivered++
	}
	return delivered, nil
}

func (s *outboxService) Start(ctx context.Context, interval time.Duration, limit int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	tickCount := 0
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("outbox relay stopped")
			return
		case <-ticker.C:
			if _, err := s.PublishPending(ctx, limit); err != nil {
				s.logger.Error("outbox relay cycle failed", zap.Error(err))
			}
			tickCount++
			// Retention runs periodically (not every tick) to avoid scanning the
			// outbox table on every relay cycle; it purges delivered/dead events
			// whose terminal state is older than the retention window.
			if tickCount%OutboxRetentionEveryTicks == 0 {
				if removed, err := s.repository.DeleteOld(ctx, time.Now().Add(-OutboxRetention)); err != nil {
					s.logger.Error("outbox retention cleanup failed", zap.Error(err))
				} else if removed > 0 {
					s.logger.Info("outbox retention cleanup", zap.Int64("removed", removed))
				}
			}
		}
	}
}
