-- +goose Up
-- +goose StatementBegin
CREATE TABLE "merchant_documents" (
    "document_id" SERIAL PRIMARY KEY,
    "merchant_id" INT NOT NULL REFERENCES "merchants" ("merchant_id") ON DELETE CASCADE,
    "document_type" VARCHAR(50) NOT NULL,
    "document_url" TEXT NOT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'pending',
    "note" TEXT DEFAULT NULL, 
    "uploaded_at" TIMESTAMP DEFAULT current_timestamp,
    "created_at" TIMESTAMP DEFAULT current_timestamp,
    "updated_at" TIMESTAMP DEFAULT current_timestamp,
    "deleted_at" TIMESTAMP DEFAULT NULL
);


CREATE INDEX idx_merchant_documents_merchant_id ON merchant_documents (merchant_id);

CREATE INDEX idx_merchant_documents_document_type ON merchant_documents (document_type);

CREATE INDEX idx_merchant_documents_status ON merchant_documents (status);

CREATE INDEX idx_merchant_documents_merchant_type ON merchant_documents (merchant_id, document_type);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_merchant_documents_merchant_id;

DROP INDEX IF EXISTS idx_merchant_documents_document_type;

DROP INDEX IF EXISTS idx_merchant_documents_status;

DROP INDEX IF EXISTS idx_merchant_documents_merchant_type;

DROP TABLE IF EXISTS "merchant_documents";

-- +goose StatementEnd
