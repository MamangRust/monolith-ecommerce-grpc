-- name: GetMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetActiveMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NULL
    AND status != 'deleted'
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetTrashedMerchantDocuments :many
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchant_documents
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR document_type ILIKE '%' || $1 || '%'
        OR status ILIKE '%' || $1 || '%'
        OR note ILIKE '%' || $1 || '%'
    )
ORDER BY document_id
LIMIT $2
OFFSET
    $3;

-- name: GetMerchantDocument :one
SELECT
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at
FROM merchant_documents
WHERE
    document_id = $1
    AND deleted_at IS NULL;

-- name: CreateMerchantDocument :one
INSERT INTO
    merchant_documents (
        merchant_id,
        document_type,
        document_url,
        status,
        note,
        uploaded_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        current_timestamp,
        current_timestamp
    )
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at;

-- name: UpdateMerchantDocument :one
UPDATE merchant_documents
SET
    document_type = $2,
    document_url = $3,
    status = $4,
    note = $5,
    updated_at = current_timestamp
WHERE
    document_id = $1
    AND deleted_at IS NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at;

-- name: UpdateMerchantDocumentStatus :one
UPDATE merchant_documents
SET
    status = $2,
    note = $3,
    updated_at = current_timestamp
WHERE
    document_id = $1
    AND deleted_at IS NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at;

-- name: TrashMerchantDocument :one
UPDATE merchant_documents
SET
    deleted_at = current_timestamp
WHERE
    document_id = $1
    AND deleted_at IS NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at;

-- name: RestoreMerchantDocument :one
UPDATE merchant_documents
SET
    deleted_at = NULL
WHERE
    document_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    document_id,
    merchant_id,
    document_type,
    document_url,
    status,
    note,
    uploaded_at,
    created_at,
    updated_at,
    deleted_at;

-- name: DeleteMerchantDocumentPermanently :exec
DELETE FROM merchant_documents
WHERE
    document_id = $1
    AND deleted_at IS NOT NULL;

-- name: RestoreAllMerchantDocuments :exec
UPDATE merchant_documents
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- name: DeleteAllPermanentMerchantDocuments :exec
DELETE FROM merchant_documents WHERE deleted_at IS NOT NULL;