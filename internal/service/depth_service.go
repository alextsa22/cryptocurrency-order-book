package service

import (
	"context"
	"sync"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DepthService defines an interface for working with fetchers who collect order books.
type DepthService interface {
	// GetSymbols returns a slice of symbols that fetchers collect.
	GetSymbols() []string
	// GetDataChannel returns the channel into which data about the specified character is written.
	GetDataChannel(symbol string) (<-chan *domain.Depth, error)
	// RunFetchers launches fetchers for each symbol.
	RunFetchers(wg *sync.WaitGroup) (context.Context, context.CancelFunc)
}
