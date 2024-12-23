-- name: AddCampaign :one
INSERT INTO campaign (
  cid,
  name,
  img,
  cta
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCampaign :one
SELECT *
FROM campaign
WHERE cid = $1;

-- name: ListCampaigns :many
SELECT *
FROM campaign;

-- name: ListActiveCampaigns :many
SELECT *
FROM campaign
WHERE status = 'active'::status_type;

-- name: toggleStatus :one
UPDATE campaign
SET status = CASE 
    WHEN status = 'active'::status_type THEN 'inactive'::status_type
    ELSE 'active'::status_type
END
WHERE cid = $1
RETURNING status;

-- name: updateCampaignName :one
UPDATE campaign
SET name = $2
WHERE cid = $1
RETURNING *;

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