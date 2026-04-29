-- GetOrders: Retrieves paginated list of active orders with search capability
-- Purpose: List all active orders for management UI
-- Parameters:
--   $1: search_term - Optional text to filter orders by ID or total price (NULL for no filter)
--   $2: limit - Maximum number of records to return (pagination limit)
--   $3: offset - Number of records to skip (pagination offset)
-- Returns:
--   All order fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted orders (deleted_at IS NULL)
--   - Supports partial text matching on order_id and total_price fields (case-insensitive)
--   - Returns newest orders first (created_at DESC)
--   - Provides total_count for client-side pagination
--   - Uses window function COUNT(*) OVER() for efficient total count
-- name: GetOrders :many
SELECT
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrdersActive: Retrieves paginated list of active orders (identical to GetOrders)
-- Purpose: Maintains consistent API pattern with other active/trashed endpoints
-- Parameters:
--   $1: search_term - Optional filter text for order ID or price
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Active order records with total_count
-- Business Logic:
--   - Same functionality as GetOrders
--   - Exists for consistency in API design patterns
-- Note: Could be consolidated with GetOrders if duplicate functionality is undesired
-- name: GetOrdersActive :many
SELECT
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrdersTrashed: Retrieves paginated list of soft-deleted orders
-- Purpose: View and manage deleted orders for potential restoration
-- Parameters:
--   $1: search_term - Optional text to filter trashed orders
--   $2: limit - Maximum records per page
--   $3: offset - Records to skip
-- Returns:
--   Trashed order records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Maintains same search functionality as active order queries
--   - Preserves chronological sorting (newest first)
--   - Used in order recovery/audit interfaces
--   - Includes total_count for pagination in trash management UI
-- name: GetOrdersTrashed :many
SELECT
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrdersByMerchant: Retrieves merchant-specific orders with pagination
-- Purpose: List orders filtered by merchant ID
-- Parameters:
--   $1: search_term - Optional text to filter orders
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
--   $4: merchant_id - Optional merchant UUID to filter by (NULL for all merchants)
-- Returns:
--   Order records with total_count
-- Business Logic:
--   - Combines merchant filtering with search functionality
--   - Maintains same sorting and pagination as other order queries
--   - Useful for merchant-specific order dashboards
--   - NULL merchant_id parameter returns all merchants' orders
-- name: GetOrdersByMerchant :many
SELECT
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM orders
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR total_price::TEXT ILIKE '%' || $1 || '%'
    )
    AND (
        $4::INT IS NULL
        OR merchant_id = $4
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetMonthlyTotalRevenue: Retrieves monthly total revenue across two custom date ranges
-- Purpose: Compare total revenue between two time periods (e.g., current month vs previous month)
-- Parameters:
--   $1: Start date of first period
--   $2: End date of first period
--   $3: Start date of second period
--   $4: End date of second period
-- Returns:
--   year: The year of the revenue data
--   month: The full month name (e.g., "January")
--   total_revenue: Total revenue (SUM of order totals) for that month (0 if no revenue)
-- Business Logic:
--   - Compares revenue between two customizable time periods
--   - Ensures all selected months appear even if no revenue (gap filling)
--   - Includes only non-deleted orders and order items
--   - Output formatted for charting or reporting tools
-- name: GetMonthlyTotalRevenue :many
WITH
    monthly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                (
                    o.created_at >= $1
                    AND o.created_at <= $2
                )
                OR (
                    o.created_at >= $3
                    AND o.created_at <= $4
                )
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            ),
            EXTRACT(
                MONTH
                FROM o.created_at
            )
    ),
    all_months AS (
        SELECT EXTRACT(
                YEAR
                FROM $1
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $1
            )::integer AS month, TO_CHAR($1, 'FMMonth') AS month_name
        UNION
        SELECT EXTRACT(
                YEAR
                FROM $3
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $3
            )::integer AS month, TO_CHAR($3, 'FMMonth') AS month_name
    )
SELECT COALESCE(
        am.year, EXTRACT(
            YEAR
            FROM $1
        )::TEXT
    ) AS year, COALESCE(
        am.month_name, TO_CHAR($1, 'FMMonth')
    ) AS month, COALESCE(mr.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_revenue mr ON am.year = mr.year
    AND am.month = mr.month
ORDER BY am.year DESC, am.month DESC;

-- GetYearlyTotalRevenue: Retrieves yearly total revenue for current and previous year
-- Purpose: Show year-over-year revenue trends
-- Parameters:
--   $1: The current year (integer)
-- Returns:
--   year: Year (as string)
--   total_revenue: Total revenue (SUM of order totals) for the year (0 if no revenue)
-- Business Logic:
--   - Automatically compares revenue between current and previous year
--   - Includes zero-value years for complete data visualization
--   - Filters only active/non-deleted orders and order items
-- name: GetYearlyTotalRevenue :many
WITH
    yearly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer - 1
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    ),
    all_years AS (
        SELECT $1 AS year
        UNION
        SELECT $1 - 1 AS year
    )
SELECT ay.year::text AS year, COALESCE(yr.total_revenue, 0) AS total_revenue
FROM
    all_years ay
    LEFT JOIN yearly_revenue yr ON ay.year = yr.year
ORDER BY ay.year DESC;

-- GetMonthlyTotalRevenueById: Retrieves monthly total revenue across two custom date ranges by order_id
-- Purpose: Compare total revenue between two time periods (e.g., current month vs previous month)
-- Parameters:
--   $1: Start date of first period
--   $2: End date of first period
--   $3: Start date of second period
--   $4: End date of second period
--   $5: Order ID
-- Returns:
--   year: The year of the revenue data
--   month: The full month name (e.g., "January")
--   total_revenue: Total revenue (SUM of order totals) for that month (0 if no revenue)
-- Business Logic:
--   - Compares revenue between two customizable time periods
--   - Ensures all selected months appear even if no revenue (gap filling)
--   - Includes only non-deleted orders and order items
--   - Output formatted for charting or reporting tools
-- name: GetMonthlyTotalRevenueById :many
WITH
    monthly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                (
                    o.created_at >= $1
                    AND o.created_at <= $2
                )
                OR (
                    o.created_at >= $3
                    AND o.created_at <= $4
                )
            )
            AND o.order_id = $5
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            ),
            EXTRACT(
                MONTH
                FROM o.created_at
            )
    ),
    all_months AS (
        SELECT EXTRACT(
                YEAR
                FROM $1
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $1
            )::integer AS month, TO_CHAR($1, 'FMMonth') AS month_name
        UNION
        SELECT EXTRACT(
                YEAR
                FROM $3
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $3
            )::integer AS month, TO_CHAR($3, 'FMMonth') AS month_name
    )
SELECT COALESCE(
        am.year, EXTRACT(
            YEAR
            FROM $1
        )::TEXT
    ) AS year, COALESCE(
        am.month_name, TO_CHAR($1, 'FMMonth')
    ) AS month, COALESCE(mr.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_revenue mr ON am.year = mr.year
    AND am.month = mr.month
ORDER BY am.year DESC, am.month DESC;

-- GetYearlyTotalRevenueById: Retrieves yearly total revenue for current and previous year by order_id
-- Purpose: Show year-over-year revenue trends
-- Parameters:
--   $1: The current year (integer)
--   $2: Order ID
-- Returns:
--   year: Year (as string)
--   total_revenue: Total revenue (SUM of order totals) for the year (0 if no revenue)
-- Business Logic:
--   - Automatically compares revenue between current and previous year
--   - Includes zero-value years for complete data visualization
--   - Filters only active/non-deleted orders and order items
-- name: GetYearlyTotalRevenueById :many
WITH
    yearly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer - 1
            )
            AND o.order_id = $2
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    ),
    all_years AS (
        SELECT $1 AS year
        UNION
        SELECT $1 - 1 AS year
    )
SELECT ay.year::text AS year, COALESCE(yr.total_revenue, 0) AS total_revenue
FROM
    all_years ay
    LEFT JOIN yearly_revenue yr ON ay.year = yr.year
ORDER BY ay.year DESC;

-- GetMonthlyTotalRevenueByMerchant: Retrieves monthly total revenue across two custom date ranges by merchant_id
-- Purpose: Compare total revenue between two time periods (e.g., current month vs previous month)
-- Parameters:
--   $1: Start date of first period
--   $2: End date of first period
--   $3: Start date of second period
--   $4: End date of second period
--   $5: Order ID
-- Returns:
--   year: The year of the revenue data
--   month: The full month name (e.g., "January")
--   total_revenue: Total revenue (SUM of order totals) for that month (0 if no revenue)
-- Business Logic:
--   - Compares revenue between two customizable time periods
--   - Ensures all selected months appear even if no revenue (gap filling)
--   - Includes only non-deleted orders and order items
--   - Output formatted for charting or reporting tools
-- name: GetMonthlyTotalRevenueByMerchant :many
WITH
    monthly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                (
                    o.created_at >= $1
                    AND o.created_at <= $2
                )
                OR (
                    o.created_at >= $3
                    AND o.created_at <= $4
                )
            )
            AND o.merchant_id = $5
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            ),
            EXTRACT(
                MONTH
                FROM o.created_at
            )
    ),
    all_months AS (
        SELECT EXTRACT(
                YEAR
                FROM $1
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $1
            )::integer AS month, TO_CHAR($1, 'FMMonth') AS month_name
        UNION
        SELECT EXTRACT(
                YEAR
                FROM $3
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM $3
            )::integer AS month, TO_CHAR($3, 'FMMonth') AS month_name
    )
SELECT COALESCE(
        am.year, EXTRACT(
            YEAR
            FROM $1
        )::TEXT
    ) AS year, COALESCE(
        am.month_name, TO_CHAR($1, 'FMMonth')
    ) AS month, COALESCE(mr.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_revenue mr ON am.year = mr.year
    AND am.month = mr.month
ORDER BY am.year DESC, am.month DESC;

-- GetYearlyTotalRevenueByMerchant: Retrieves yearly total revenue for current and previous year by merchant_id
-- Purpose: Show year-over-year revenue trends
-- Parameters:
--   $1: The current year (integer)
--   $2: Order ID
-- Returns:
--   year: Year (as string)
--   total_revenue: Total revenue (SUM of order totals) for the year (0 if no revenue)
-- Business Logic:
--   - Automatically compares revenue between current and previous year
--   - Includes zero-value years for complete data visualization
--   - Filters only active/non-deleted orders and order items
-- name: GetYearlyTotalRevenueByMerchant :many
WITH
    yearly_revenue AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND (
                EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer
                OR EXTRACT(
                    YEAR
                    FROM o.created_at
                ) = $1::integer - 1
            )
            AND o.merchant_id = $2
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    ),
    all_years AS (
        SELECT $1 AS year
        UNION
        SELECT $1 - 1 AS year
    )
SELECT ay.year::text AS year, COALESCE(yr.total_revenue, 0) AS total_revenue
FROM
    all_years ay
    LEFT JOIN yearly_revenue yr ON ay.year = yr.year
ORDER BY ay.year DESC;

-- GetMonthlyOrder: Retrieves monthly order summary within a 1-year period
-- Purpose: Provides monthly sales performance metrics for trend and operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis window
-- Returns:
--   month: 3-letter abbreviation of the month (e.g., 'Jan')
--   order_count: Total number of orders in the month
--   total_revenue: Sum of total price from all orders
--   total_items_sold: Total quantity of items sold in that month
-- Business Logic:
--   - Analyzes a 12-month period starting from the month of the reference date
--   - Ignores soft-deleted records for accurate reporting
--   - Aggregates data by month for visualizations and monthly performance tracking
--   - Uses short month format for dashboard/chart compactness
--   - Sorts chronologically by month
-- name: GetMonthlyOrder :many
WITH
    date_range AS (
        SELECT date_trunc('month', $1::timestamp) AS start_date, date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
    ),
    monthly_orders AS (
        SELECT
            date_trunc('month', o.created_at) AS activity_month,
            COUNT(o.order_id) AS order_count,
            SUM(o.total_price)::NUMERIC AS total_revenue,
            SUM(oi.quantity) AS total_items_sold
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND o.created_at BETWEEN (
                SELECT start_date
                FROM date_range
            ) AND (
                SELECT end_date
                FROM date_range
            )
        GROUP BY
            activity_month
    )
SELECT TO_CHAR(mo.activity_month, 'Mon') AS month, mo.order_count, mo.total_revenue, mo.total_items_sold
FROM monthly_orders mo
ORDER BY mo.activity_month;

-- GetYearlyOrder: Retrieves yearly order summary over the past 5 years
-- Purpose: Enables long-term trend analysis of sales performance
-- Parameters:
--   $1: Reference date (timestamp) - defines the 5-year analysis window
-- Returns:
--   year: 4-digit year as string
--   order_count: Total number of orders in the year
--   total_revenue: Sum of total price from all orders
--   total_items_sold: Total quantity of items sold in the year
--   active_cashiers: Number of distinct cashier IDs involved in transactions
--   unique_products_sold: Number of unique products sold
-- Business Logic:
--   - Covers a rolling 5-year window up to the reference year
--   - Filters out deleted records to ensure data consistency
--   - Useful for high-level KPI tracking, forecasting, and strategic planning
--   - Includes both volume and revenue metrics for comprehensive reporting
--   - Results sorted by year in ascending order
-- name: GetYearlyOrder :many
WITH
    last_five_years AS (
        SELECT
            EXTRACT(
                YEAR
                FROM o.created_at
            )::text AS year,
            COUNT(o.order_id) AS order_count,
            SUM(o.total_price)::NUMERIC AS total_revenue,
            SUM(oi.quantity) AS total_items_sold,
            COUNT(DISTINCT oi.product_id) AS unique_products_sold
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND EXTRACT(
                YEAR
                FROM o.created_at
            ) BETWEEN (
                EXTRACT(
                    YEAR
                    FROM $1::timestamp
                ) - 4
            ) AND EXTRACT(
                YEAR
                FROM $1::timestamp
            )
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    )
SELECT
    year,
    order_count,
    total_revenue,
    total_items_sold,
    unique_products_sold
FROM last_five_years
ORDER BY year;

-- GetMonthlyOrderByMerchant: Retrieves monthly order summary within a 1-year period by merchant_id
-- Purpose: Provides monthly sales performance metrics for trend and operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis window
--   $2: Merchant ID
-- Returns:
--   month: 3-letter abbreviation of the month (e.g., 'Jan')
--   order_count: Total number of orders in the month
--   total_revenue: Sum of total price from all orders
--   total_items_sold: Total quantity of items sold in that month
-- Business Logic:
--   - Analyzes a 12-month period starting from the month of the reference date
--   - Ignores soft-deleted records for accurate reporting
--   - Aggregates data by month for visualizations and monthly performance tracking
--   - Uses short month format for dashboard/chart compactness
--   - Sorts chronologically by month
-- name: GetMonthlyOrderByMerchant :many
WITH
    date_range AS (
        SELECT date_trunc('month', $1::timestamp) AS start_date, date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
    ),
    monthly_orders AS (
        SELECT
            date_trunc('month', o.created_at) AS activity_month,
            COUNT(o.order_id) AS order_count,
            SUM(o.total_price)::NUMERIC AS total_revenue,
            SUM(oi.quantity) AS total_items_sold
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND o.created_at BETWEEN (
                SELECT start_date
                FROM date_range
            ) AND (
                SELECT end_date
                FROM date_range
            )
            AND o.merchant_id = $2
        GROUP BY
            activity_month
    )
SELECT TO_CHAR(mo.activity_month, 'Mon') AS month, mo.order_count, mo.total_revenue, mo.total_items_sold
FROM monthly_orders mo
ORDER BY mo.activity_month;

-- GetYearlyOrderByMerchant: Retrieves yearly order summary over the past 5 years by merchant_id
-- Purpose: Enables long-term trend analysis of sales performance
-- Parameters:
--   $1: Reference date (timestamp) - defines the 5-year analysis window
-- Returns:
--   year: 4-digit year as string
--   order_count: Total number of orders in the year
--   total_revenue: Sum of total price from all orders
--   total_items_sold: Total quantity of items sold in the year
--   active_cashiers: Number of distinct cashier IDs involved in transactions
--   unique_products_sold: Number of unique products sold
-- Business Logic:
--   - Covers a rolling 5-year window up to the reference year
--   - Filters out deleted records to ensure data consistency
--   - Useful for high-level KPI tracking, forecasting, and strategic planning
--   - Includes both volume and revenue metrics for comprehensive reporting
--   - Results sorted by year in ascending order
-- name: GetYearlyOrderByMerchant :many
WITH
    last_five_years AS (
        SELECT
            EXTRACT(
                YEAR
                FROM o.created_at
            )::text AS year,
            COUNT(o.order_id) AS order_count,
            SUM(o.total_price)::NUMERIC AS total_revenue,
            SUM(oi.quantity) AS total_items_sold,
            COUNT(DISTINCT oi.product_id) AS unique_products_sold
        FROM orders o
            JOIN order_items oi ON o.order_id = oi.order_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND EXTRACT(
                YEAR
                FROM o.created_at
            ) BETWEEN (
                EXTRACT(
                    YEAR
                    FROM $1::timestamp
                ) - 4
            ) AND EXTRACT(
                YEAR
                FROM $1::timestamp
            )
            AND o.merchant_id = $2
        GROUP BY
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    )
SELECT
    year,
    order_count,
    total_revenue,
    total_items_sold,
    unique_products_sold
FROM last_five_years
ORDER BY year;

-- CreateOrder: Creates a new order record
-- Purpose: Register a new transaction in the system
-- Parameters:
--   $1: merchant_id - UUID of the merchant associated with the order
--   $2: user_id - ID of the user processing the order
--   $3: total_price - Numeric total amount of the order
-- Returns: The complete created order record
-- Business Logic:
--   - Automatically sets created_at timestamp
--   - Requires merchant_id, user_id and total_price
--   - Typically followed by order item creation
-- name: CreateOrder :one
INSERT INTO
    orders (
        merchant_id,
        user_id,
        total_price
    )
VALUES ($1, $2, $3)
RETURNING
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at;

-- GetOrderByID: Retrieves an active order by ID
-- Purpose: Fetch order details for display/processing
-- Parameters:
--   $1: order_id - UUID of the order to retrieve
-- Returns: Full order record if found and active
-- Business Logic:
--   - Excludes soft-deleted orders
--   - Used for order viewing, receipts, and processing
--   - Typically joined with order_items in application
-- name: GetOrderByID :one
SELECT
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at
FROM orders
WHERE
    order_id = $1
    AND deleted_at IS NULL;

-- UpdateOrder: Modifies order information
-- Purpose: Update order details (primarily total price)
-- Parameters:
--   $1: order_id - UUID of order to update
--   $2: total_price - New total amount
-- Returns: Updated order record
-- Business Logic:
--   - Auto-updates updated_at timestamp
--   - Only modifies active (non-deleted) orders
--   - Used when order items change
--   - Should trigger recalculation of total_price
-- name: UpdateOrder :one
UPDATE orders
SET
    total_price = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    order_id = $1
    AND deleted_at IS NULL
RETURNING
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at;

-- TrashedOrder: Soft-deletes an order
-- Purpose: Cancel/void an order without permanent deletion
-- Parameters:
--   $1: order_id - UUID of order to cancel
-- Returns: The soft-deleted order record
-- Business Logic:
--   - Sets deleted_at to current timestamp
--   - Preserves order data for reporting
--   - Only processes active orders
--   - Can be restored via RestoreOrder
-- name: TrashedOrder :one
UPDATE orders
SET
    deleted_at = current_timestamp
WHERE
    order_id = $1
    AND deleted_at IS NULL
RETURNING
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    deleted_at;

-- RestoreOrder: Recovers a soft-deleted order
-- Purpose: Reactivate a cancelled order
-- Parameters:
--   $1: order_id - UUID of order to restore
-- Returns: The restored order record
-- Business Logic:
--   - Nullifies deleted_at field
--   - Only works on previously cancelled orders
--   - Maintains all original order data
-- name: RestoreOrder :one
UPDATE orders
SET
    deleted_at = NULL
WHERE
    order_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    order_id,
    user_id,
    merchant_id,
    total_price,
    created_at,
    updated_at,
    deleted_at;

-- DeleteOrderPermanently: Hard-deletes an order
-- Purpose: Completely remove order from database
-- Parameters:
--   $1: order_id - UUID of order to delete
-- Business Logic:
--   - Permanent deletion of already cancelled orders
--   - No return value (exec-only operation)
--   - Irreversible action - use with caution
--   - Should trigger deletion of related order_items
-- name: DeleteOrderPermanently :exec
DELETE FROM orders WHERE order_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllOrders: Mass restoration of cancelled orders
-- Purpose: Recover all trashed orders at once
-- Business Logic:
--   - Reactivates all soft-deleted orders
--   - No parameters needed (bulk operation)
--   - Typically used during system recovery
-- name: RestoreAllOrders :exec
UPDATE orders
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentOrders: Purges all cancelled orders
-- Purpose: Clean up all soft-deleted order records
-- Business Logic:
--   - Irreversible bulk deletion operation
--   - Only affects already cancelled orders
--   - Typically used during database maintenance
--   - Should be restricted to admin users
-- name: DeleteAllPermanentOrders :exec
DELETE FROM orders WHERE deleted_at IS NOT NULL;