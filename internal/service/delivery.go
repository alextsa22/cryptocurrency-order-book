package service

import (
	"context"
	"sync"
	
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DeliveryService implement DeliveryManager
type DeliveryService struct {}

// GetSymbols returns the set of available symbols.
func (s *DeliveryService) GetSymbols() []string {
	return nil
}

// GetDataChannel returns the channel into which data about the specified character is written.
func (s *DeliveryService) GetDataChannel(symbol string) (<-chan *domain.Depth, error) {
	return nil, nil
}

// RunProviders launches providers.
func (s *DeliveryService) RunProviders(wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	return nil, nil
}
