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

-- name: ToggleStatus :exec
UPDATE campaign
SET status = CASE 
    WHEN status = 'active' THEN 'inactive'
    ELSE 'inactive'
END
WHERE cid = $1;

-- name: UpdateCampaignImage :one
UPDATE campaign
SET img = $2
WHERE cid = $1
RETURNING *;

-- name: UpdateCampaignCta :one
UPDATE campaign
SET cta = $2
WHERE cid = $1
RETURNING *;

-- name: DeleteCampaign :exec
DELETE FROM campaign
WHERE cid = $1;