package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Linear LinearConfig
}

type LinearConfig struct {
	APIKey string
	APIURL string
	TeamID string
}

func New() (*Config, error) {
	v := viper.New()

	v.SetDefault("linear.api_url", "https://api.linear.app/graphql")
	v.SetDefault("linear.team_id", "")
	v.SetDefault("linear.api_key", "")

	v.SetConfigName(".lcli")
	v.SetConfigType("yaml")

	// Where config file is resolved from
	v.AddConfigPath(".")
	v.AddConfigPath("$HOME")

	v.SetEnvPrefix("LCLI")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	cfg := &Config{
		Linear: LinearConfig{
			APIKey: v.GetString("linear.api_key"),
			APIURL: v.GetString("linear.api_url"),
			TeamID: v.GetString("linear.team_id"),
		},
	}

	return cfg, nil
}
