package config

import "github.com/spf13/viper"

type ServiceConfig struct {
	Symbols     []string
	Limit       int
	FetcherRate int `mapstructure:"rate"`
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
	return &cfg, nil
}
