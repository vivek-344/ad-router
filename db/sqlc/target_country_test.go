package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

func createRandomTargetCountry(t *testing.T, cid string) db.TargetCountry {
	arg := db.CreateTargetCountryParams{
		Cid:     cid,
		Country: util.RandomCountry(),
		Rule:    db.RuleType(util.RandomRule()),
	}

	target_country, err := testQueries.CreateTargetCountry(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Cid, target_country.Cid)
	require.Equal(t, arg.Country, target_country.Country)
	require.Equal(t, arg.Rule, target_country.Rule)

	return target_country
}

func TestCreateTargetCountry(t *testing.T) {
	campaign := createRandomCampaign(t)
	createRandomTargetCountry(t, campaign.Cid)
	testQueries.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestGetTargetCountry(t *testing.T) {
	campaign := createRandomCampaign(t)
	target_country := createRandomTargetCountry(t, campaign.Cid)

	get_target_country, err := testQueries.GetTargetCountry(context.Background(), campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, target_country, get_target_country)

	testQueries.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestDeleteTargetCountry(t *testing.T) {
	campaign := createRandomCampaign(t)
	createRandomTargetCountry(t, campaign.Cid)

	err := testQueries.DeleteTargetCountry(context.Background(), campaign.Cid)
	require.NoError(t, err)

	target_country, err := testQueries.GetTargetCountry(context.Background(), campaign.Cid)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())
	require.Empty(t, target_country)

	testQueries.DeleteCampaign(context.Background(), campaign.Cid)
}
