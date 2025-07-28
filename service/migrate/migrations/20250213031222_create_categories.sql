-- +goose Up
-- +goose StatementBegin
CREATE TABLE "categories" (
    "category_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "description" TEXT,
    "slug_category" VARCHAR(100) UNIQUE,
    "image_category" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_categories_slug" ON "categories" ("slug_category");
CREATE INDEX "idx_categories_created_at" ON "categories" ("created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_categories_slug" ON "categories" ("slug_category");
DROP INDEX IF EXISTS "idx_categories_created_at" ON "categories" ("created_at");
DROP TABLE IF EXISTS "categories";
-- +goose StatementEnd
