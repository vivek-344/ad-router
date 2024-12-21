-- name: CreateCampaign :one
INSERT INTO campaign (
  cid,
  img,
  cta
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCampaign :one
SELECT *
FROM campaign
WHERE cid = $1;

-- name: ListAllCampaign :many
SELECT *
FROM campaign;

-- name: ListAllActiveCampaign :many
SELECT *
FROM campaign
WHERE status = 'active';

-- name: toggleStatus :one
UPDATE campaign
SET status = CASE 
    WHEN status = 'active' THEN 'inactive'
    ELSE 'inactive'
END
WHERE cid = $1
RETURNING status;

-- name: updateCampaignImage :one
UPDATE campaign
SET img = $2
WHERE cid = $1
RETURNING *;

-- name: updateCampaignCta :one
UPDATE campaign
SET cta = $2
WHERE cid = $1
RETURNING *;

-- name: DeleteCampaign :exec
DELETE FROM campaign
WHERE cid = $1;