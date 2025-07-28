-- +goose Up
-- +goose StatementBegin
CREATE TABLE "banners" (
    "banner_id" SERIAL PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "start_date" DATE NOT NULL,       
    "end_date" DATE NOT NULL,          
    "start_time" TIME NOT NULL,     
    "end_time" TIME NOT NULL,          
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP WITH TIME ZONE NULL,
    CONSTRAINT "valid_date_range" CHECK (
        "end_date" IS NULL OR 
        "start_date" IS NULL OR 
        "end_date" >= "start_date"
    ),
    CONSTRAINT "valid_time_range" CHECK (
        "end_time" IS NULL OR 
        "start_time" IS NULL OR 
        "end_time" > "start_time"
    )
);

CREATE INDEX "idx_banners_active" ON "banners"("is_active") WHERE "deleted_at" IS NULL;
CREATE INDEX "idx_banners_date_range" ON "banners"("start_date", "end_date");
CREATE INDEX "idx_banners_time_range" ON "banners"("start_time", "end_time");
CREATE INDEX "idx_banners_created_at" ON "banners"("created_at");
CREATE INDEX "idx_banners_deleted_at" ON "banners"("deleted_at") WHERE "deleted_at" IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_banners_active";
DROP INDEX IF EXISTS "idx_banners_date_range";
DROP INDEX IF EXISTS "idx_banners_time_range";
DROP INDEX IF EXISTS "idx_banners_created_at";
DROP INDEX IF EXISTS "idx_banners_deleted_at";
DROP TABLE IF EXISTS "banners";
-- +goose StatementEnd