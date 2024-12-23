package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListCampaignHistory(t *testing.T) {
	campaign := addRandomCampaign(t)

	for range 10 {
		testStore.ToggleStatus(context.Background(), campaign.Cid)
	}

	history, err := testStore.ListCampaignHistory(context.Background(), campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, len(history), 10)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}
