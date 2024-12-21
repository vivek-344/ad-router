-- name: CreateTargetApp :one
INSERT INTO target_app (
    cid,
    app_id,
    rule
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTargetApp :one
SELECT *
FROM target_app
WHERE cid = $1;

-- name: UpdateTargetApp :one
UPDATE target_app
SET app_id = $2, rule = $3
WHERE cid = $1
RETURNING *;

-- name: DeleteTargetApp :exec
DELETE FROM target_app
WHERE cid = $1;