package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

func TestListAllCampaign(t *testing.T) {
	var all_campaign []db.Campaign
	for range 10 {
		campaign := createRandomCampaign(t)
		all_campaign = append(all_campaign, campaign)
	}

	listed_campaigns, err := testQueries.ListAllCampaign(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(listed_campaigns), 10)

	for _, campaign := range all_campaign {
		require.Contains(t, listed_campaigns, campaign)
		testQueries.DeleteCampaign(context.Background(), campaign.Cid)
	}
}

func TestListAllActiveCampaign(t *testing.T) {
	store := db.NewStore(testDB)
	var all_campaign []db.Campaign
	var active_campaigns []db.Campaign
	var active_count int
	for range 10 {
		campaign := createRandomCampaign(t)
		if util.RandomBool() {
			store.ToggleStatus(context.Background(), campaign.Cid)
			campaign, err := testQueries.GetCampaign(context.Background(), campaign.Cid)
			require.NoError(t, err)
			all_campaign = append(all_campaign, campaign)
		} else {
			all_campaign = append(all_campaign, campaign)
			active_campaigns = append(active_campaigns, campaign)
			active_count++
		}
	}

	listed_campaigns, err := testQueries.ListAllActiveCampaign(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(listed_campaigns), active_count)

	for _, campaign := range all_campaign {
		if campaign.Status == db.StatusType("active") {
			require.Contains(t, listed_campaigns, campaign)
		}
		testQueries.DeleteCampaign(context.Background(), campaign.Cid)
	}
}

func TestDeleteCampaign(t *testing.T) {
	campaign := createRandomCampaign(t)
	err := testQueries.DeleteCampaign(context.Background(), campaign.Cid)
	require.NoError(t, err)

	campaign, err = testQueries.GetCampaign(context.Background(), campaign.Cid)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, campaign)
}
