-- GetSliders: Retrieves all sliders (active & trashed) with optional name search and pagination
-- Purpose: General listing of sliders regardless of status
-- Parameters:
--   $1: search_term - Optional text to filter sliders by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All slider fields plus total_count of matching records
-- Business Logic:
--   - Includes both active and trashed sliders
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest sliders first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetSliders :many
SELECT
    slider_id,
    name,
    image,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM sliders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetSlidersActive: Retrieves active sliders with optional name search and pagination
-- Purpose: Display active sliders for frontend/backend UI
-- Parameters:
--   $1: search_term - Optional text to filter sliders by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All slider fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted sliders (deleted_at IS NULL)
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest sliders first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetSlidersActive :many
SELECT
    slider_id,
    name,
    image,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM sliders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetSlidersTrashed: Retrieves trashed sliders with optional name search and pagination
-- Purpose: Display deleted sliders in admin recycle bin
-- Parameters:
--   $1: search_term - Optional text to filter sliders by name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All slider fields plus total_count of matching records
-- Business Logic:
--   - Only includes soft-deleted sliders (deleted_at IS NOT NULL)
--   - Supports partial text matching on name field (case-insensitive)
--   - Returns newest sliders first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetSlidersTrashed :many
SELECT
    slider_id,
    name,
    image,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM sliders
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- CreateSlider: Creates a new slider entry
-- Membuat slider baru
-- Parameters:
--   $1: name - Nama/judul slider
--   $2: image - URL gambar slider
-- Returns:
--   The complete created slider record
--   Seluruh record slider yang baru dibuat
-- Business Logic:
--   - Creates a new active slider
--   - Membuat slider aktif baru
--   - Requires both name and image
--   - Memerlukan nama dan gambar
-- name: CreateSlider :one
INSERT INTO
    sliders (name, image)
VALUES ($1, $2)
RETURNING
    slider_id,
    name,
    image,
    created_at,
    updated_at;

-- GetSliderByID: Retrieves a single active slider by ID
-- Mengambil satu slider aktif berdasarkan ID
-- Parameters:
--   $1: slider_id - ID slider yang akan diambil
-- Returns:
--   Complete slider record if found and active
--   Seluruh record slider jika ditemukan dan aktif
-- Business Logic:
--   - Only returns non-deleted (active) sliders
--   - Hanya mengembalikan slider yang tidak terhapus (aktif)
-- name: GetSliderByID :one
SELECT
    slider_id,
    name,
    image,
    created_at,
    updated_at
FROM sliders
WHERE
    slider_id = $1
    AND deleted_at IS NULL;

-- UpdateSlider: Modifies an existing slider
-- Memperbarui slider yang sudah ada
-- Parameters:
--   $1: slider_id - ID slider yang akan diperbarui
--   $2: name - Nama baru slider
--   $3: image - URL gambar baru slider
-- Returns:
--   The updated slider record
--   Record slider yang telah diperbarui
-- Business Logic:
--   - Updates both name and image
--   - Memperbarui nama dan gambar
--   - Only works on active sliders
--   - Hanya bekerja pada slider aktif
--   - Automatically updates timestamp
--   - Otomatis memperbarui timestamp
-- name: UpdateSlider :one
UPDATE sliders
SET
    name = $2,
    image = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    slider_id = $1
    AND deleted_at IS NULL
RETURNING
    slider_id,
    name,
    image,
    created_at,
    updated_at;

-- TrashSlider: Soft-deletes a slider
-- Menghapus sementara slider (soft delete)
-- Parameters:
--   $1: slider_id - ID slider yang akan dihapus
-- Returns:
--   The trashed slider record
--   Record slider yang telah dihapus sementara
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Menandai deleted_at timestamp
--   - Only works on active sliders
--   - Hanya bekerja pada slider aktif
-- name: TrashSlider :one
UPDATE sliders
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    slider_id = $1
    AND deleted_at IS NULL
RETURNING
    slider_id,
    name,
    image,
    created_at,
    updated_at,
    deleted_at;

-- RestoreSlider: Recovers a soft-deleted slider
-- Memulihkan slider yang dihapus sementara
-- Parameters:
--   $1: slider_id - ID slider yang akan dipulihkan
-- Returns:
--   The restored slider record
--   Record slider yang telah dipulihkan
-- Business Logic:
--   - Clears the deleted_at field
--   - Membersihkan field deleted_at
--   - Only works on trashed sliders
--   - Hanya bekerja pada slider yang dihapus sementara
-- name: RestoreSlider :one
UPDATE sliders
SET
    deleted_at = NULL
WHERE
    slider_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    slider_id,
    name,
    image,
    created_at,
    updated_at,
    deleted_at;

-- DeleteSliderPermanently: Permanently removes a trashed slider
-- Menghapus permanen slider yang sudah di-trash
-- Parameters:
--   $1: slider_id - ID slider yang akan dihapus permanen
-- Returns:
--   Nothing (exec-only)
--   Tidak mengembalikan apa pun (exec-only)
-- Business Logic:
--   - Physical deletion from database
--   - Penghapusan fisik dari database
--   - Only works on already trashed sliders
--   - Hanya bekerja pada slider yang sudah di-trash
--   - Irreversible operation
--   - Operasi tidak dapat dibatalkan
-- name: DeleteSliderPermanently :exec
DELETE FROM sliders WHERE slider_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllSliders: Recovers all trashed sliders
-- Memulihkan semua slider yang dihapus sementara
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Admin-level bulk restore operation
--   - Operasi pemulihan massal level admin
-- name: RestoreAllSliders :exec
UPDATE sliders SET deleted_at = NULL WHERE deleted_at IS NOT NULL;

-- DeleteAllPermanentSliders: Permanently removes all trashed sliders
-- Menghapus permanen semua slider yang di-trash
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Admin-level bulk deletion
--   - Operasi penghapusan massal level admin
--   - Irreversible operation
--   - Operasi tidak dapat dibatalkan
-- name: DeleteAllPermanentSliders :exec
DELETE FROM sliders WHERE deleted_at IS NOT NULL;