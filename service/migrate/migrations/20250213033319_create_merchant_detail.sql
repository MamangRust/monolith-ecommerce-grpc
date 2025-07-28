-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_details" (
    "merchant_detail_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id") ON DELETE CASCADE,
    "display_name" VARCHAR(255),
    "cover_image_url" VARCHAR(255),
    "logo_url" VARCHAR(255),
    "short_description" TEXT,
    "website_url" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_merchant_details_merchant" ON "merchant_details" ("merchant_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_merchant_details_merchant";
DROP TABLE IF EXISTS "merchant_details";
-- +goose StatementEnd