package util_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vivek-344/AdRouter/util"
)

func TestLoadConfig(t *testing.T) {
	config, err := util.LoadConfig("..")

	require.NoError(t, err)
	require.NotEmpty(t, config.DBSource)
	require.NotEmpty(t, config.RedisSource)
	require.NotEmpty(t, config.ServerAddress)
	require.NotEmpty(t, config.PostgresPassword)
	require.NotEmpty(t, config.PostgresUser)
}
