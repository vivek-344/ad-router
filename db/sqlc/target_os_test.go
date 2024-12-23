package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

func addRandomTargetOs(t *testing.T, cid string) db.TargetOs {
	arg := db.AddTargetOsParams{
		Cid:  cid,
		Os:   util.RandomOs(),
		Rule: db.RuleType(util.RandomRule()),
	}

	target_os, err := testStore.AddTargetOs(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Cid, target_os.Cid)
	require.Equal(t, arg.Os, target_os.Os)
	require.Equal(t, arg.Rule, target_os.Rule)

	return target_os
}

func TestAddTargetOs(t *testing.T) {
	campaign := addRandomCampaign(t)
	addRandomTargetOs(t, campaign.Cid)
	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestGetTargetOs(t *testing.T) {
	campaign := addRandomCampaign(t)
	target_os := addRandomTargetOs(t, campaign.Cid)

	get_target_os, err := testStore.GetTargetOs(context.Background(), campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, target_os, get_target_os)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestDeleteTargetOs(t *testing.T) {
	campaign := addRandomCampaign(t)
	addRandomTargetOs(t, campaign.Cid)

	err := testStore.DeleteTargetOs(context.Background(), campaign.Cid)
	require.NoError(t, err)

	target_os, err := testStore.GetTargetOs(context.Background(), campaign.Cid)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, target_os)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}
