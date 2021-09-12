package service

import (
	"context"
	"sync"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DepthService defines an interface for working with fetchers who collect order books.
type DepthService interface {
	GetSymbols() []string
	GetDataChannel(symbol string) (<-chan *domain.Depth, error)
	RunFetchers(wg *sync.WaitGroup) (context.Context, context.CancelFunc)
}