package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
)

// OutboxRepository provides durable persistence for events that must be published
// to Kafka after the business transaction commits. A relay service consumes
// pending events and retries until delivered or dead-lettered.
type OutboxRepository interface {
	Create(ctx context.Context, topic, key string, payload []byte) (*db.OutboxEvent, error)

	// CreateInTx persists a pending event inside the given database transaction so
	// the caller can commit the business write and the event atomically.
	CreateInTx(ctx context.Context, tx pgx.Tx, topic, key string, payload []byte) (*db.OutboxEvent, error)

	GetPending(ctx context.Context, limit int) ([]*db.OutboxEvent, error)

	// Claim atomically claims up to limit pending events whose retry window has
	// elapsed by extending next_attempt_at to leaseUntil. The FOR UPDATE SKIP
	// LOCKED guard ensures concurrent relay instances never publish the same
	// event; a crashed worker's claim simply expires when the lease passes.
	Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*db.OutboxEvent, error)

	MarkDelivered(ctx context.Context, outboxID int64) (*db.OutboxEvent, error)
	MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*db.OutboxEvent, error)
	MarkDead(ctx context.Context, outboxID int64) (*db.OutboxEvent, error)
	DeleteOld(ctx context.Context, cutoff time.Time) (int64, error)
}

type outboxRepository struct {
	db *db.Queries
}

func NewOutboxRepository(queries *db.Queries) OutboxRepository {
	return &outboxRepository{db: queries}
}

func (r *outboxRepository) Create(ctx context.Context, topic, key string, payload []byte) (*db.OutboxEvent, error) {
	return r.db.CreateOutboxEvent(ctx, db.CreateOutboxEventParams{
		Topic:    topic,
		EventKey: key,
		Payload:  payload,
	})
}

func (r *outboxRepository) CreateInTx(ctx context.Context, tx pgx.Tx, topic, key string, payload []byte) (*db.OutboxEvent, error) {
	return r.db.WithTx(tx).CreateOutboxEvent(ctx, db.CreateOutboxEventParams{
		Topic:    topic,
		EventKey: key,
		Payload:  payload,
	})
}

func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*db.OutboxEvent, error) {
	return r.db.GetPendingOutboxEvents(ctx, int32(limit))
}

func (r *outboxRepository) Claim(ctx context.Context, limit int, leaseUntil time.Time) ([]*db.OutboxEvent, error) {
	return r.db.ClaimPendingOutboxEvents(ctx, db.ClaimPendingOutboxEventsParams{
		Limit:         int32(limit),
		NextAttemptAt: leaseUntil,
	})
}

func (r *outboxRepository) MarkDelivered(ctx context.Context, outboxID int64) (*db.OutboxEvent, error) {
	return r.db.MarkOutboxEventDelivered(ctx, outboxID)
}

func (r *outboxRepository) MarkFailed(ctx context.Context, outboxID int64, nextAttemptAt time.Time) (*db.OutboxEvent, error) {
	return r.db.MarkOutboxEventFailed(ctx, db.MarkOutboxEventFailedParams{
		OutboxID:      outboxID,
		NextAttemptAt: nextAttemptAt,
	})
}

func (r *outboxRepository) MarkDead(ctx context.Context, outboxID int64) (*db.OutboxEvent, error) {
	return r.db.MarkOutboxEventDead(ctx, outboxID)
}

func (r *outboxRepository) DeleteOld(ctx context.Context, cutoff time.Time) (int64, error) {
	return r.db.DeleteOldOutboxEvents(ctx, pgtype.Timestamp{Time: cutoff, Valid: true})
}
