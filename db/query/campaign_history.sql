-- name: createCampaignHistory :exec
INSERT INTO campaign_history (
    cid,
    field_changed,
    old_value,
    new_value
) VALUES (
    $1, $2, $3, $4
);
