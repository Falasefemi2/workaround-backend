-- name: CreateDesignation :one
INSERT INTO designations (
name, code
) VALUES (
   $1,$2
)
RETURNING *;

-- name: GetDesignationByID :one
SELECT * FROM designations
WHERE id = $1;

-- name: ListDesignations :many
SELECT * FROM designations
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: DeleteDesignation :exec  
DELETE FROM designations
WHERE id = $1;

-- name: UpdateDesignation :one
UPDATE designations
SET name = $2, code = $3
WHERE id = $1
RETURNING *;


