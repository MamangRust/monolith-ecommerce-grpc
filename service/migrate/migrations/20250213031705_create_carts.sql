-- +goose Up
-- +goose StatementBegin
CREATE TABLE "carts" (
    "cart_id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES "users" ("user_id"),
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "name" VARCHAR(255) NOT NULL,
    "price" INT NOT NULL,
    "image" VARCHAR(255) NOT NULL,
    "quantity" INT NOT NULL,
    "weight" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);
CREATE INDEX "idx_carts_user_id" ON "carts"("user_id");
CREATE INDEX "idx_carts_product_id" ON "carts"("product_id");
CREATE INDEX "idx_carts_name" ON "carts"(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_carts_user_id";
DROP INDEX IF EXISTS "idx_carts_product_id";
DROP INDEX IF EXISTS "idx_carts_name";
DROP TABLE IF EXISTS "carts";
-- +goose StatementEnd
