-- name: AddTargetCountry :one
INSERT INTO target_country (
    cid,
    country,
    rule
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTargetCountry :one
SELECT *
FROM target_country
WHERE cid = $1;

-- name: updateTargetCountry :one
UPDATE target_country
SET country = $2, rule = $3
WHERE cid = $1
RETURNING *;

-- name: DeleteTargetCountry :exec
DELETE FROM target_country
WHERE cid = $1;