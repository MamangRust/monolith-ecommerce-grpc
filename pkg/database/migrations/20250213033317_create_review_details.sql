-- +goose Up
-- +goose StatementBegin
CREATE TABLE "review_details" (
    "review_detail_id" SERIAL PRIMARY KEY,
    "review_id" INT NOT NULL REFERENCES "reviews" ("review_id") ON DELETE CASCADE,
    "type" VARCHAR(20) NOT NULL CHECK (type IN ('photo', 'video')),
    "url" TEXT NOT NULL,
    "caption" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_review_details_review_id" ON "review_details"("review_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_review_details_review_id";
DROP TABLE IF EXISTS "review_details";
-- +goose StatementEnd
