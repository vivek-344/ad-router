package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBSource         string `mapstructure:"DB_SOURCE"`
	RedisSource      string `mapstructure:"REDIS_SOURCE"`
	ServerAddress    string `mapstructure:"SERVER_ADDRESS"`
	PostgresPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresUser     string `mapstructure:"POSTGRES_USER"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
