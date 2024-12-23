package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/vivek-344/AdRouter/db/sqlc"
	"github.com/vivek-344/AdRouter/util"
)

func TestDelivery(t *testing.T) {
	campaigns := make([]db.Campaign, 3)

	campaigns[0] = addRandomCampaign(t)
	arg1 := db.AddTargetAppParams{
		Cid:   campaigns[0].Cid,
		AppID: "app1,app2,app3",
		Rule:  "include",
	}
	_, err := testStore.AddTargetApp(context.Background(), arg1)
	require.NoError(t, err)

	campaigns[1] = addRandomCampaign(t)
	arg2 := db.AddTargetCountryParams{
		Cid:     campaigns[1].Cid,
		Country: "US,UK,CA",
		Rule:    "exclude",
	}
	_, err = testStore.AddTargetCountry(context.Background(), arg2)
	require.NoError(t, err)

	campaigns[2] = addRandomCampaign(t)
	arg3 := db.AddTargetOsParams{
		Cid:  campaigns[2].Cid,
		Os:   "android,ios",
		Rule: "include",
	}
	_, err = testStore.AddTargetOs(context.Background(), arg3)
	require.NoError(t, err)

	testCases := []struct {
		name         string
		deliveryArgs db.DeliveryParams
		checkResults func(t *testing.T, results []db.DeliveryResult)
	}{
		{
			name: "Match app inclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app1",
				Country: "IN",
				Os:      "android",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.Contains(t, extractCids(results), campaigns[0].Cid)
			},
		},
		{
			name: "No match for app inclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app4",
				Country: "IN",
				Os:      "android",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.NotContains(t, extractCids(results), campaigns[0].Cid)
			},
		},
		{
			name: "Match country exclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app1",
				Country: "IN",
				Os:      "android",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.Contains(t, extractCids(results), campaigns[1].Cid)
			},
		},
		{
			name: "No match for country exclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app1",
				Country: "US",
				Os:      "android",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.NotContains(t, extractCids(results), campaigns[1].Cid)
			},
		},
		{
			name: "Match OS inclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app1",
				Country: "IN",
				Os:      "android",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.Contains(t, extractCids(results), campaigns[2].Cid)
			},
		},
		{
			name: "No match for OS inclusion rule",
			deliveryArgs: db.DeliveryParams{
				AppID:   "app1",
				Country: "IN",
				Os:      "windows",
			},
			checkResults: func(t *testing.T, results []db.DeliveryResult) {
				require.NotContains(t, extractCids(results), campaigns[2].Cid)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results, err := testStore.Delivery(context.Background(), tc.deliveryArgs)
			require.NoError(t, err)
			tc.checkResults(t, results)
		})
	}

	for _, campaign := range campaigns {
		err := testStore.DeleteCampaign(context.Background(), campaign.Cid)
		require.NoError(t, err)
	}
}

func extractCids(results []db.DeliveryResult) []string {
	cids := make([]string, len(results))
	for i, result := range results {
		cids[i] = result.Cid
	}
	return cids
}

func addRandomCampaign(t *testing.T) db.Campaign {
	arg := db.AddCampaignParams{
		Cid:  util.RandomCid(),
		Name: util.RandomName(),
		Img:  util.RandomImg(),
		Cta:  util.RandomCta(),
	}

	campaign, err := testStore.AddCampaign(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Cid, campaign.Cid)
	require.Equal(t, arg.Name, campaign.Name)
	require.Equal(t, arg.Img, campaign.Img)
	require.Equal(t, arg.Cta, campaign.Cta)
	require.Equal(t, db.StatusType("active"), campaign.Status)
	require.NotEmpty(t, campaign.CreatedAt)

	return campaign
}

func createRandomCampaign(t *testing.T) db.CreateCampaignResult {
	arg := db.CreateCampaignParams{
		Cid:  util.RandomCid(),
		Name: util.RandomName(),
		Img:  util.RandomImg(),
		Cta:  util.RandomCta(),
	}
	if util.RandomBool() {
		arg.AppID = util.RandomAppID()
		arg.AppRule = db.RuleType(util.RandomRule())
	}
	if util.RandomBool() {
		arg.Country = util.RandomCountry()
		arg.CountryRule = db.RuleType(util.RandomRule())
	}
	if util.RandomBool() {
		arg.Os = util.RandomOs()
		arg.OsRule = db.RuleType(util.RandomRule())
	}

	campaign, err := testStore.CreateCampaign(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Cid, campaign.Cid)
	require.Equal(t, arg.Name, campaign.Name)
	require.Equal(t, arg.Img, campaign.Img)
	require.Equal(t, arg.Cta, campaign.Cta)
	require.Equal(t, arg.AppID, campaign.AppID)
	require.Equal(t, arg.AppRule, campaign.AppRule)
	require.Equal(t, arg.Country, campaign.Country)
	require.Equal(t, arg.CountryRule, campaign.CountryRule)
	require.Equal(t, arg.Os, campaign.Os)
	require.Equal(t, arg.OsRule, campaign.OsRule)
	require.Equal(t, db.StatusType("active"), campaign.Status)
	require.NotEmpty(t, campaign.CreatedAt)

	return campaign
}

func TestCreateCampaign(t *testing.T) {
	campaign := createRandomCampaign(t)
	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestReadCampaign(t *testing.T) {
	campaign := createRandomCampaign(t)
	read_campaign, err := testStore.ReadCampaign(context.Background(), campaign.Cid)
	require.NoError(t, err)

	require.Equal(t, campaign.Cid, read_campaign.Cid)
	require.Equal(t, campaign.Name, read_campaign.Name)
	require.Equal(t, campaign.Img, read_campaign.Img)
	require.Equal(t, campaign.Cta, read_campaign.Cta)
	require.Equal(t, campaign.AppID, read_campaign.AppID)
	require.Equal(t, campaign.AppRule, read_campaign.AppRule)
	require.Equal(t, campaign.Country, read_campaign.Country)
	require.Equal(t, campaign.CountryRule, read_campaign.CountryRule)
	require.Equal(t, campaign.Os, read_campaign.Os)
	require.Equal(t, campaign.OsRule, read_campaign.OsRule)
	require.Equal(t, campaign.Status, read_campaign.Status)
	require.Equal(t, campaign.CreatedAt, read_campaign.CreatedAt)

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestToggleStatus(t *testing.T) {
	old_campaign := addRandomCampaign(t)

	err := testStore.ToggleStatus(context.Background(), old_campaign.Cid)
	require.NoError(t, err)

	updated_campaign, err := testStore.GetCampaign(context.Background(), old_campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, old_campaign.Cid, updated_campaign.Cid)
	require.Equal(t, old_campaign.Name, updated_campaign.Name)
	require.Equal(t, old_campaign.Img, updated_campaign.Img)
	require.Equal(t, old_campaign.Cta, updated_campaign.Cta)
	require.Equal(t, db.StatusType("inactive"), updated_campaign.Status)
	require.Equal(t, old_campaign.CreatedAt, updated_campaign.CreatedAt)

	campaignHistory, err := testStore.GetCampaignHistory(context.Background(), old_campaign.Cid)
	require.NoError(t, err)
	require.NotEmpty(t, campaignHistory.ID)
	require.Equal(t, updated_campaign.Cid, campaignHistory.Cid)
	require.Equal(t, "active", campaignHistory.OldValue)
	require.Equal(t, string(updated_campaign.Status), campaignHistory.NewValue)
	require.Equal(t, "status", campaignHistory.FieldChanged)
	require.NotEmpty(t, campaignHistory.UpdatedAt)

	old_campaign = updated_campaign
	err = testStore.ToggleStatus(context.Background(), old_campaign.Cid)
	require.NoError(t, err)

	updated_campaign, err = testStore.GetCampaign(context.Background(), old_campaign.Cid)
	require.NoError(t, err)
	require.Equal(t, old_campaign.Cid, updated_campaign.Cid)
	require.Equal(t, old_campaign.Name, updated_campaign.Name)
	require.Equal(t, old_campaign.Img, updated_campaign.Img)
	require.Equal(t, old_campaign.Cta, updated_campaign.Cta)
	require.Equal(t, db.StatusType("active"), updated_campaign.Status)
	require.Equal(t, old_campaign.CreatedAt, updated_campaign.CreatedAt)

	campaignHistory, err = testStore.GetCampaignHistory(context.Background(), updated_campaign.Cid)
	require.NoError(t, err)
	require.NotEmpty(t, campaignHistory.ID)
	require.Equal(t, updated_campaign.Cid, campaignHistory.Cid)
	require.Equal(t, "inactive", campaignHistory.OldValue)
	require.Equal(t, string(updated_campaign.Status), campaignHistory.NewValue)
	require.Equal(t, "status", campaignHistory.FieldChanged)
	require.NotEmpty(t, campaignHistory.UpdatedAt)

	testStore.DeleteCampaign(context.Background(), updated_campaign.Cid)
}

func TestUpdateCampaignName(t *testing.T) {
	old_campaign := addRandomCampaign(t)

	var newName string
	for {
		newName = util.RandomName()
		if newName != old_campaign.Name {
			break
		}
	}

	arg := db.UpdateCampaignNameParams{
		Cid:  old_campaign.Cid,
		Name: newName,
	}

	updated_campaign, err := testStore.UpdateCampaignName(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, old_campaign.Cid, updated_campaign.Cid)
	require.Equal(t, arg.Name, updated_campaign.Name)
	require.Equal(t, old_campaign.Img, updated_campaign.Img)
	require.Equal(t, old_campaign.Cta, updated_campaign.Cta)
	require.Equal(t, old_campaign.Status, updated_campaign.Status)
	require.Equal(t, old_campaign.CreatedAt, updated_campaign.CreatedAt)

	campaignHistory, err := testStore.GetCampaignHistory(context.Background(), arg.Cid)
	require.NoError(t, err)
	require.NotEmpty(t, campaignHistory.ID)
	require.Equal(t, arg.Cid, campaignHistory.Cid)
	require.Equal(t, old_campaign.Name, campaignHistory.OldValue)
	require.Equal(t, updated_campaign.Name, campaignHistory.NewValue)
	require.Equal(t, "name", campaignHistory.FieldChanged)
	require.NotEmpty(t, campaignHistory.UpdatedAt)

	testStore.DeleteCampaign(context.Background(), arg.Cid)
}

func TestUpdateCampaignCta(t *testing.T) {
	old_campaign := addRandomCampaign(t)

	var newCta string
	for {
		newCta = util.RandomCta()
		if newCta != old_campaign.Cta {
			break
		}
	}

	arg := db.UpdateCampaignCtaParams{
		Cid: old_campaign.Cid,
		Cta: newCta,
	}

	updated_campaign, err := testStore.UpdateCampaignCta(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, old_campaign.Cid, updated_campaign.Cid)
	require.Equal(t, old_campaign.Name, updated_campaign.Name)
	require.Equal(t, old_campaign.Img, updated_campaign.Img)
	require.Equal(t, arg.Cta, updated_campaign.Cta)
	require.Equal(t, old_campaign.Status, updated_campaign.Status)
	require.Equal(t, old_campaign.CreatedAt, updated_campaign.CreatedAt)

	campaignHistory, err := testStore.GetCampaignHistory(context.Background(), arg.Cid)
	require.NoError(t, err)
	require.NotEmpty(t, campaignHistory.ID)
	require.Equal(t, arg.Cid, campaignHistory.Cid)
	require.Equal(t, old_campaign.Cta, campaignHistory.OldValue)
	require.Equal(t, updated_campaign.Cta, campaignHistory.NewValue)
	require.Equal(t, "cta", campaignHistory.FieldChanged)
	require.NotEmpty(t, campaignHistory.UpdatedAt)

	testStore.DeleteCampaign(context.Background(), arg.Cid)
}

func TestUpdateCampaignImage(t *testing.T) {
	old_campaign := addRandomCampaign(t)

	var newImg string
	for {
		newImg = util.RandomImg()
		if newImg != old_campaign.Img {
			break
		}
	}

	arg := db.UpdateCampaignImageParams{
		Cid: old_campaign.Cid,
		Img: newImg,
	}

	updated_campaign, err := testStore.UpdateCampaignImage(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, old_campaign.Cid, updated_campaign.Cid)
	require.Equal(t, old_campaign.Name, updated_campaign.Name)
	require.Equal(t, arg.Img, updated_campaign.Img)
	require.Equal(t, old_campaign.Cta, updated_campaign.Cta)
	require.Equal(t, old_campaign.Status, updated_campaign.Status)
	require.Equal(t, old_campaign.CreatedAt, updated_campaign.CreatedAt)

	campaignHistory, err := testStore.GetCampaignHistory(context.Background(), arg.Cid)
	require.NoError(t, err)
	require.NotEmpty(t, campaignHistory.ID)
	require.Equal(t, arg.Cid, campaignHistory.Cid)
	require.Equal(t, old_campaign.Img, campaignHistory.OldValue)
	require.Equal(t, updated_campaign.Img, campaignHistory.NewValue)
	require.Equal(t, "img", campaignHistory.FieldChanged)
	require.NotEmpty(t, campaignHistory.UpdatedAt)

	testStore.DeleteCampaign(context.Background(), arg.Cid)
}

func TestUpdateTargetApp(t *testing.T) {
	campaign := addRandomCampaign(t)

	new_arg := db.AddTargetAppParams{
		Cid:   campaign.Cid,
		AppID: util.RandomAppID(),
		Rule:  db.RuleType(util.RandomRule()),
	}

	old_target_app, err := testStore.AddTargetApp(context.Background(), new_arg)
	require.NoError(t, err)
	require.Equal(t, new_arg.Cid, old_target_app.Cid)
	require.Equal(t, new_arg.AppID, old_target_app.AppID)
	require.Equal(t, new_arg.Rule, old_target_app.Rule)

	var newAppID string
	for {
		newAppID = util.RandomAppID()
		if newAppID != old_target_app.AppID {
			break
		}
	}

	update_arg := db.UpdateTargetAppParams{
		Cid:   campaign.Cid,
		AppID: newAppID,
		Rule:  db.RuleType(util.RandomRule()),
	}

	updated_target_app, err := testStore.UpdateTargetApp(context.Background(), update_arg)
	require.NoError(t, err)
	require.Equal(t, old_target_app.Cid, updated_target_app.Cid)
	require.Equal(t, update_arg.AppID, updated_target_app.AppID)
	require.Equal(t, update_arg.Rule, updated_target_app.Rule)

	campaignHistory, err := testStore.GetLastTwoCampaignHistory(context.Background(), campaign.Cid)
	require.NoError(t, err)

	for _, history := range campaignHistory {
		expected_changes := []string{"app_id", "app_rule"}
		require.NotEmpty(t, history.ID)
		require.Equal(t, campaign.Cid, history.Cid)
		require.Contains(t, expected_changes, history.FieldChanged)
		require.NotEmpty(t, history.UpdatedAt)
		if history.FieldChanged == "app_id" {
			require.Equal(t, "app_id", history.FieldChanged)
			require.Equal(t, old_target_app.AppID, history.OldValue)
			require.Equal(t, updated_target_app.AppID, history.NewValue)
		} else {
			require.Equal(t, "app_rule", history.FieldChanged)
			require.Equal(t, string(old_target_app.Rule), history.OldValue)
			require.Equal(t, string(updated_target_app.Rule), history.NewValue)
		}
	}

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestUpdateTargetOs(t *testing.T) {
	campaign := addRandomCampaign(t)

	new_arg := db.AddTargetOsParams{
		Cid:  campaign.Cid,
		Os:   util.RandomOs(),
		Rule: db.RuleType(util.RandomRule()),
	}

	old_target_os, err := testStore.AddTargetOs(context.Background(), new_arg)
	require.NoError(t, err)
	require.Equal(t, new_arg.Cid, old_target_os.Cid)
	require.Equal(t, new_arg.Os, old_target_os.Os)
	require.Equal(t, new_arg.Rule, old_target_os.Rule)

	var newOs string
	for {
		newOs = util.RandomOs()
		if newOs != old_target_os.Os {
			break
		}
	}

	update_arg := db.UpdateTargetOsParams{
		Cid:  campaign.Cid,
		Os:   newOs,
		Rule: db.RuleType(util.RandomRule()),
	}

	updated_target_os, err := testStore.UpdateTargetOs(context.Background(), update_arg)
	require.NoError(t, err)
	require.Equal(t, campaign.Cid, updated_target_os.Cid)
	require.Equal(t, update_arg.Os, updated_target_os.Os)
	require.Equal(t, update_arg.Rule, updated_target_os.Rule)

	campaignHistory, err := testStore.GetLastTwoCampaignHistory(context.Background(), campaign.Cid)
	require.NoError(t, err)

	for _, history := range campaignHistory {
		expected_changes := []string{"os", "os_rule"}
		require.NotEmpty(t, history.ID)
		require.Equal(t, campaign.Cid, history.Cid)
		require.Contains(t, expected_changes, history.FieldChanged)
		require.NotEmpty(t, history.UpdatedAt)
		if history.FieldChanged == "os" {
			require.Equal(t, "os", history.FieldChanged)
			require.Equal(t, old_target_os.Os, history.OldValue)
			require.Equal(t, updated_target_os.Os, history.NewValue)
		} else {
			require.Equal(t, "os_rule", history.FieldChanged)
			require.Equal(t, string(old_target_os.Rule), history.OldValue)
			require.Equal(t, string(updated_target_os.Rule), history.NewValue)
		}
	}

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}

func TestUpdateTargetCountry(t *testing.T) {
	campaign := addRandomCampaign(t)

	new_arg := db.AddTargetCountryParams{
		Cid:     campaign.Cid,
		Country: util.RandomCountry(),
		Rule:    db.RuleType(util.RandomRule()),
	}

	old_target_country, err := testStore.AddTargetCountry(context.Background(), new_arg)
	require.NoError(t, err)
	require.Equal(t, new_arg.Cid, old_target_country.Cid)
	require.Equal(t, new_arg.Country, old_target_country.Country)
	require.Equal(t, new_arg.Rule, old_target_country.Rule)

	var newCountry string
	for {
		newCountry = util.RandomCountry()
		if newCountry != old_target_country.Country {
			break
		}
	}

	update_arg := db.UpdateTargetCountryParams{
		Cid:     campaign.Cid,
		Country: newCountry,
		Rule:    db.RuleType(util.RandomRule()),
	}

	updated_target_country, err := testStore.UpdateTargetCountry(context.Background(), update_arg)
	require.NoError(t, err)
	require.Equal(t, campaign.Cid, updated_target_country.Cid)
	require.Equal(t, update_arg.Country, updated_target_country.Country)
	require.Equal(t, update_arg.Rule, updated_target_country.Rule)

	campaignHistory, err := testStore.GetLastTwoCampaignHistory(context.Background(), campaign.Cid)
	require.NoError(t, err)

	for _, history := range campaignHistory {
		expected_changes := []string{"country", "country_rule"}
		require.NotEmpty(t, history.ID)
		require.Equal(t, campaign.Cid, history.Cid)
		require.Contains(t, expected_changes, history.FieldChanged)
		require.NotEmpty(t, history.UpdatedAt)
		if history.FieldChanged == "country" {
			require.Equal(t, "country", history.FieldChanged)
			require.Equal(t, old_target_country.Country, history.OldValue)
			require.Equal(t, updated_target_country.Country, history.NewValue)
		} else {
			require.Equal(t, "country_rule", history.FieldChanged)
			require.Equal(t, string(old_target_country.Rule), history.OldValue)
			require.Equal(t, string(updated_target_country.Rule), history.NewValue)
		}
	}

	testStore.DeleteCampaign(context.Background(), campaign.Cid)
}
