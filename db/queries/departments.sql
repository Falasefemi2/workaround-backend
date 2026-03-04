-- name: CreateDepartment :one
INSERT INTO departments (
name, code, hod_id
) VALUES (
	$1,$2,$3
)
RETURNING *;

-- name: GetDepartmentByID :one 
SELECT * FROM departments
WHERE id = $1; 

-- name: ListDepartments :many
SELECT * FROM departments
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteDepartment :exec  
DELETE FROM departments
WHERE id = $1;

-- name: UpdateDepartment :one
UPDATE departments
SET name = $2, code = $3, hod_id= $4
WHERE id = $1
RETURNING *;

-- name: AssignHodToDepartment :one
UPDATE departments
SET hod_id = $2
WHERE id = $1
RETURNING *;

