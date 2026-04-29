-- GetCategories: Retrieves paginated list of active categories with search capability
-- Purpose: List all active product categories for management UI
-- Parameters:
--   $1: search_term - Optional text to filter categories by name or slug (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All category fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted categories (deleted_at IS NULL)
--   - Supports partial text matching on name and slug_category fields (case-insensitive)
--   - Returns newest categories first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCategories :many
SELECT
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM categories
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR slug_category ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetCategoriesActive: Retrieves paginated list of active categories with search capability
-- Purpose: List all active product categories for management UI
-- Parameters:
--   $1: search_term - Optional text to filter categories by name or slug (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All category fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted categories (deleted_at IS NULL)
--   - Supports partial text matching on name and slug_category fields (case-insensitive)
--   - Returns newest categories first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetCategoriesActive :many
SELECT
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM categories
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR slug_category ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetCategoriesTrashed: Retrieves paginated list of soft-deleted categories
-- Purpose: View/manage deleted categories for potential restoration
-- Parameters:
--   $1: search_term - Optional filter text (NULL for all trashed categories)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Trashed category records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Same search functionality as active categories
--   - Maintains consistent sorting with active records
--   - Used in trash management/recovery interfaces
-- name: GetCategoriesTrashed :many
SELECT
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM categories
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR slug_category ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetMonthlyTotalPrice: Retrieves monthly revenue totals across two comparison periods
-- Purpose: Provides month-over-month revenue analytics for financial reporting
-- Parameters:
--   $1: Start date of first comparison period
--   $2: End date of first comparison period
--   $3: Start date of second comparison period
--   $4: End date of second comparison period
-- Returns:
--   year: Year of revenue data (text format)
--   month_name: Full month name (e.g. "January")
--   total_revenue: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares revenue between two customizable date ranges
--   - Joins with order_items to ensure accurate order calculations
--   - Excludes deleted orders and order items for data integrity
--   - Uses gap-filling to show all months in both periods
--   - Formats output for financial dashboards
-- name: GetMonthlyTotalPrice :many
WITH
    monthly_totals AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
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
    ) AS month, COALESCE(mt.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_totals mt ON am.year = mt.year
    AND am.month = mt.month
ORDER BY am.year::INT DESC, am.month DESC;

-- GetYearlyTotalPrice: Retrieves annual revenue with category/product validation
-- Purpose: Provides year-over-year revenue analysis with product hierarchy verification
-- Parameters:
--   $1: Reference year for comparison (current year)
-- Returns:
--   year: Year as text
--   total_revenue: Annual revenue total (0 if no sales)
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Validates product/category relationships through joins
--   - Excludes deleted records across all joined tables
--   - Ensures complete year reporting even with no sales
--   - Orders results by most recent year first
-- name: GetYearlyTotalPrice :many
WITH
    yearly_data AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
SELECT a.year::text AS year, COALESCE(yd.total_revenue, 0) AS total_revenue
FROM all_years a
    LEFT JOIN yearly_data yd ON a.year = yd.year
ORDER BY a.year DESC;

-- GetMonthlyTotalPriceByMerchant: Retrieves monthly revenue totals across two comparison periods by merchant_id
-- Purpose: Provides month-over-month revenue analytics for financial reporting
-- Parameters:
--   $1: Start date of first comparison period
--   $2: End date of first comparison period
--   $3: Start date of second comparison period
--   $4: End date of second comparison period
--   $5: Merchant ID
-- Returns:
--   year: Year of revenue data (text format)
--   month_name: Full month name (e.g. "January")
--   total_revenue: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares revenue between two customizable date ranges
--   - Joins with order_items to ensure accurate order calculations
--   - Excludes deleted orders and order items for data integrity
--   - Uses gap-filling to show all months in both periods
--   - Formats output for financial dashboards
-- name: GetMonthlyTotalPriceByMerchant :many
WITH
    monthly_totals AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
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
    ) AS month, COALESCE(mt.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_totals mt ON am.year = mt.year
    AND am.month = mt.month
ORDER BY am.year::INT DESC, am.month DESC;

-- GetYearlyTotalPriceByMerchant: Retrieves annual revenue with category/product validation by merchant_id
-- Purpose: Provides year-over-year revenue analysis with product hierarchy verification
-- Parameters:
--   $1: Reference year for comparison (current year)
--   $2: Merchant ID
-- Returns:
--   year: Year as text
--   total_revenue: Annual revenue total (0 if no sales)
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Validates product/category relationships through joins
--   - Excludes deleted records across all joined tables
--   - Ensures complete year reporting even with no sales
--   - Orders results by most recent year first
-- name: GetYearlyTotalPriceByMerchant :many
WITH
    yearly_data AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
SELECT a.year::text AS year, COALESCE(yd.total_revenue, 0) AS total_revenue
FROM all_years a
    LEFT JOIN yearly_data yd ON a.year = yd.year
ORDER BY a.year DESC;

-- GetMonthlyTotalPriceById: Retrieves monthly revenue totals across two comparison periods by category_id
-- Purpose: Provides month-over-month revenue analytics for financial reporting
-- Parameters:
--   $1: Start date of first comparison period
--   $2: End date of first comparison period
--   $3: Start date of second comparison period
--   $4: End date of second comparison period
--   $5: Category ID
-- Returns:
--   year: Year of revenue data (text format)
--   month_name: Full month name (e.g. "January")
--   total_revenue: Sum of order totals for that month (0 if no sales)
-- Business Logic:
--   - Compares revenue between two customizable date ranges
--   - Joins with order_items to ensure accurate order calculations
--   - Excludes deleted orders and order items for data integrity
--   - Uses gap-filling to show all months in both periods
--   - Formats output for financial dashboards
-- name: GetMonthlyTotalPriceById :many
WITH
    monthly_totals AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::TEXT AS year, EXTRACT(
                MONTH
                FROM o.created_at
            )::integer AS month, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
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
            AND c.category_id = $5
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
    ) AS month, COALESCE(mt.total_revenue, 0) AS total_revenue
FROM
    all_months am
    LEFT JOIN monthly_totals mt ON am.year = mt.year
    AND am.month = mt.month
ORDER BY am.year::INT DESC, am.month DESC;

-- GetYearlyTotalPriceById: Retrieves annual revenue with category/product validation by category_id
-- Purpose: Provides year-over-year revenue analysis with product hierarchy verification
-- Parameters:
--   $1: Reference year for comparison (current year)
--   $2: Category ID
-- Returns:
--   year: Year as text
--   total_revenue: Annual revenue total (0 if no sales)
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Validates product/category relationships through joins
--   - Excludes deleted records across all joined tables
--   - Ensures complete year reporting even with no sales
--   - Orders results by most recent year first
-- name: GetYearlyTotalPriceById :many
WITH
    yearly_data AS (
        SELECT EXTRACT(
                YEAR
                FROM o.created_at
            )::integer AS year, COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
            AND c.category_id = $2
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
SELECT a.year::text AS year, COALESCE(yd.total_revenue, 0) AS total_revenue
FROM all_years a
    LEFT JOIN yearly_data yd ON a.year = yd.year
ORDER BY a.year DESC;

-- GetMonthlyCategory: Retrieves monthly sales activity for all categories within a 1-year period
-- Purpose: Provides category performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   category_id: Unique identifier for the category
--   category_name: Display name of the category
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders associated with the category
--   items_sold: Total quantity of items sold from the category
--   total_revenue: Total revenue generated from category items
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted orders, items, products, and categories to ensure valid data
--   - Aggregates by category and month for trend tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Results sorted by month and revenue for time-series analysis
-- name: GetMonthlyCategory :many
WITH
    date_range AS (
        SELECT date_trunc('month', $1::timestamp) AS start_date, date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
    ),
    monthly_category_stats AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            date_trunc('month', o.created_at) AS activity_month,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
            AND o.created_at BETWEEN (
                SELECT start_date
                FROM date_range
            ) AND (
                SELECT end_date
                FROM date_range
            )
        GROUP BY
            c.category_id,
            c.name,
            activity_month
    )
SELECT TO_CHAR(mcs.activity_month, 'Mon') AS month, mcs.category_id, mcs.category_name, mcs.order_count, mcs.items_sold, mcs.total_revenue
FROM monthly_category_stats mcs
ORDER BY mcs.activity_month, mcs.total_revenue DESC;

-- GetYearlyCategory: Retrieves annual sales performance for categories over a 5-year span
-- Purpose: Enables long-term product category performance trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   category_id: Unique category identifier
--   category_name: Display name of the category
--   order_count: Annual number of orders involving this category
--   items_sold: Quantity of products sold from this category
--   total_revenue: Total sales revenue from category products
--   unique_products_sold: Count of unique products sold within the category
-- Business Logic:
--   - Covers the current year and previous four years (5-year window)
--   - Filters out soft-deleted data from all related tables
--   - Provides both volume and value metrics for category-level evaluation
--   - Results sorted by year and revenue to show historical trends
--   - Suitable for business reviews and strategic category planning
-- name: GetYearlyCategory :many
WITH
    last_five_years AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )::text AS year,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue,
            COUNT(DISTINCT oi.product_id) AS unique_products_sold
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
            c.category_id,
            c.name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    )
SELECT
    year,
    category_id,
    category_name,
    order_count,
    items_sold,
    total_revenue,
    unique_products_sold
FROM last_five_years
ORDER BY year, total_revenue DESC;

-- GetMonthlyCategoryByMerchant: Retrieves monthly sales activity for all categories within a 1-year period by merchant_id
-- Purpose: Provides category performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
--   $2: Merchant ID
-- Returns:
--   category_id: Unique identifier for the category
--   category_name: Display name of the category
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders associated with the category
--   items_sold: Total quantity of items sold from the category
--   total_revenue: Total revenue generated from category items
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted orders, items, products, and categories to ensure valid data
--   - Aggregates by category and month for trend tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Results sorted by month and revenue for time-series analysis
-- name: GetMonthlyCategoryByMerchant :many
WITH
    date_range AS (
        SELECT date_trunc('month', $1::timestamp) AS start_date, date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
    ),
    monthly_category_stats AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            date_trunc('month', o.created_at) AS activity_month,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
            AND o.created_at BETWEEN (
                SELECT start_date
                FROM date_range
            ) AND (
                SELECT end_date
                FROM date_range
            )
            AND o.merchant_id = $2
        GROUP BY
            c.category_id,
            c.name,
            activity_month
    )
SELECT TO_CHAR(mcs.activity_month, 'Mon') AS month, mcs.category_id, mcs.category_name, mcs.order_count, mcs.items_sold, mcs.total_revenue
FROM monthly_category_stats mcs
ORDER BY mcs.activity_month, mcs.total_revenue DESC;

-- GetYearlyCategoryByMerchant: Retrieves annual sales performance for categories over a 5-year span by merchant_id
-- Purpose: Enables long-term product category performance trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
--   $2: Merchant ID
-- Returns:
--   year: 4-digit year as text
--   category_id: Unique category identifier
--   category_name: Display name of the category
--   order_count: Annual number of orders involving this category
--   items_sold: Quantity of products sold from this category
--   total_revenue: Total sales revenue from category products
--   unique_products_sold: Count of unique products sold within the category
-- Business Logic:
--   - Covers the current year and previous four years (5-year window)
--   - Filters out soft-deleted data from all related tables
--   - Provides both volume and value metrics for category-level evaluation
--   - Results sorted by year and revenue to show historical trends
--   - Suitable for business reviews and strategic category planning
-- name: GetYearlyCategoryByMerchant :many
WITH
    last_five_years AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )::text AS year,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue,
            COUNT(DISTINCT oi.product_id) AS unique_products_sold
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
            c.category_id,
            c.name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    )
SELECT
    year,
    category_id,
    category_name,
    order_count,
    items_sold,
    total_revenue,
    unique_products_sold
FROM last_five_years
ORDER BY year, total_revenue DESC;

-- GetMonthlyCategoryById: Retrieves monthly sales activity for all categories within a 1-year period by category_id
-- Purpose: Provides category performance metrics by month for operational analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
--   $2: Category ID
-- Returns:
--   category_id: Unique identifier for the category
--   category_name: Display name of the category
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   order_count: Number of orders associated with the category
--   items_sold: Total quantity of items sold from the category
--   total_revenue: Total revenue generated from category items
-- Business Logic:
--   - Analyzes a rolling 12-month period from the reference date
--   - Excludes deleted orders, items, products, and categories to ensure valid data
--   - Aggregates by category and month for trend tracking
--   - Uses abbreviated month names for compact visual reporting
--   - Results sorted by month and revenue for time-series analysis
-- name: GetMonthlyCategoryById :many
WITH
    date_range AS (
        SELECT date_trunc('month', $1::timestamp) AS start_date, date_trunc('month', $1::timestamp) + interval '1 year' - interval '1 day' AS end_date
    ),
    monthly_category_stats AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            date_trunc('month', o.created_at) AS activity_month,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
            AND o.created_at BETWEEN (
                SELECT start_date
                FROM date_range
            ) AND (
                SELECT end_date
                FROM date_range
            )
            AND c.category_id = $2
        GROUP BY
            c.category_id,
            c.name,
            activity_month
    )
SELECT TO_CHAR(mcs.activity_month, 'Mon') AS month, mcs.category_id, mcs.category_name, mcs.order_count, mcs.items_sold, mcs.total_revenue
FROM monthly_category_stats mcs
ORDER BY mcs.activity_month, mcs.total_revenue DESC;

-- GetYearlyCategoryById: Retrieves annual sales performance for categories over a 5-year span by category_id
-- Purpose: Enables long-term product category performance trend analysis
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
--   $2: Category ID
-- Returns:
--   year: 4-digit year as text
--   category_id: Unique category identifier
--   category_name: Display name of the category
--   order_count: Annual number of orders involving this category
--   items_sold: Quantity of products sold from this category
--   total_revenue: Total sales revenue from category products
--   unique_products_sold: Count of unique products sold within the category
-- Business Logic:
--   - Covers the current year and previous four years (5-year window)
--   - Filters out soft-deleted data from all related tables
--   - Provides both volume and value metrics for category-level evaluation
--   - Results sorted by year and revenue to show historical trends
--   - Suitable for business reviews and strategic category planning
-- name: GetYearlyCategoryById :many
WITH
    last_five_years AS (
        SELECT
            c.category_id,
            c.name AS category_name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )::text AS year,
            COUNT(DISTINCT o.order_id) AS order_count,
            SUM(oi.quantity) AS items_sold,
            COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue,
            COUNT(DISTINCT oi.product_id) AS unique_products_sold
        FROM
            orders o
            JOIN order_items oi ON o.order_id = oi.order_id
            JOIN products p ON oi.product_id = p.product_id
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            o.deleted_at IS NULL
            AND oi.deleted_at IS NULL
            AND p.deleted_at IS NULL
            AND c.deleted_at IS NULL
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
            AND c.category_id = $2
        GROUP BY
            c.category_id,
            c.name,
            EXTRACT(
                YEAR
                FROM o.created_at
            )
    )
SELECT
    year,
    category_id,
    category_name,
    order_count,
    items_sold,
    total_revenue,
    unique_products_sold
FROM last_five_years
ORDER BY year, total_revenue DESC;

-- CreateCategory: Inserts a new category into the system
-- Purpose: Adds a new product category for classification and reporting
-- Parameters:
--   $1: Category name
--   $2: Category description
--   $3: Slug for category (URL-friendly identifier)
-- Returns:
--   Full category record including generated ID
-- Business Logic:
--   - Assumes unique slug for identification in URLs
--   - Automatically populates timestamps via default DB behavior (if configured)
-- name: CreateCategory :one
INSERT INTO
    categories (
        name,
        description,
        slug_category,
        image_category
    )
VALUES ($1, $2, $3, $4)
RETURNING
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at;

-- GetCategoryByID: Fetches a single category by its ID
-- Purpose: Retrieve details of an active (non-deleted) category
-- Parameters:
--   $1: Category ID to search for
-- Returns:
--   Full category record if found and not deleted
-- Business Logic:
--   - Excludes soft-deleted categories
-- name: GetCategoryByID :one
SELECT
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at
FROM categories
WHERE
    category_id = $1
    AND deleted_at IS NULL;

-- GetCategoryByIDTrashed: Fetches a single category by its ID
-- Purpose: Retrieve details of an active (non-deleted) category
-- Parameters:
--   $1: Category ID to search for
-- Returns:
--   Full category record if found and not deleted
-- Business Logic:
--   - Excludes soft-deleted categories
-- name: GetCategoryByIDTrashed :one
SELECT *
FROM categories
WHERE
    category_id = $1
    AND deleted_at IS NOT NULL;

-- UpdateCategory: Updates category details
-- Purpose: Modify existing category data while maintaining soft delete integrity
-- Parameters:
--   $1: Category ID
--   $2: Updated name
--   $3: Updated description
--   $4: Updated slug
-- Returns:
--   Updated category record
-- Business Logic:
--   - Automatically updates the updated_at field
--   - Skips if category has been soft-deleted
-- name: UpdateCategory :one
UPDATE categories
SET
    name = $2,
    description = $3,
    slug_category = $4,
    image_category = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE
    category_id = $1
    AND deleted_at IS NULL
RETURNING
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at;

-- TrashCategory: Soft-deletes a category
-- Purpose: Moves category to trash without permanent deletion
-- Parameters:
--   $1: Category ID
-- Returns:
--   The soft-deleted category record
-- Business Logic:
--   - Updates deleted_at with current timestamp
--   - Prevents repeat trashing of already-deleted records
-- name: TrashCategory :one
UPDATE categories
SET
    deleted_at = current_timestamp
WHERE
    category_id = $1
    AND deleted_at IS NULL
RETURNING
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at,
    deleted_at;

-- RestoreCategory: Recovers a previously trashed category
-- Purpose: Restores a soft-deleted category for reuse
-- Parameters:
--   $1: Category ID
-- Returns:
--   Restored category record
-- Business Logic:
--   - Only applies to categories currently marked as deleted
-- name: RestoreCategory :one
UPDATE categories
SET
    deleted_at = NULL
WHERE
    category_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    category_id,
    name,
    description,
    slug_category,
    image_category,
    created_at,
    updated_at, deleted_at;

-- DeleteCategoryPermanently: Removes a soft-deleted category permanently
-- Purpose: Final cleanup of trashed categories
-- Parameters:
--   $1: Category ID
-- Returns:
--   Nothing (command only)
-- Business Logic:
--   - Ensures category is deleted only if it has been soft-deleted
-- name: DeleteCategoryPermanently :exec
DELETE FROM categories
WHERE
    category_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllCategories: Recovers all trashed categories
-- Purpose: Bulk restore of all soft-deleted category records
-- Parameters: None
-- Returns: None
-- Business Logic:
--   - Resets deleted_at for all soft-deleted records
-- name: RestoreAllCategories :exec
UPDATE categories
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentCategories: Permanently deletes all trashed categories
-- Purpose: Bulk purge of all soft-deleted category records
-- Parameters: None
-- Returns: None
-- Business Logic:
--   - Only affects records marked as deleted
-- name: DeleteAllPermanentCategories :exec
DELETE FROM categories WHERE deleted_at IS NOT NULL;