package delivery

import (
	"context"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DepthDeliveryProvider defines methods for providers.
type DepthDeliveryProvider interface {
	// DepthDelivery launches depth delivery.
	DepthDelivery(ctx context.Context, symbol string, limit int, dataCh chan *domain.Depth)
}
