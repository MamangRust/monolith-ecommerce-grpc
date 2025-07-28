-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_policies" (
    "merchant_policy_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id") ON DELETE CASCADE,
    "policy_type" VARCHAR(100) NOT NULL, 
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_policies_merchant_id" ON "merchant_policies"("merchant_id");
CREATE INDEX "idx_policies_policy_type" ON "merchant_policies"("policy_type");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_policies_policy_type";
DROP INDEX IF EXISTS "idx_policies_merchant_id";
DROP TABLE IF EXISTS "merchant_policies";
-- +goose StatementEnd