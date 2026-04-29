-- +goose Up
-- +goose StatementBegin
CREATE TABLE "sliders" (
    "slider_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "image" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

CREATE INDEX "idx_sliders_name" ON "sliders"("name");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_sliders_name";

DROP TABLE IF EXISTS "sliders";
-- +goose StatementEnd
