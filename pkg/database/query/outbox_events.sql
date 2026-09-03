-- CreateOutboxEvent: Inserts a pending outbox event for durable Kafka publish.
-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (topic, event_key, payload, status, next_attempt_at)
VALUES ($1, $2, $3, 'pending', CURRENT_TIMESTAMP)
RETURNING outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at;

-- GetPendingOutboxEvents: Returns pending events whose retry window has elapsed.
-- name: GetPendingOutboxEvents :many
SELECT outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at
FROM outbox_events
WHERE status = 'pending' AND next_attempt_at <= CURRENT_TIMESTAMP
ORDER BY outbox_id
LIMIT $1;

-- ClaimPendingOutboxEvents: Atomically claims up to $1 pending events whose retry
-- window has elapsed. The claim extends next_attempt_at to $2 (the lease expiry)
-- so concurrent relay instances can never publish the same event twice; a crashed
-- worker's claim simply expires when the lease passes and the event is retried.
-- name: ClaimPendingOutboxEvents :many
UPDATE outbox_events
SET next_attempt_at = $2, updated_at = CURRENT_TIMESTAMP
WHERE outbox_id IN (
    SELECT outbox_id
    FROM outbox_events
    WHERE status = 'pending' AND next_attempt_at <= CURRENT_TIMESTAMP
    ORDER BY outbox_id
    LIMIT $1
    FOR UPDATE SKIP LOCKED
)
RETURNING outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at;

-- MarkOutboxEventDelivered: Marks a pending event as delivered.
-- name: MarkOutboxEventDelivered :one
UPDATE outbox_events
SET status = 'delivered', updated_at = CURRENT_TIMESTAMP
WHERE outbox_id = $1 AND status = 'pending'
RETURNING outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at;

-- MarkOutboxEventFailed: Increments the attempt counter and schedules the next
-- attempt at the provided timestamp (backoff computed by the caller).
-- name: MarkOutboxEventFailed :one
UPDATE outbox_events
SET attempts = attempts + 1,
    next_attempt_at = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE outbox_id = $1 AND status = 'pending'
RETURNING outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at;

-- MarkOutboxEventDead: Marks a pending event as dead after exhausting retries.
-- name: MarkOutboxEventDead :one
UPDATE outbox_events
SET status = 'dead', updated_at = CURRENT_TIMESTAMP
WHERE outbox_id = $1 AND status = 'pending'
RETURNING outbox_id, topic, event_key, payload, status, attempts, next_attempt_at, created_at, updated_at;

-- DeleteOldOutboxEvents: Purges delivered/dead events older than the cutoff.
-- name: DeleteOldOutboxEvents :execrows
DELETE FROM outbox_events
WHERE status IN ('delivered', 'dead')
  AND updated_at < $1;
