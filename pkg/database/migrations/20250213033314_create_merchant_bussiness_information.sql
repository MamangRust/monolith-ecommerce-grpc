-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_business_information" (
    "merchant_business_info_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id") ON DELETE CASCADE,
    "business_type" VARCHAR(100),
    "tax_id" VARCHAR(50),
    "established_year" INT,
    "number_of_employees" INT,
    "website_url" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_merchant_business_info_merchant_id" ON "merchant_business_information"("merchant_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_merchant_business_info_merchant_id";
DROP TABLE IF EXISTS "merchant_business_information";
-- +goose StatementEnd