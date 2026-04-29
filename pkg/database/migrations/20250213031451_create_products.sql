-- +goose Up
-- +goose StatementBegin
CREATE TABLE "products" (
    "product_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id"),
    "category_id" INT NOT NULL REFERENCES "categories" ("category_id"),
    "name" VARCHAR(255) NOT NULL,
    "description" TEXT,
    "price" INT NOT NULL,
    "count_in_stock" INT NOT NULL DEFAULT 0,
    "brand" VARCHAR(100),
    "weight" INT,
    "rating" FLOAT DEFAULT 0.0 NULL,
    "slug_product" VARCHAR(100) UNIQUE,
    "image_product" TEXT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_products_merchant_id" ON "products" ("merchant_id");

CREATE INDEX "idx_products_category_id" ON "products" ("category_id");

CREATE INDEX "idx_products_slug" ON "products" ("slug_product");

CREATE INDEX "idx_products_price" ON "products" ("price");

CREATE INDEX "idx_products_created_at" ON "products" ("created_at");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_products_merchant_id" ON "products" ("merchant_id");

DROP INDEX IF EXISTS "idx_products_category_id" ON "products" ("category_id");

DROP INDEX IF EXISTS "idx_products_slug" ON "products" ("slug_product");

DROP INDEX IF EXISTS "idx_products_price" ON "products" ("price");

DROP INDEX IF EXISTS "idx_products_created_at" ON "products" ("created_at");

DROP TABLE IF EXISTS "products";
-- +goose StatementEnd
