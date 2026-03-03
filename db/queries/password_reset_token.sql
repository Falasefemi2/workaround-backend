-- name: CreatePasswordResetToken :one
INSERT INTO password_reset_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPasswordResetToken :one
SELECT * FROM password_reset_tokens
WHERE token = $1;

-- name: MarkTokenUsed :exec
UPDATE password_reset_tokens
SET used = TRUE
WHERE token = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM password_reset_tokens
WHERE expires_at < NOW() OR used = TRUE;
