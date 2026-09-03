package outbox

import (
	"context"
	"errors"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// ConsumerInbox is the durable deduplication contract used by Kafka handlers
// Durable consumer-side deduplication. It replaces in-memory-only deduplication:
// reservations survive consumer restarts and rebalances, so at-least-once
// redelivery cannot send the same email twice.
type ConsumerInbox interface {
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}

// InboxExecutor is sqlc's generated DBTX contract. Using it here guarantees
// inbox persistence follows the same query-file/code-generation rule as the
// domain repositories.
type InboxExecutor interface {
	db.DBTX
}

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
func Reserve(ctx context.Context, tx InboxExecutor, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if tx == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}
	reservation, err := db.New(tx).ReserveConsumerInbox(ctx, db.ReserveConsumerInboxParams{
		ConsumerName: consumerName, EventKey: eventKey, Topic: topic,
		PartitionID: partition, MessageOffset: offset,
	})
	if err != nil {
		return false, false, 0, err
	}
	return reservation.Reserved, reservation.Processed, reservation.ReservationVersion, nil
}

func MarkProcessed(ctx context.Context, tx InboxExecutor, consumerName, eventKey string, reservationVersion int64) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	return db.New(tx).MarkConsumerInboxProcessed(ctx, db.MarkConsumerInboxProcessedParams{
		ConsumerName: consumerName, EventKey: eventKey, ReservationVersion: reservationVersion,
	})
}

func Release(ctx context.Context, tx InboxExecutor, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if tx == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	return db.New(tx).ReleaseConsumerInbox(ctx, db.ReleaseConsumerInboxParams{
		ConsumerName: consumerName, EventKey: eventKey, LastError: lastError, ReservationVersion: reservationVersion,
	})
}

// PostgresInbox adapts a pgx pool to ConsumerInbox. Reservation and completion
// are committed independently because an external side effect cannot share a
// PostgreSQL transaction with the Kafka consumer.
type PostgresInbox struct {
	pool *pgxpool.Pool
}

func NewPostgresInbox(pool *pgxpool.Pool) (*PostgresInbox, error) {
	if pool == nil {
		return nil, errors.New("inbox pool is nil")
	}
	return &PostgresInbox{pool: pool}, nil
}

func (i *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	tx, err := i.pool.Begin(ctx)
	if err != nil {
		return false, false, 0, err
	}
	defer tx.Rollback(ctx)
	reserved, processed, reservationVersion, err := Reserve(ctx, tx, consumerName, eventKey, topic, partition, offset)
	if err != nil {
		return false, false, 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return false, false, 0, err
	}
	return reserved, processed, reservationVersion, nil
}

func (i *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	tx, err := i.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := MarkProcessed(ctx, tx, consumerName, eventKey, reservationVersion); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (i *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	tx, err := i.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := Release(ctx, tx, consumerName, eventKey, reservationVersion, processingErr); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
