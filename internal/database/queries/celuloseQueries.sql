-- name: GetLatest :many
SELECT
    *
FROM
    loads
ORDER BY
    created_at DESC
LIMIT 10;

-- name: GetSummary :many
SELECT
    material,
    SUM(average_weight) AS total_weight
FROM
    loads
WHERE
    created_at::date = $1
GROUP BY
    material;

-- name: GetFiltered :many
SELECT
    *
FROM
    loads
WHERE (material = $1
    OR NULLIF ($1, '') IS NULL)
AND ((NULLIF ($2, '') IS NOT NULL
        AND NULLIF ($3, '') IS NOT NULL
        AND created_at::date BETWEEN $2::date AND $3::date)
    OR (NULLIF ($2, '') IS NOT NULL
        AND NULLIF ($3, '') IS NULL
        AND created_at::date = $2::date)
    OR (NULLIF ($2, '') IS NULL
        AND NULLIF ($3, '') IS NULL))
ORDER BY
    created_at ASC;

-- name: CreateLoad :exec
INSERT INTO loads (id, material, average_weight, unit, created_at, timezone, operator, shift)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateLoad :exec
UPDATE
    loads
SET
    material = $1,
    created_at = $2,
    operator = $3,
    shift = $4
WHERE
    id = $5;

-- name: DeleteLoad :exec
DELETE FROM loads
WHERE id = $1;

