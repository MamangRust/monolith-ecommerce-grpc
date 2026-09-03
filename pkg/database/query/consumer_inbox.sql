-- ReserveConsumerInbox: Claims an event and returns whether this call owns the reservation.
-- name: ReserveConsumerInbox :one
WITH reserved AS (
    INSERT INTO consumer_inbox (
        consumer_name, event_key, topic, partition_id, message_offset,
        status, attempts, reservation_version, lease_until, last_error, processed_at
    )
    VALUES ($1, $2, $3, $4, $5, 'processing', 1, 1,
            current_timestamp + interval '1 minute', '', NULL)
    ON CONFLICT (consumer_name, event_key) DO UPDATE
    SET status = 'processing',
        attempts = consumer_inbox.attempts + 1,
        reservation_version = consumer_inbox.reservation_version + 1,
        lease_until = current_timestamp + interval '1 minute',
        last_error = '',
        topic = EXCLUDED.topic,
        partition_id = EXCLUDED.partition_id,
        message_offset = EXCLUDED.message_offset
    WHERE consumer_inbox.status <> 'processed'
      AND consumer_inbox.lease_until <= current_timestamp
    RETURNING reservation_version
)
SELECT
    EXISTS (SELECT 1 FROM reserved) AS reserved,
    EXISTS (
        SELECT 1
        FROM consumer_inbox ci
        WHERE ci.consumer_name = $1
          AND ci.event_key = $2
          AND ci.status = 'processed'
    ) AS processed,
    COALESCE(
        (SELECT reservation_version FROM reserved),
        (SELECT ci.reservation_version FROM consumer_inbox ci WHERE ci.consumer_name = $1 AND ci.event_key = $2)
    )::BIGINT AS reservation_version;

-- MarkConsumerInboxProcessed: Completes only the active reservation.
-- name: MarkConsumerInboxProcessed :exec
UPDATE consumer_inbox
SET status = 'processed', processed_at = current_timestamp,
    lease_until = current_timestamp, last_error = ''
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $3;

-- ReleaseConsumerInbox: Releases only the active reservation.
-- name: ReleaseConsumerInbox :exec
UPDATE consumer_inbox
SET status = 'pending', lease_until = current_timestamp,
    last_error = $3
WHERE consumer_name = $1 AND event_key = $2
  AND status = 'processing' AND reservation_version = $4;
