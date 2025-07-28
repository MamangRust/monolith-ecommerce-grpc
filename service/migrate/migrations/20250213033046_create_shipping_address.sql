-- +goose Up
-- +goose StatementBegin
CREATE TABLE "shipping_addresses" (
    "shipping_address_id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL REFERENCES "orders" ("order_id"),
    "alamat" TEXT NOT NULL,
    "provinsi" VARCHAR(255) NOT NULL,
    "negara" VARCHAR(255) NOT NULL,
    "kota" VARCHAR(255) NOT NULL,
    "courier" VARCHAR(100) NOT NULL,
    "shipping_method" VARCHAR(255) NOT NULL,
    "shipping_cost" DECIMAL(10,2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_shipping_addresses_order_id" ON "shipping_addresses"("order_id");
CREATE INDEX "idx_shipping_addresses_provinsi" ON "shipping_addresses"("provinsi");
CREATE INDEX "idx_shipping_addresses_negara" ON "shipping_addresses"("negara");
CREATE INDEX "idx_shipping_addresses_kota" ON "shipping_addresses"("kota");
CREATE INDEX "idx_shipping_addresses_method" ON "shipping_addresses"("shipping_method");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_shipping_addresses_order_id";
DROP INDEX IF EXISTS "idx_shipping_addresses_provinsi";
DROP INDEX IF EXISTS "idx_shipping_addresses_negara";
DROP INDEX IF EXISTS "idx_shipping_addresses_kota";
DROP INDEX IF EXISTS "idx_shipping_addresses_method";
DROP TABLE IF EXISTS "shipping_addresses";
-- +goose StatementEnd
