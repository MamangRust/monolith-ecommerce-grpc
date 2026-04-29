-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchants" (
    "merchant_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "address" TEXT,
    "contact_email" VARCHAR(100),
    "contact_phone" VARCHAR(20),
    "status" VARCHAR(20) NOT NULL DEFAULT 'active',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_merchants_user_id" ON "merchants" ("user_id");

CREATE INDEX "idx_merchants_status" ON "merchants" ("status");

CREATE INDEX "idx_merchants_created_at" ON "merchants" ("created_at");
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_merchants_user_id" ON "merchants" ("user_id");

DROP INDEX IF EXISTS "idx_merchants_status" ON "merchants" ("status");

DROP INDEX IF EXISTS "idx_merchants_created_at" ON "merchants" ("created_at");

DROP TABLE IF EXISTS "merchants";

-- +goose StatementEnd
