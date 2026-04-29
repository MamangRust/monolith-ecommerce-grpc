-- GetShippingAddress: Retrieves paginated list of all shipping addresses (active and trashed) with search capability
-- Purpose: Display shipping addresses in admin dashboard
-- Parameters:
--   $1: search_term - Optional text to filter addresses by ID or alamat (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All shipping address fields plus total_count of matching records
-- Business Logic:
--   - Includes both active and trashed addresses
--   - Supports partial text matching on shipping_address_id (cast as text) and alamat fields (case-insensitive)
--   - Returns newest addresses first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetShippingAddress :many
SELECT
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM shipping_addresses
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR shipping_address_id::TEXT ILIKE '%' || $1 || '%'
        OR alamat ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetShippingAddressActive: Retrieves paginated list of active shipping addresses with search capability
-- Purpose: Display active addresses in checkout/address book UI
-- Parameters:
--   $1: search_term - Optional text to filter addresses by ID or alamat (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All shipping address fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted addresses (deleted_at IS NULL)
--   - Supports partial text matching on shipping_address_id (cast as text) and alamat fields (case-insensitive)
--   - Returns newest addresses first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetShippingAddressActive :many
SELECT
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM shipping_addresses
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR shipping_address_id::TEXT ILIKE '%' || $1 || '%'
        OR alamat ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetShippingAddressTrashed: Retrieves paginated list of trashed shipping addresses with search capability
-- Purpose: Display deleted addresses in admin recycle bin
-- Parameters:
--   $1: search_term - Optional text to filter addresses by ID or alamat (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All shipping address fields plus total_count of matching records
-- Business Logic:
--   - Only includes soft-deleted addresses (deleted_at IS NOT NULL)
--   - Supports partial text matching on shipping_address_id (cast as text) and alamat fields (case-insensitive)
--   - Returns newest addresses first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetShippingAddressTrashed :many
SELECT
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM shipping_addresses
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR shipping_address_id::TEXT ILIKE '%' || $1 || '%'
        OR alamat ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetShippingByID: Retrieves a single active shipping address by ID
-- Mengambil satu alamat pengiriman aktif berdasarkan ID
-- Parameters:
--   $1: shipping_address_id - ID of the shipping address to retrieve
--   $1: shipping_address_id - ID alamat pengiriman yang akan diambil
-- Returns:
--   Complete shipping address record if found and active
--   Seluruh record alamat pengiriman jika ditemukan dan aktif
-- Business Logic:
--   - Only returns non-deleted (active) addresses
--   - Hanya mengembalikan alamat yang tidak terhapus (aktif)
-- name: GetShippingByID :one
SELECT
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at
FROM shipping_addresses
WHERE
    shipping_address_id = $1
    AND deleted_at IS NULL;

-- GetShippingAddressByOrderID: Retrieves shipping address for a specific order
-- Mengambil alamat pengiriman untuk pesanan tertentu
-- Parameters:
--   $1: order_id - ID of the order to find shipping address for
--   $1: order_id - ID pesanan untuk mencari alamat pengiriman
-- Returns:
--   Complete shipping address record associated with the order
--   Seluruh record alamat pengiriman terkait pesanan
-- Business Logic:
--   - Used to display shipping info for completed orders
--   - Digunakan untuk menampilkan info pengiriman pesanan selesai
-- name: GetShippingAddressByOrderID :one
SELECT
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at
FROM shipping_addresses
WHERE
    order_id = $1
    AND deleted_at IS NULL;

-- CreateShippingAddress: Creates a new shipping address record
-- Membuat record alamat pengiriman baru
-- Parameters:
--   $1: order_id - Associated order ID
--   $2: alamat - Complete street address
--   $3: provinsi - Province/state
--   $4: negara - Country
--   $5: kota - City
--   $6: courier - Shipping courier name
--   $7: shipping_method - Shipping service type
--   $8: shipping_cost - Calculated shipping cost
-- Returns:
--   The complete created shipping address record
--   Seluruh record alamat pengiriman yang baru dibuat
-- Business Logic:
--   - Stores all necessary shipping information for an order
--   - Menyimpan semua informasi pengiriman yang diperlukan untuk pesanan
-- name: CreateShippingAddress :one
INSERT INTO
    shipping_addresses (
        order_id,
        alamat,
        provinsi,
        negara,
        kota,
        courier,
        shipping_method,
        shipping_cost
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    )
RETURNING
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at;

-- UpdateShippingAddress: Modifies an existing shipping address
-- Memperbarui alamat pengiriman yang sudah ada
-- Parameters:
--   $1: shipping_address_id - ID of address to update
--   $2-$8: Updated address fields (alamat to shipping_cost)
-- Returns:
--   The updated shipping address record
--   Record alamat pengiriman yang telah diperbarui
-- Business Logic:
--   - Only updates active (non-deleted) addresses
--   - Automatically updates the timestamp
--   - Hanya memperbarui alamat aktif (tidak terhapus)
--   - Secara otomatis memperbarui timestamp
-- name: UpdateShippingAddress :one
UPDATE shipping_addresses
SET
    alamat = $2,
    provinsi = $3,
    negara = $4,
    kota = $5,
    courier = $6,
    shipping_method = $7,
    shipping_cost = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE
    shipping_address_id = $1
    AND deleted_at IS NULL
RETURNING
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at;

-- TrashShippingAddress: Soft-deletes a shipping address
-- Menghapus sementara alamat pengiriman (soft delete)
-- Parameters:
--   $1: shipping_address_id - ID of address to trash
-- Returns:
--   The trashed shipping address record
--   Record alamat pengiriman yang telah dihapus sementara
-- Business Logic:
--   - Sets deleted_at timestamp for soft deletion
--   - Menandai deleted_at untuk penghapusan sementara
-- name: TrashShippingAddress :one
UPDATE shipping_addresses
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    shipping_address_id = $1
    AND deleted_at IS NULL
RETURNING
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at,
    deleted_at;

-- RestoreShippingAddress: Recovers a soft-deleted address
-- Memulihkan alamat pengiriman yang dihapus sementara
-- Parameters:
--   $1: shipping_address_id - ID of address to restore
-- Returns:
--   The restored shipping address record
--   Record alamat pengiriman yang telah dipulihkan
-- Business Logic:
--   - Clears the deleted_at field
--   - Membersihkan field deleted_at
-- name: RestoreShippingAddress :one
UPDATE shipping_addresses
SET
    deleted_at = NULL
WHERE
    shipping_address_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    shipping_address_id,
    order_id,
    alamat,
    provinsi,
    negara,
    kota,
    courier,
    shipping_method,
    shipping_cost,
    created_at,
    updated_at,
    deleted_at;

-- DeleteShippingAddressPermanently: Permanently removes a trashed address
-- Menghapus permanen alamat pengiriman yang sudah di-trash
-- Parameters:
--   $1: shipping_address_id - ID of address to delete
-- Returns:
--   Nothing (exec-only)
--   Tidak mengembalikan apa pun (exec-only)
-- Business Logic:
--   - Only works on already trashed addresses
--   - Irreversible deletion
--   - Hanya bekerja pada alamat yang sudah di-trash
--   - Penghapusan permanen tidak dapat dibatalkan
-- name: DeleteShippingAddressPermanently :exec
DELETE FROM shipping_addresses
WHERE
    shipping_address_id = $1
    AND deleted_at IS NOT NULL;

-- DeleteShippingAddressByOrderPermanent: Permanently removes a trashed address by OrderID
-- name: DeleteShippingAddressByOrderPermanent :exec
DELETE FROM shipping_addresses
WHERE
    order_id = $1;

-- RestoreAllShippingAddress: Recovers all trashed shipping addresses
-- Memulihkan semua alamat pengiriman yang dihapus sementara
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Admin-level bulk restore operation
--   - Operasi pemulihan massal level admin
-- name: RestoreAllShippingAddress :exec
UPDATE shipping_addresses
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentShippingAddress: Permanently removes all trashed addresses
-- Menghapus permanen semua alamat pengiriman yang di-trash
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Admin-level bulk deletion
--   - Irreversible operation
--   - Operasi penghapusan massal level admin
--   - Tidak dapat dibatalkan
-- name: DeleteAllPermanentShippingAddress :exec
DELETE FROM shipping_addresses
WHERE
    deleted_at IS NOT NULL
    OR order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL);