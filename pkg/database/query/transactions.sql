-- GetTransactions: Retrieves paginated list of active transactions with search capability
-- Purpose: List all active transactions for management UI
-- Parameters:
--   $1: search_term - Optional text to filter transactions by payment method or status (NULL for no filter)
--   $2: limit - Maximum number of records to return (pagination limit)
--   $3: offset - Number of records to skip (pagination offset)
-- Returns:
--   All transaction fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted transactions (deleted_at IS NULL)
--   - Supports partial text matching on payment_method and payment_status fields (case-insensitive)
--   - Returns newest transactions first (created_at DESC)
--   - Provides total_count for client-side pagination
--   - Uses window function COUNT(*) OVER() for efficient total count
-- name: GetTransactions :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM transactions
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR payment_method ILIKE '%' || $1 || '%' OR payment_status ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- GetTransactionsActive: Retrieves paginated list of active transactions (identical to GetTransactions)
-- Purpose: Maintains consistent API pattern with other active/trashed endpoints
-- Parameters:
--   $1: search_term - Optional filter text for payment method/status
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Active transaction records with total_count
-- Business Logic:
--   - Same functionality as GetTransactions
--   - Exists for consistency in API design patterns
-- Note: Could be consolidated with GetTransactions if duplicate functionality is undesired
-- name: GetTransactionsActive :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM transactions
WHERE deleted_at IS NULL
AND ($1::TEXT IS NULL OR payment_method ILIKE '%' || $1 || '%' OR payment_status ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- GetTransactionsTrashed: Retrieves paginated list of soft-deleted transactions
-- Purpose: View and manage deleted transactions for audit/recovery
-- Parameters:
--   $1: search_term - Optional text to filter trashed transactions
--   $2: limit - Maximum records per page
--   $3: offset - Records to skip
-- Returns:
--   Trashed transaction records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Maintains same search functionality as active transaction queries
--   - Preserves chronological sorting (newest first)
--   - Used in transaction recovery/audit interfaces
--   - Includes total_count for pagination in trash management UI
-- name: GetTransactionsTrashed :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM transactions
WHERE deleted_at IS NOT NULL
AND ($1::TEXT IS NULL OR payment_method ILIKE '%' || $1 || '%' OR payment_status ILIKE '%' || $1 || '%')
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- GetTransactionByMerchant: Retrieves merchant-specific transactions with pagination
-- Purpose: List transactions filtered by merchant ID
-- Parameters:
--   $1: search_term - Optional text to filter transactions
--   $2: merchant_id - Optional merchant ID to filter by (NULL for all merchants)
--   $3: limit - Pagination limit
--   $4: offset - Pagination offset
-- Returns:
--   Transaction records with total_count
-- Business Logic:
--   - Combines merchant filtering with search functionality
--   - Maintains same sorting and pagination as other transaction queries
--   - Useful for merchant-specific transaction reporting
--   - NULL merchant_id parameter returns all merchants' transactions
-- name: GetTransactionByMerchant :many
SELECT
    *,
    COUNT(*) OVER() AS total_count
FROM transactions
WHERE deleted_at IS NULL
  AND ($1::TEXT IS NULL OR payment_method ILIKE '%' || $1 || '%' OR payment_status ILIKE '%' || $1 || '%')
  AND ($2::INT IS NULL OR merchant_id = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;


-- GetMonthlyAmountTransactionSuccess: Retrieves monthly success transaction metrics
-- Purpose: Generate monthly reports of successful transactions for analysis
-- Parameters:
--   $1: Start date of first comparison period (timestamp)
--   $2: End date of first comparison period (timestamp)
--   $3: Start date of second comparison period (timestamp)
--   $4: End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful (payment_status = 'success') transactions
--   - Excludes deleted transactions
--   - Compares two customizable time periods
--   - Includes gap-filling for months with no transactions
--   - Returns 0 values for months with no successful transactions
--   - Orders by most recent year/month first
-- name: GetMonthlyAmountTransactionSuccess :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        EXTRACT(MONTH FROM t.created_at)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'success'
        AND (
            (t.created_at >= $1::timestamp AND t.created_at <= $2::timestamp)
            OR (t.created_at >= $3::timestamp AND t.created_at <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at),
        EXTRACT(MONTH FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_success,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- GetYearlyAmountTransactionSuccess: Retrieves yearly success transaction metrics
-- Purpose: Generate annual reports of successful transactions
-- Parameters:
--   $1: Reference year for comparison (current year as integer)
-- Returns:
--   year: Year as text
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Only includes successful (payment_status = 'success') transactions
--   - Excludes deleted transactions
--   - Includes gap-filling for years with no transactions
--   - Returns 0 values for years with no successful transactions
--   - Orders by most recent year first
-- name: GetYearlyAmountTransactionSuccess :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'success'
        AND (
            EXTRACT(YEAR FROM t.created_at) = $1::integer
            OR EXTRACT(YEAR FROM t.created_at) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        total_success::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- GetMonthlyAmountTransactionFailed: Retrieves monthly failed transaction metrics
-- Purpose: Generate monthly reports of failed transactions for analysis
-- Parameters:
--   $1: Start date of first comparison period (timestamp)
--   $2: End date of first comparison period (timestamp)
--   $3: Start date of second comparison period (timestamp)
--   $4: End date of second comparison period (timestamp)
-- Returns:
--   year: Year as text
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   total_failed: Count of failed transactions
--   total_amount: Sum of failed transaction amounts
-- Business Logic:
--   - Only includes failed (payment_status = 'failed') transactions
--   - Excludes deleted transactions
--   - Compares two customizable time periods
--   - Includes gap-filling for months with no failed transactions
--   - Returns 0 values for months with no failed transactions
--   - Orders by most recent year/month first
-- name: GetMonthlyAmountTransactionFailed :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        EXTRACT(MONTH FROM t.created_at)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'failed'
        AND (
            (t.created_at >= $1::timestamp AND t.created_at <= $2::timestamp)
            OR (t.created_at >= $3::timestamp AND t.created_at <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at),
        EXTRACT(MONTH FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_failed,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- GetYearlyAmountTransactionFailed: Retrieves yearly failed transaction metrics
-- Purpose: Generate annual reports of failed transactions
-- Parameters:
--   $1: Reference year for comparison (current year as integer)
-- Returns:
--   year: Year as text
--   total_failed: Count of failed transactions
--   total_amount: Sum of failed transaction amounts
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Only includes failed (payment_status = 'failed') transactions
--   - Excludes deleted transactions
--   - Includes gap-filling for years with no transactions
--   - Returns 0 values for years with no failed transactions
--   - Orders by most recent year first
-- name: GetYearlyAmountTransactionFailed :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'failed'
        AND (
            EXTRACT(YEAR FROM t.created_at) = $1::integer
            OR EXTRACT(YEAR FROM t.created_at) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        total_failed::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- GetMonthlyTransactionMethodsSuccess: Analyzes successful payment method usage by month
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   payment_method: The payment method used
--   total_transactions: Count of successful transactions
--   total_amount: Total amount processed by this method
-- name: GetMonthlyTransactionMethodsSuccess :many
WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT 
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT 
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON 
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY 
    ac.activity_month, 
    ac.payment_method;

-- GetMonthlyTransactionMethodsFailed: Analyzes failed payment method usage by month
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   payment_method: The payment method used
--   total_transactions: Count of failed transactions
--   total_amount: Total amount that failed processing
-- name: GetMonthlyTransactionMethodsFailed :many
WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT 
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT 
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON 
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY 
    ac.activity_month, 
    ac.payment_method;

-- GetYearlyTransactionMethodsSuccess: Analyzes successful payment method usage by year
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   payment_method: The payment method used
--   total_transactions: Count of successful transactions
--   total_amount: Total amount processed by this method
-- name: GetYearlyTransactionMethodsSuccess :many
WITH
    year_range AS (
        SELECT 
            EXTRACT(YEAR FROM $1::timestamp)::int - 1 AS start_year,
            EXTRACT(YEAR FROM $1::timestamp)::int AS end_year
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_years AS (
        SELECT generate_series(
            (SELECT start_year FROM year_range),
            (SELECT end_year FROM year_range)
        )::int AS year
    ),
    all_combinations AS (
        SELECT 
            ay.year::text AS year,  
            pm.payment_method
        FROM all_years ay
        CROSS JOIN payment_methods pm
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::text AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (SELECT start_year FROM year_range) AND (SELECT end_year FROM year_range)
        GROUP BY
            EXTRACT(YEAR FROM t.created_at),
            t.payment_method
    )
SELECT 
    ac.year,  
    ac.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN yearly_transactions yt ON 
    ac.year = yt.year
    AND ac.payment_method = yt.payment_method
ORDER BY 
    ac.year,
    ac.payment_method;

-- GetYearlyTransactionMethodsFailed: Analyzes failed payment method usage by year
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   payment_method: The payment method used
--   total_transactions: Count of failed transactions
--   total_amount: Total amount that failed processing
-- name: GetYearlyTransactionMethodsFailed :many
WITH
    year_range AS (
        SELECT 
            EXTRACT(YEAR FROM $1::timestamp)::int - 1 AS start_year,
            EXTRACT(YEAR FROM $1::timestamp)::int AS end_year
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_years AS (
        SELECT generate_series(
            (SELECT start_year FROM year_range),
            (SELECT end_year FROM year_range)
        )::int AS year
    ),
    all_combinations AS (
        SELECT 
            ay.year::text AS year,  
            pm.payment_method
        FROM all_years ay
        CROSS JOIN payment_methods pm
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::text AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (SELECT start_year FROM year_range) AND (SELECT end_year FROM year_range)
        GROUP BY
            EXTRACT(YEAR FROM t.created_at),
            t.payment_method
    )
SELECT 
    ac.year, 
    ac.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN yearly_transactions yt ON 
    ac.year = yt.year
    AND ac.payment_method = yt.payment_method
ORDER BY 
    ac.year,  
    ac.payment_method;



-- GetMonthlyAmountTransactionSuccessByMerchant: Retrieves monthly success transaction metrics by merchant_id
-- Purpose: Generate monthly reports of successful transactions for analysis
-- Parameters:
--   $1: Start date of first comparison period (timestamp)
--   $2: End date of first comparison period (timestamp)
--   $3: Start date of second comparison period (timestamp)
--   $4: End date of second comparison period (timestamp)
--   $5: Merchant ID
-- Returns:
--   year: Year as text
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Only includes successful (payment_status = 'success') transactions
--   - Excludes deleted transactions
--   - Compares two customizable time periods
--   - Includes gap-filling for months with no transactions
--   - Returns 0 values for months with no successful transactions
--   - Orders by most recent year/month first
-- name: GetMonthlyAmountTransactionSuccessByMerchant :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        EXTRACT(MONTH FROM t.created_at)::integer AS month,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'success'
        AND t.merchant_id = $5
        AND (
            (t.created_at >= $1::timestamp AND t.created_at <= $2::timestamp)
            OR (t.created_at >= $3::timestamp AND t.created_at <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at),
        EXTRACT(MONTH FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_success,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_success,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- GetYearlyAmountTransactionSuccessByMerchant: Retrieves yearly success transaction metrics
-- Purpose: Generate annual reports of successful transactions by merchant_id
-- Parameters:
--   $1: Reference year for comparison (current year as integer)
--   $2: Merchant ID
-- Returns:
--   year: Year as text
--   total_success: Count of successful transactions
--   total_amount: Sum of successful transaction amounts
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Only includes successful (payment_status = 'success') transactions
--   - Excludes deleted transactions
--   - Includes gap-filling for years with no transactions
--   - Returns 0 values for years with no successful transactions
--   - Orders by most recent year first
-- name: GetYearlyAmountTransactionSuccessByMerchant :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        COUNT(*) AS total_success,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'success'
        AND t.merchant_id = $2
        AND (
            EXTRACT(YEAR FROM t.created_at) = $1::integer
            OR EXTRACT(YEAR FROM t.created_at) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        total_success::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_success,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;



-- GetMonthlyAmountTransactionFailedByMerchant: Retrieves monthly failed transaction metrics
-- Purpose: Generate monthly reports of failed transactions for analysis by merchant_id
-- Parameters:
--   $1: Start date of first comparison period (timestamp)
--   $2: End date of first comparison period (timestamp)
--   $3: Start date of second comparison period (timestamp)
--   $4: End date of second comparison period (timestamp)
--   $5: Merchant ID
-- Returns:
--   year: Year as text
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   total_failed: Count of failed transactions
--   total_amount: Sum of failed transaction amounts
-- Business Logic:
--   - Only includes failed (payment_status = 'failed') transactions
--   - Excludes deleted transactions
--   - Compares two customizable time periods
--   - Includes gap-filling for months with no failed transactions
--   - Returns 0 values for months with no failed transactions
--   - Orders by most recent year/month first
-- name: GetMonthlyAmountTransactionFailedByMerchant :many
WITH monthly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        EXTRACT(MONTH FROM t.created_at)::integer AS month,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'failed'
        AND t.merchant_id = $5
        AND (
            (t.created_at >= $1::timestamp AND t.created_at <= $2::timestamp)
            OR (t.created_at >= $3::timestamp AND t.created_at <= $4::timestamp)
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at),
        EXTRACT(MONTH FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        TO_CHAR(TO_DATE(month::text, 'MM'), 'Mon') AS month,
        total_failed,
        total_amount
    FROM
        monthly_data

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $1::timestamp)::text AS year,
        TO_CHAR($1::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $1::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $1::timestamp)::integer
    )

    UNION ALL

    SELECT
        EXTRACT(YEAR FROM $3::timestamp)::text AS year,
        TO_CHAR($3::timestamp, 'Mon') AS month,
        0 AS total_failed,
        0 AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM monthly_data
        WHERE year = EXTRACT(YEAR FROM $3::timestamp)::integer
        AND month = EXTRACT(MONTH FROM $3::timestamp)::integer
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC,
    TO_DATE(month, 'Mon') DESC;


-- GetYearlyAmountTransactionFailedByMerchant: Retrieves yearly failed transaction metrics
-- Purpose: Generate annual reports of failed transactions by merchant_id
-- Parameters:
--   $1: Reference year for comparison (current year as integer)
--   $2: Merchant ID
-- Returns:
--   year: Year as text
--   total_failed: Count of failed transactions
--   total_amount: Sum of failed transaction amounts
-- Business Logic:
--   - Compares current year with previous year automatically
--   - Only includes failed (payment_status = 'failed') transactions
--   - Excludes deleted transactions
--   - Includes gap-filling for years with no transactions
--   - Returns 0 values for years with no failed transactions
--   - Orders by most recent year first

-- name: GetYearlyAmountTransactionFailedByMerchant :many
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM t.created_at)::integer AS year,
        COUNT(*) AS total_failed,
        COALESCE(SUM(t.amount), 0)::integer AS total_amount
    FROM
        transactions t
    WHERE
        t.deleted_at IS NULL
        AND t.payment_status = 'failed'
        AND t.merchant_id = $2
        AND (
            EXTRACT(YEAR FROM t.created_at) = $1::integer
            OR EXTRACT(YEAR FROM t.created_at) = $1::integer - 1
        )
    GROUP BY
        EXTRACT(YEAR FROM t.created_at)
), formatted_data AS (
    SELECT
        year::text,
        total_failed::integer,
        total_amount::integer
    FROM
        yearly_data

    UNION ALL

    SELECT
        $1::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer
    )

    UNION ALL

    SELECT
        ($1::integer - 1)::text AS year,
        0::integer AS total_failed,
        0::integer AS total_amount
    WHERE NOT EXISTS (
        SELECT 1
        FROM yearly_data
        WHERE year = $1::integer - 1
    )
)
SELECT * FROM formatted_data
ORDER BY
    year DESC;


-- GetMonthlyTransactionMethodsByMerchantSuccess: Analyzes successful transactions by merchant and payment method monthly
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   merchant_id: The merchant identifier
--   merchant_name: The merchant's name
--   payment_method: The payment method used
--   total_transactions: Count of successful transactions
--   total_amount: Total amount processed by this method
-- name: GetMonthlyTransactionMethodsByMerchantSuccess :many
WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT 
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND t.merchant_id = $5  
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT 
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON 
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY 
    ac.activity_month, 
    ac.payment_method;

-- GetMonthlyTransactionMethodsByMerchantFailed: Analyzes failed transactions by merchant and payment method monthly
-- Parameters:
--   $1: Reference date (timestamp) - determines the 12-month analysis period
-- Returns:
--   month: 3-letter month abbreviation (e.g. 'Jan')
--   merchant_id: The merchant identifier
--   merchant_name: The merchant's name
--   payment_method: The payment method used
--   total_transactions: Count of failed transactions
--   total_amount: Total amount that failed processing
-- name: GetMonthlyTransactionMethodsByMerchantFailed :many
WITH
    date_ranges AS (
        SELECT
            $1::timestamp AS range1_start,
            $2::timestamp AS range1_end,
            $3::timestamp AS range2_start,
            $4::timestamp AS range2_end
    ),
    payment_methods AS (
        SELECT DISTINCT
            payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    ),
    all_months AS (
        SELECT generate_series(
            date_trunc('month', LEAST(
                (SELECT range1_start FROM date_ranges),
                (SELECT range2_start FROM date_ranges)
            )),
            date_trunc('month', GREATEST(
                (SELECT range1_end FROM date_ranges),
                (SELECT range2_end FROM date_ranges)
            )),
            interval '1 month'
        )::date AS activity_month
    ),
    all_combinations AS (
        SELECT 
            am.activity_month,
            pm.payment_method
        FROM all_months am
        CROSS JOIN payment_methods pm
    ),
    monthly_transactions AS (
        SELECT
            date_trunc('month', t.created_at)::date AS activity_month,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            COALESCE(SUM(t.amount), 0)::NUMERIC AS total_amount
        FROM transactions t
        JOIN date_ranges dr ON (
            t.created_at BETWEEN dr.range1_start AND dr.range1_end
            OR t.created_at BETWEEN dr.range2_start AND dr.range2_end
        )
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND t.merchant_id = $5  
        GROUP BY
            date_trunc('month', t.created_at),
            t.payment_method
    )
SELECT 
    TO_CHAR(ac.activity_month, 'Mon') AS month,
    ac.payment_method,
    COALESCE(mt.total_transactions, 0) AS total_transactions,
    COALESCE(mt.total_amount, 0) AS total_amount
FROM all_combinations ac
LEFT JOIN monthly_transactions mt ON 
    ac.activity_month = mt.activity_month
    AND ac.payment_method = mt.payment_method
ORDER BY 
    ac.activity_month, 
    ac.payment_method;

-- GetYearlyTransactionMethodsByMerchantSuccess: Analyzes successful transactions by merchant and payment method yearly
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   merchant_id: The merchant identifier
--   merchant_name: The merchant's name
--   payment_method: The payment method used
--   total_transactions: Count of successful transactions
--   total_amount: Total amount processed by this method
-- name: GetYearlyTransactionMethodsByMerchantSuccess :many
WITH
    year_series AS (
        SELECT generate_series(
            EXTRACT(YEAR FROM $1::timestamp)::integer - 2,
            EXTRACT(YEAR FROM $1::timestamp)::integer,
            1
        ) AS year
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::integer AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            SUM(t.amount)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'success'
            AND t.merchant_id = $2
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 1) AND EXTRACT(YEAR FROM $1::timestamp)
        GROUP BY
            year,
            t.payment_method
    ),
    payment_methods AS (
        SELECT DISTINCT payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    )
SELECT
    ys.year::text AS year,
    pm.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM year_series ys
CROSS JOIN payment_methods pm
LEFT JOIN yearly_transactions yt
    ON ys.year = yt.year
    AND pm.payment_method = yt.payment_method
ORDER BY ys.year, pm.payment_method;

-- GetYearlyTransactionMethodsByMerchantFailed: Analyzes failed transactions by merchant and payment method yearly
-- Parameters:
--   $1: Reference date (timestamp) - determines the 5-year analysis window
-- Returns:
--   year: 4-digit year as text
--   merchant_id: The merchant identifier
--   merchant_name: The merchant's name
--   payment_method: The payment method used
--   total_transactions: Count of failed transactions
--   total_amount: Total amount that failed processing
-- name: GetYearlyTransactionMethodsByMerchantFailed :many
WITH
    year_series AS (
        SELECT generate_series(
            EXTRACT(YEAR FROM $1::timestamp)::integer - 1,
            EXTRACT(YEAR FROM $1::timestamp)::integer,
            1
        ) AS year
    ),
    yearly_transactions AS (
        SELECT
            EXTRACT(YEAR FROM t.created_at)::integer AS year,
            t.payment_method,
            COUNT(t.transaction_id) AS total_transactions,
            SUM(t.amount)::NUMERIC AS total_amount
        FROM transactions t
        WHERE
            t.deleted_at IS NULL
            AND t.payment_status = 'failed'
            AND t.merchant_id = $2
            AND EXTRACT(YEAR FROM t.created_at) BETWEEN (EXTRACT(YEAR FROM $1::timestamp) - 1) AND EXTRACT(YEAR FROM $1::timestamp)
        GROUP BY
            year,
            t.payment_method
    ),
    payment_methods AS (
        SELECT DISTINCT payment_method
        FROM transactions
        WHERE deleted_at IS NULL
    )
SELECT
    ys.year::text AS year,
    pm.payment_method,
    COALESCE(yt.total_transactions, 0) AS total_transactions,
    COALESCE(yt.total_amount, 0) AS total_amount
FROM year_series ys
CROSS JOIN payment_methods pm
LEFT JOIN yearly_transactions yt
    ON ys.year = yt.year
    AND pm.payment_method = yt.payment_method
ORDER BY ys.year, pm.payment_method;


-- CreateTransactions: Creates a new transaction record
-- Purpose: Record a payment transaction in the system
-- Parameters:
--   $1: order_id - Reference to the associated order
--   $2: merchant_id - ID of the merchant receiving payment
--   $3: payment_method - Payment method used (e.g., 'cash', 'credit_card')
--   $4: amount - Total amount of the transaction
--   $5: change_amount - Amount of change given (for cash payments)
--   $6: payment_status - Status of payment ('pending', 'success', 'failed')
-- Returns: The complete created transaction record
-- Business Logic:
--   - Sets created_at timestamp automatically
--   - Validates all required payment fields
--   - Typically created during checkout process
-- name: CreateTransaction :one
INSERT INTO transactions (
    order_id, merchant_id, payment_method, amount, payment_status
) VALUES ($1, $2, $3, $4, $5)
RETURNING transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at;


-- GetTransactionByOrderID: Retrieves transaction by order reference
-- Purpose: Lookup transaction associated with specific order
-- Parameters:
--   $1: order_id - The order ID to search by
-- Returns: Transaction record if found and active
-- Business Logic:
--   - Only returns non-deleted transactions
--   - Used for order payment verification
--   - Helps prevent duplicate payments
-- name: GetTransactionByOrderID :one
SELECT 
    transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at
FROM transactions
WHERE order_id = $1
  AND deleted_at IS NULL;

-- GetTransactionByID: Retrieves transaction by transaction ID
-- Purpose: Fetch specific transaction details
-- Parameters:
--   $1: transaction_id - The unique transaction ID
-- Returns: Full transaction record if active
-- Business Logic:
--   - Excludes deleted transactions
--   - Used for transaction details/receipts
--   - Primary lookup for transaction management
-- name: GetTransactionByID :one
SELECT transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at
FROM transactions
WHERE transaction_id = $1
  AND deleted_at IS NULL;


-- UpdateTransaction: Modifies transaction details
-- Purpose: Update transaction information
-- Parameters:
--   $1: transaction_id - ID of transaction to update
--   $2: merchant_id - Updated merchant reference
--   $3: payment_method - Updated payment method
--   $4: amount - Updated transaction amount
--   $5: change_amount - Updated change amount
--   $6: payment_status - Updated payment status
--   $7: order_id - Updated order reference
-- Returns: Updated transaction record
-- Business Logic:
--   - Auto-updates updated_at timestamp
--   - Only modifies active transactions
--   - Validates all payment fields
--   - Used for payment corrections/updates
-- name: UpdateTransaction :one
UPDATE transactions
SET merchant_id = $2,
    payment_method = $3,
    amount = $4,
    payment_status = $5,
    order_id = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE transaction_id = $1
  AND deleted_at IS NULL
RETURNING transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at;


-- TrashTransaction: Soft-deletes a transaction
-- Purpose: Void/cancel a transaction without permanent deletion
-- Parameters:
--   $1: transaction_id - ID of transaction to cancel
-- Returns: The soft-deleted transaction record
-- Business Logic:
--   - Sets deleted_at to current timestamp
--   - Preserves transaction for reporting
--   - Only processes active transactions
--   - Can be restored if needed
-- name: TrashTransaction :one
UPDATE transactions
SET
    deleted_at = current_timestamp
WHERE
    transaction_id = $1
    AND deleted_at IS NULL
    RETURNING transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at,
    deleted_at;


-- RestoreTransaction: Recovers a soft-deleted transaction
-- Purpose: Reactivate a cancelled transaction
-- Parameters:
--   $1: transaction_id - ID of transaction to restore
-- Returns: The restored transaction record
-- Business Logic:
--   - Nullifies deleted_at field
--   - Only works on previously cancelled transactions
--   - Maintains all original transaction data
-- name: RestoreTransaction :one
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    transaction_id = $1
    AND deleted_at IS NOT NULL
  RETURNING transaction_id,
    order_id,
    merchant_id,
    payment_method,
    amount,
    payment_status,
    created_at,
    updated_at,
    deleted_at;


-- DeleteTransactionPermanently: Hard-deletes a transaction
-- Purpose: Completely remove transaction from database
-- Parameters:
--   $1: transaction_id - ID of transaction to delete
-- Business Logic:
--   - Permanent deletion of already cancelled transactions
--   - No return value (exec-only operation)
--   - Irreversible action - use with caution
--   - Should be restricted to admin users
-- name: DeleteTransactionPermanently :exec
DELETE FROM transactions WHERE transaction_id = $1 AND deleted_at IS NOT NULL;


-- RestoreAllTransactions: Mass restoration of cancelled transactions
-- Purpose: Recover all trashed transactions at once
-- Business Logic:
--   - Reactivates all soft-deleted transactions
--   - No parameters needed (bulk operation)
--   - Typically used during system recovery
-- name: RestoreAllTransactions :exec
UPDATE transactions
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- DeleteAllPermanentTransactions: Purges all cancelled transactions
-- Purpose: Clean up all soft-deleted transaction records
-- Business Logic:
--   - Irreversible bulk deletion operation
--   - Only affects already cancelled transactions
--   - Typically used during database maintenance
--   - Should be restricted to admin users
-- name: DeleteAllPermanentTransactions :exec
DELETE FROM transactions
WHERE
    deleted_at IS NOT NULL
    OR order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL);

-- DeleteTransactionByOrderPermanent: Permanently deletes all transactions for a specific order
-- Purpose: Cascading cleanup during order permanent deletion
-- Parameters:
--   $1: order_id
-- name: DeleteTransactionByOrderPermanent :exec
DELETE FROM transactions
WHERE
    order_id = $1;



