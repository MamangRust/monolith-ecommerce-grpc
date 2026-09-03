package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MamangRust/monolith-ecommerce-grpc-transaction/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// noopLogger implements logger.LoggerInterface for unit tests.
type noopLogger struct{}

func (noopLogger) Info(string, ...zap.Field)                         {}
func (noopLogger) Fatal(string, ...zap.Field)                        {}
func (noopLogger) Debug(string, ...zap.Field)                        {}
func (noopLogger) Error(string, ...zap.Field)                        {}
func (noopLogger) Warn(string, ...zap.Field)                         {}
func (noopLogger) Check(zapcore.Level, string) *zapcore.CheckedEntry { return nil }
func (noopLogger) With(...zap.Field) logger.LoggerInterface          { return noopLogger{} }
func (noopLogger) Sync() error                                       { return nil }

// fakeOutboxRepository is an in-memory implementation of repository.OutboxRepository.
type fakeOutboxRepository struct {
	pending      []*db.OutboxEvent
	deliveredIDs map[int64]bool
	deadIDs      map[int64]bool
	failedIDs    map[int64]time.Time
}

func (f *fakeOutboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*db.OutboxEvent, error) {
	ev := &db.OutboxEvent{
		OutboxID:      int64(len(f.pending) + 1),
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		Attempts:      0,
		NextAttemptAt: time.Now(),
	}
	f.pending = append(f.pending, ev)
	return ev, nil
}

func (f *fakeOutboxRepository) CreateInTx(ctx context.Context, tx pgx.Tx, topic, key string, payload []byte) (*db.OutboxEvent, error) {
	return f.Create(ctx, topic, key, payload)
}

func (f *fakeOutboxRepository) GetPending(ctx context.Context, limit int) ([]*db.OutboxEvent, error) {
	return f.pending, nil
}

func (f *fakeOutboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*db.OutboxEvent, error) {
	claimed := f.pending
	if limit > 0 && len(claimed) > limit {
		claimed = claimed[:limit]
	}
	f.pending = nil
	return claimed, nil
}

func (f *fakeOutboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*db.OutboxEvent, error) {
	if f.deliveredIDs == nil {
		f.deliveredIDs = map[int64]bool{}
	}
	f.deliveredIDs[outboxID] = true
	return &db.OutboxEvent{OutboxID: outboxID, Status: "delivered"}, nil
}

func (f *fakeOutboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*db.OutboxEvent, error) {
	if f.failedIDs == nil {
		f.failedIDs = map[int64]time.Time{}
	}
	f.failedIDs[outboxID] = nextAttemptAt
	return &db.OutboxEvent{OutboxID: outboxID, Status: "pending"}, nil
}

func (f *fakeOutboxRepository) MarkDead(ctx context.Context, outboxID int64) (*db.OutboxEvent, error) {
	if f.deadIDs == nil {
		f.deadIDs = map[int64]bool{}
	}
	f.deadIDs[outboxID] = true
	return &db.OutboxEvent{OutboxID: outboxID, Status: "dead"}, nil
}

func (f *fakeOutboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	return 0, nil
}

// fakeOutboxPublisher fails for configured keys and records everything it sends.
type fakeOutboxPublisher struct {
	failKeys map[string]bool
	sent     []string
}

func (p *fakeOutboxPublisher) SendMessage(topic, key string, value []byte) error {
	p.sent = append(p.sent, key)
	if p.failKeys != nil && p.failKeys[key] {
		return errors.New("kafka unavailable")
	}
	return nil
}

func newOutboxServiceForTest(repo *fakeOutboxRepository, pub *fakeOutboxPublisher) OutboxService {
	return NewOutboxService(repo, pub, noopLogger{})
}

func TestOutboxEnqueueAndDeliver(t *testing.T) {
	repo := &fakeOutboxRepository{}
	pub := &fakeOutboxPublisher{}
	svc := newOutboxServiceForTest(repo, pub)

	err := svc.Enqueue(context.Background(), "topic-a", "key-1", []byte(`{"x":1}`))
	require.NoError(t, err)
	require.Len(t, repo.pending, 1)
	require.Equal(t, "pending", repo.pending[0].Status)

	delivered, err := svc.PublishPending(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, 1, delivered)
	require.True(t, repo.deliveredIDs[1], "event must be marked delivered")
	require.Equal(t, []string{"key-1"}, pub.sent)
}

func TestOutboxRetryBackoffOnFailure(t *testing.T) {
	repo := &fakeOutboxRepository{
		pending: []*db.OutboxEvent{{
			OutboxID: 7, Topic: "topic-a", EventKey: "fail", Payload: []byte(`{}`),
			Status: "pending", Attempts: 2,
		}},
	}
	pub := &fakeOutboxPublisher{failKeys: map[string]bool{"fail": true}}
	svc := newOutboxServiceForTest(repo, pub)

	delivered, err := svc.PublishPending(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, 0, delivered)

	nextAttempt, ok := repo.failedIDs[7]
	require.True(t, ok, "failed event must be scheduled for retry")
	require.True(t, nextAttempt.After(time.Now()), "retry must be scheduled in the future (backoff)")
	require.False(t, repo.deadIDs[7], "event must not be dead-lettered before max attempts")
}

func TestOutboxDeadLetterAfterMaxAttempts(t *testing.T) {
	repo := &fakeOutboxRepository{
		pending: []*db.OutboxEvent{{
			OutboxID: 9, Topic: "topic-a", EventKey: "dead", Payload: []byte(`{}`),
			Status: "pending", Attempts: OutboxMaxAttempts - 1, // one more failure exhausts retries
		}},
	}
	pub := &fakeOutboxPublisher{failKeys: map[string]bool{"dead": true}}
	svc := newOutboxServiceForTest(repo, pub)

	delivered, err := svc.PublishPending(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, 0, delivered)
	require.True(t, repo.deadIDs[9], "event must be dead-lettered after max attempts")
	_, scheduled := repo.failedIDs[9]
	require.False(t, scheduled, "dead-lettered event must not be scheduled for retry")
}

func TestOutboxNilPublisherIsNoOp(t *testing.T) {
	repo := &fakeOutboxRepository{}
	svc := NewOutboxService(repo, nil, noopLogger{})

	delivered, err := svc.PublishPending(context.Background(), 10)
	require.NoError(t, err)
	require.Equal(t, 0, delivered)
}

var _ repository.OutboxRepository = (*fakeOutboxRepository)(nil)
