-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_certifications_and_awards" (
    "merchant_certification_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id") ON DELETE CASCADE,
    "title" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "issued_by" VARCHAR(255),
    "issue_date" DATE,
    "expiry_date" DATE,
    "certificate_url" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_certifications_merchant_id" ON "merchant_certifications_and_awards"("merchant_id");
CREATE INDEX "idx_certifications_issue_date" ON "merchant_certifications_and_awards"("issue_date");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_certifications_issue_date";
DROP INDEX IF EXISTS "idx_certifications_merchant_id";
DROP TABLE IF EXISTS "merchant_certifications_and_awards";
-- +goose StatementEnd