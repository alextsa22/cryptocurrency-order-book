package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Symbols []string
	Limit   int
	Rate    int
}

func InitServiceConfig() (*ServiceConfig, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg ServiceConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if len(cfg.Symbols) < 1 {
		return nil, fmt.Errorf("unable to run 0 providers, please add symbols into config.yml")
	}

	return &cfg, nil
}
