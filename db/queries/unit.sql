-- name: CreateUnit :one
INSERT INTO units (
    department_id, name, unit_lead_id
) VALUES (
   $1,$2,$3
)
RETURNING *;

-- name: GetUnitByID :one
SELECT * FROM units
WHERE id = $1;

-- name: ListUnits :many
SELECT * FROM units
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: DeleteUnit :exec
DELETE FROM units
WHERE id = $1;

-- name: UpdateUnit :one
UPDATE units 
SET name = $2
WHERE id = $1
RETURNING *;

-- name: AssignUnitLead :one
UPDATE units
SET unit_lead_id = $2
WHERE id = $1
RETURNING *;
