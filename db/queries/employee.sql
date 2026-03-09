-- name: CreateEmployee :one
INSERT INTO employees (
    user_id, employee_number, department_id, unit_id, designation_id,
    level_id, employment_type, employment_status, date_of_employment
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9
)
RETURNING *;

-- name: GetEmployeeByID :one
SELECT * FROM employees
WHERE id = $1;

-- name: GetEmployeeByUserID :one
SELECT * FROM employees
WHERE user_id = $1;

-- name: GetEmployeeByNumber :one
SELECT * FROM employees
WHERE employee_number = $1;

-- name: ListEmployees :many
SELECT * FROM employees
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateEmployee :one
UPDATE employees
SET
    department_id = $2,
    unit_id = $3,
    designation_id = $4,
    level_id = $5,
    employment_type = $6,
    employment_status = $7
WHERE id = $1
RETURNING *;

-- name: DeleteEmployee :exec
DELETE FROM employees
WHERE id = $1;
