package binance

import (
	"github.com/spf13/viper"
)

type ApiConfig struct {
	PingUrl string
	REST
	WS
}

type REST struct {
	GetDepthUrl string
}

type WS struct {
	Scheme          string
	Host            string
	Path            string
	StreamDepthName string
}

func InitAPIConfig() (*ApiConfig, error) {
	viper.AddConfigPath("configs/binance")
	viper.SetConfigName("api")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg ApiConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
