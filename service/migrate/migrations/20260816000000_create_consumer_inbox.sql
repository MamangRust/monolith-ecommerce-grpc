-- +goose Up
-- +goose StatementBegin

-- Durable consumer-side deduplication (Phase 3). A message is marked only after
-- the handler's side effect succeeds; malformed/non-retryable messages may be
-- recorded separately by the consumer if desired. The reservation_version is
-- used for fencing: only the consumer that owns the active reservation may
-- mark or release it.
CREATE TABLE IF NOT EXISTS consumer_inbox (
    consumer_name VARCHAR(128) NOT NULL,
    event_key VARCHAR(255) NOT NULL,
    topic VARCHAR(255) NOT NULL DEFAULT '',
    partition_id INTEGER NOT NULL DEFAULT -1,
    message_offset BIGINT NOT NULL DEFAULT -1,
    status VARCHAR(16) NOT NULL DEFAULT 'processing'
        CHECK (status IN ('processing', 'processed', 'pending')),
    attempts INTEGER NOT NULL DEFAULT 1 CHECK (attempts >= 0),
    lease_until TIMESTAMPTZ NOT NULL DEFAULT current_timestamp + interval '1 minute',
    last_error TEXT NOT NULL DEFAULT '',
    processed_at TIMESTAMPTZ,
    reservation_version BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (consumer_name, event_key)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS consumer_inbox;
-- +goose StatementEnd
