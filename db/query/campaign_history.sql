-- name: createCampaignHistory :exec
INSERT INTO campaign_history (
    cid,
    field_changed,
    old_value,
    new_value
) VALUES (
    $1, $2, $3, $4
);

-- name: GetCampaignHistory :one
SELECT *
FROM campaign_history
WHERE cid = $1
ORDER BY updated_at DESC
LIMIT 1;

-- name: GetLastTwoCampaignHistory :many
SELECT *
FROM campaign_history
WHERE cid = $1
ORDER BY updated_at DESC
LIMIT 2;

-- name: GetAllCampaignHistory :many
SELECT *
FROM campaign_history
WHERE cid = $1
ORDER BY updated_at DESC;