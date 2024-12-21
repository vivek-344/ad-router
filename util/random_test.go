package util_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vivek-344/AdRouter/util"
)

func TestRandomBool(t *testing.T) {
	trueCount, falseCount := 0, 0

	for i := 0; i < 100; i++ {
		if util.RandomBool() {
			trueCount++
		} else {
			falseCount++
		}
	}

	require.Greater(t, trueCount, 0)
	require.Greater(t, falseCount, 0)
}

func TestRandomCid(t *testing.T) {
	name := util.RandomCid()
	require.GreaterOrEqual(t, len(name), 6)
	require.LessOrEqual(t, len(name), 8)
}

func TestRandomName(t *testing.T) {
	name := util.RandomName()
	require.GreaterOrEqual(t, len(name), 7)
	require.LessOrEqual(t, len(name), 10)
}

func TestRandomImg(t *testing.T) {
	require.Equal(t, "https://example.com", util.RandomImg())
}

func TestRandomCta(t *testing.T) {
	cta := []string{"Download", "Install", "Get", "Play"}
	require.Contains(t, cta, util.RandomCta())
}

func TestRandomAppID(t *testing.T) {
	appID := util.RandomAppID()
	appIDs := csvToSlice(appID)

	require.GreaterOrEqual(t, len(appIDs), 1)
	require.LessOrEqual(t, len(appIDs), 10)

	for _, appID := range appIDs {
		id := strings.Split(appID, ".")
		require.Equal(t, "com", id[0])
		require.GreaterOrEqual(t, len(id[1]), 5)
		require.LessOrEqual(t, len(id[1]), 10)
		require.GreaterOrEqual(t, len(id[2]), 5)
		require.LessOrEqual(t, len(id[2]), 10)
	}

	if len(appIDs) > 1 {
		joined := strings.Join(appIDs, ", ")
		require.Equal(t, appID, joined)
	}
}

func TestRandomCountry(t *testing.T) {
	allCountries := []string{"Russia", "Canada", "China", "United States", "Brazil", "Australia", "India", "Argentina", "Kazakhstan", "Algeria"}

	country := util.RandomCountry()
	countries := csvToSlice(country)

	require.GreaterOrEqual(t, len(countries), 1)
	require.LessOrEqual(t, len(countries), 10)

	for _, country := range countries {
		require.Contains(t, allCountries, country)
	}

	if len(countries) > 1 {
		joined := strings.Join(countries, ", ")
		require.Equal(t, country, joined)
	}
}

func TestRandomOs(t *testing.T) {
	allOs := []string{"Android", "iOS", "Web"}

	os := util.RandomOs()
	osList := csvToSlice(os)

	require.GreaterOrEqual(t, len(osList), 1)
	require.LessOrEqual(t, len(osList), 3)

	for _, os := range osList {
		require.Contains(t, allOs, os)
	}

	if len(osList) > 1 {
		joined := strings.Join(osList, ", ")
		require.Equal(t, os, joined)
	}
}

func TestRandomRule(t *testing.T) {
	rule := []string{"include", "exclude"}
	require.Contains(t, rule, util.RandomRule())
}

func csvToSlice(csv string) []string {
	return strings.Split(csv, ", ")
}
