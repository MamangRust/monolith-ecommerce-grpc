-- +goose Up
-- +goose StatementBegin
CREATE TABLE "product_stock_adjustments" (
    "operation_id" VARCHAR(255) PRIMARY KEY,
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "delta" INT NOT NULL CHECK ("delta" <> 0),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "idx_product_stock_adjustments_product_id" ON "product_stock_adjustments"("product_id");

CREATE TABLE "order_stock_reservations" (
    "reservation_id" SERIAL PRIMARY KEY,
    "order_id" INT NOT NULL REFERENCES "orders" ("order_id"),
    "product_id" INT NOT NULL REFERENCES "products" ("product_id"),
    "quantity" INT NOT NULL CHECK ("quantity" > 0),
    "status" VARCHAR(20) NOT NULL DEFAULT 'reserved' CHECK ("status" IN ('reserved', 'released')),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "uq_order_stock_reservation" UNIQUE ("order_id", "product_id")
);
CREATE INDEX "idx_order_stock_reservations_order_id" ON "order_stock_reservations"("order_id");
CREATE INDEX "idx_order_stock_reservations_status" ON "order_stock_reservations"("status");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "order_stock_reservations";
DROP INDEX IF EXISTS "idx_product_stock_adjustments_product_id";
DROP TABLE IF EXISTS "product_stock_adjustments";
-- +goose StatementEnd
