package main_test

import (
	"context"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/vivek-344/AdRouter/util"
)

func TestRedis(t *testing.T) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	opt, err := redis.ParseURL(config.RedisSource)
	require.NoError(t, err)

	client := redis.NewClient(opt)
	defer client.Close()

	ctx := context.Background()

	err = client.Set(ctx, "foo", "bar", 0).Err()
	require.NoError(t, err)

	val, err := client.Get(ctx, "foo").Result()
	require.NoError(t, err)
	require.Equal(t, "bar", val)
}
