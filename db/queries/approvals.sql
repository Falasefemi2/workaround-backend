-- name: CreateApprovalSetup :one
INSERT INTO approval_setups (module_type, department_id, level_order, role_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetApprovalChain :many
SELECT * FROM approval_setups
WHERE module_type = $1
ORDER BY level_order ASC;

-- name: DeleteApprovalSetup :exec
DELETE FROM approval_setups
WHERE id = $1;

-- name: CreateApproval :one
INSERT INTO approvals (module_type, reference_id, approval_level, approver_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPendingApproval :one
SELECT * FROM approvals
WHERE reference_id = $1
AND status = 'pending'
ORDER BY approval_level ASC
LIMIT 1;

-- name: ActOnApproval :one
UPDATE approvals
SET status = $2, comment = $3, acted_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetApprovalsByReference :many
SELECT * FROM approvals
WHERE reference_id = $1
ORDER BY approval_level ASC;

-- name: GetPendingApprovalsByApprover :many
SELECT * FROM approvals
WHERE approver_id = $1
AND status = 'pending'
ORDER BY created_at ASC;

