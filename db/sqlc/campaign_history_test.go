package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
)

func TestGetAllCampaignHistory(t *testing.T) {
	store := db.NewStore(testDB)
	campaign := createRandomCampaign(t)

	for range 10 {
		store.ToggleStatus(context.Background(), campaign.Cid)
	}

	history, err := testQueries.GetAllCampaignHistory(context.Background(), campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, len(history), 10)

	store.DeleteCampaign(context.Background(), campaign.Cid)
}
