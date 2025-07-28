-- +goose Up
-- +goose StatementBegin
CREATE TABLE "transactions" (
    "transaction_id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL REFERENCES "orders" ("order_id"),
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id"),
    "payment_method" VARCHAR(50) NOT NULL,
    "amount" INT NOT NULL,
    "payment_status" VARCHAR(20) NOT NULL DEFAULT 'completed',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);
CREATE INDEX "idx_transactions_order_id" ON "transactions"("order_id");
CREATE INDEX "idx_transactions_merchant_id" ON "transactions"("merchant_id");
CREATE INDEX "idx_transactions_payment_status" ON "transactions"("payment_status");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_transactions_order_id";
DROP INDEX IF EXISTS "idx_transactions_merchant_id";
DROP INDEX IF EXISTS "idx_transactions_payment_status";
DROP TABLE IF EXISTS "transactions";
-- +goose StatementEnd
