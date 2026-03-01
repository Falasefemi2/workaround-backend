-- name: CreateUser :one
INSERT INTO users (
    email, password_hash, first_name, last_name, phone, status, user_type
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, phone = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
