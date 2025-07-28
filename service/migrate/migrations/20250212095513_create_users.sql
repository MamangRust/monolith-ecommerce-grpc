-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "users" (
    "user_id" serial PRIMARY KEY,
    "firstname" VARCHAR(100) NOT NULL,
    "lastname" varchar(100) NOT NULL,
    "email" varchar(100) UNIQUE NOT NULL,
    "password" varchar(100) NOT NULL,
    "verification_code" varchar(100) NOT NULL,
    "is_verified" bool DEFAULT false,
    "created_at" timestamp DEFAULT current_timestamp,
    "updated_at" timestamp DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);


CREATE INDEX "idx_users_email" ON "users" ("email");

CREATE INDEX "idx_users_firstname" ON "users" ("firstname");

CREATE INDEX "idx_users_lastname" ON "users" ("lastname");

CREATE INDEX "idx_users_firstname_lastname" ON "users" ("firstname", "lastname");

CREATE INDEX "idx_users_created_at" ON "users" ("created_at");
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS "idx_users_email";

DROP INDEX IF EXISTS "idx_users_firstname";

DROP INDEX IF EXISTS "idx_users_lastname";

DROP INDEX IF EXISTS "idx_users_firstname_lastname";

DROP INDEX IF EXISTS "idx_users_created_at";

DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
