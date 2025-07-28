-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_social_media_links" (
    "merchant_social_id" SERIAL PRIMARY KEY,
    "merchant_detail_id" INT NOT NULL REFERENCES merchant_details("merchant_detail_id") ON DELETE CASCADE,
    "platform" VARCHAR(100) NOT NULL,
    "url" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "idx_merchant_social_media_links_merchant_detail_id" ON "merchant_social_media_links" ("merchant_detail_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_merchant_social_media_links_merchant_detail_id";
DROP TABLE IF EXISTS "merchant_social_media_links";
-- +goose StatementEnd
