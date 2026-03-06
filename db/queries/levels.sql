-- name: CreateLevel :one
INSERT INTO levels (
    name, code, annual_leave_days, annual_gross, support_total,
    minimum_leave_days, total_annual_leave_days, leave_expiration_interval,
    basic_salary, transport_allowance, domestic_allowance, utility_allowance, lunch_subsidy
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
)
RETURNING *;

-- name: GetLevelByID :one
SELECT * FROM levels
WHERE id = $1;

-- name: ListLevels :many
SELECT * FROM levels
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateLevel :one
UPDATE levels
SET name = $2, code = $3, annual_leave_days = $4, annual_gross = $5,
    support_total = $6, minimum_leave_days = $7, total_annual_leave_days = $8,
    leave_expiration_interval = $9, basic_salary = $10, transport_allowance = $11,
    domestic_allowance = $12, utility_allowance = $13, lunch_subsidy = $14
WHERE id = $1
RETURNING *;

-- name: DeleteLevel :exec
DELETE FROM levels
WHERE id = $1;

-- name: GetLevelByCode :one
SELECT * FROM levels
WHERE code = $1;

-- name: SearchLevels :many
SELECT * FROM levels
WHERE name ILIKE '%' || $1 || '%'
   OR code ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
