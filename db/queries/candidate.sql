-- name: CreateCandidate :one
INSERT INTO candidates (
    first_name,last_name,email,phone
) VALUES (
    $1,$2,$3,$4
)
RETURNING *;

-- name: GetCandidateByID :one
SELECT * FROM candidates
WHERE id = $1;

-- name: GetCandidateByEmail :one
SELECT * FROM candidates
WHERE email = $1;

-- name: ListCandidates :many
SELECT * FROM candidates
ORDER BY first_name DESC
LIMIT $1 OFFSET $2;

-- name: UpdateCandidate :one
UPDATE candidates
SET first_name = $2, last_name = $3, phone = $4
WHERE id = $1
RETURNING *;

-- name: DeleteCandidate :exec
DELETE FROM candidates
WHERE id = $1;

-- name: CreateOffer :one
INSERT INTO offers (
    candidate_id, department_id, designation_id, level_id, 
    proposed_start_date, new_start_date, offer_letter_url, created_by
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8
)
RETURNING *;

-- name: GetOfferByID :one
SELECT * FROM offers
WHERE id = $1;

-- name: ListOffers :many
SELECT * FROM offers
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOfferStatus :one
UPDATE offers
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetOffersByCandidate :many
SELECT * FROM offers
WHERE candidate_id = $1
ORDER BY created_at DESC;

-- name: UpdateCandidateStatus :one
UPDATE candidates
SET status = $2
WHERE id = $1
RETURNING *;

-- name: GetOfferStatsByMonth :many
SELECT 
    DATE_TRUNC('month', created_at) AS month,
    status,
    COUNT(*) AS total
FROM offers
GROUP BY month, status
ORDER BY month ASC;
