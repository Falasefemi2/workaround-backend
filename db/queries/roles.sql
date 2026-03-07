-- name: CreateRole :one
INSERT INTO roles (
    name, description
) VALUES (
 $1,$2
)
RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY name DESC 
LIMIT $1 OFFSET $2;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3
where id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1;

-- name: AssignRoleToUser :one
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
RETURNING *;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE user_id = $1 AND role_id = $2;

-- name: GetUserRoles :many
SELECT r.* FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1;
