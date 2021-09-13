package service

import (
	"context"
	"sync"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DeliveryManager defines an interface for working with providers who collect order books.
type DeliveryManager interface {
	// GetSymbols returns the set of available symbols.
	GetSymbols() []string
	// GetDataChannel returns the channel into which data about the specified character is written.
	GetDataChannel(symbol string) (<-chan *domain.Depth, error)
	// RunProviders launches providers.
	RunProviders(wg *sync.WaitGroup) (context.Context, context.CancelFunc)
}
