-- +goose Up
-- +goose StatementBegin
CREATE TABLE "reviews" (
    "review_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "name" VARCHAR(255) NOT NULL,
    "comment" TEXT NOT NULL,
    "rating" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_reviews_user_id" ON "reviews"("user_id");
CREATE INDEX "idx_reviews_product_id" ON "reviews"("product_id");
CREATE INDEX "idx_reviews_rating" ON "reviews"("rating");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_reviews_user_id";
DROP INDEX IF EXISTS "idx_reviews_product_id";
DROP INDEX IF EXISTS "idx_reviews_rating";
DROP TABLE IF EXISTS "reviews";
-- +goose StatementEnd
