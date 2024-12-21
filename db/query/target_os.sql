-- name: CreateTargetOs :one
INSERT INTO target_os (
    cid,
    os,
    rule
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetTargetOs :one
SELECT *
FROM target_os
WHERE cid = $1;

-- name: updateTargetOs :one
UPDATE target_os
SET os = $2, rule = $3
WHERE cid = $1
RETURNING *;

-- name: DeleteTargetOs :exec
DELETE FROM target_os
WHERE cid = $1;