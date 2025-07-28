-- +goose Up
-- +goose StatementBegin
CREATE TABLE "reset_tokens" (
    "id" SERIAL PRIMARY KEY,
    "user_id" BIGINT NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    "token" TEXT NOT NULL UNIQUE,
    "expiry_date" TIMESTAMP NOT NULL
);

CREATE INDEX idx_reset_token_token ON reset_tokens(token);
CREATE INDEX idx_reset_token_user_id ON reset_tokens(user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_reset_token_token";
DROP INDEX IF EXISTS "idx_reset_token_user_id";
DROP TABLE IF EXISTS "reset_tokens";

-- +goose StatementEnd
