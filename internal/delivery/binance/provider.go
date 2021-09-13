package binance

import (
	"context"
	"time"

	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/delivery"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

type DeliveryProvider struct {
	rate time.Duration
}

func NewDeliveryProvider(serviceConfig *config.ServiceConfig) delivery.DepthDeliveryProvider {
	return &DeliveryProvider{
		rate: time.Duration(serviceConfig.Rate) * time.Second,
	}
}

func (p *DeliveryProvider) DepthDelivery(ctx context.Context, symbol string, limit int, dataCh chan *domain.Depth) {
	p.restDepthDelivery(ctx, symbol, limit, dataCh)
}