package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

func addRandomTargetApp(t *testing.T, cid string) db.TargetApp {
	arg := db.AddTargetAppParams{
		Cid:   cid,
		AppID: util.RandomAppID(),
		Rule:  db.RuleType(util.RandomRule()),
	}

	target_app, err := testStore.AddTargetApp(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Cid, target_app.Cid)
	require.Equal(t, arg.AppID, target_app.AppID)
	require.Equal(t, arg.Rule, target_app.Rule)

	return target_app
}

func TestAddTargetApp(t *testing.T) {
	campaign := addRandomCampaign(t)
	addRandomTargetApp(t, campaign.Cid)
	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestGetTargetApp(t *testing.T) {
	campaign := addRandomCampaign(t)
	target_app := addRandomTargetApp(t, campaign.Cid)

	get_target_app, err := testStore.GetTargetApp(context.Background(), campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, target_app, get_target_app)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestDeleteTargetApp(t *testing.T) {
	campaign := addRandomCampaign(t)
	addRandomTargetApp(t, campaign.Cid)

	err := testStore.DeleteTargetApp(context.Background(), campaign.Cid)
	require.NoError(t, err)

	target_app, err := testStore.GetTargetApp(context.Background(), campaign.Cid)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, target_app)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}
