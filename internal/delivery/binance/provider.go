package binance

import (
	"context"
	"time"

	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/config/binance"
	"github.com/alextsa22/cryptocurrency-order-book/internal/delivery"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

type DeliveryMethod int

const (
	RestMethod = iota
	WebSocket
)

type DeliveryProvider struct {
	apiConfig *binance.ApiConfig
	method    DeliveryMethod
	rate      time.Duration
}

func NewDeliveryProvider(serviceConfig *config.ServiceConfig, apiConfig *binance.ApiConfig, method DeliveryMethod) delivery.DepthDeliveryProvider {
	return &DeliveryProvider{
		apiConfig: apiConfig,
		method:    method,
		rate:      time.Duration(serviceConfig.Rate) * time.Second,
	}
}

func (p *DeliveryProvider) DepthDelivery(ctx context.Context, symbol string, limit int, dataCh chan *domain.Depth) {
	switch p.method {
	case RestMethod:
		p.restDepthDelivery(ctx, symbol, limit, dataCh)
	case WebSocket:
		p.wsDepthDelivery(ctx, symbol, limit, dataCh)
	}
}
