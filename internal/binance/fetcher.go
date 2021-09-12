package binance

import (
	"context"
	"fmt"
	"sync"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	"github.com/alextsa22/cryptocurrency-order-book/internal/service"
)

type DepthFetcher struct {
	symbols []string
	// dataChannels connects a symbol and a channel to receive data about this symbol.
	dataChannels map[string]chan *domain.Depth
}

func NewDepthFetcher(symbols []string) service.DepthService {
	// create channels for receiving data
	dataChannels := make(map[string]chan *domain.Depth)
	for _, symbol := range symbols {
		dataChannels[symbol] = make(chan *domain.Depth, 10)
	}

	return &DepthFetcher{
		symbols: symbols,
		dataChannels: dataChannels,
	}
}

func (f *DepthFetcher) GetSymbols() []string {
	return f.symbols
}

func (f *DepthFetcher) GetDataChannel(symbol string) (<-chan *domain.Depth, error) {
	c, exists := f.dataChannels[symbol]
	if !exists {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}
	return c, nil
}

func (f *DepthFetcher) RunFetchers(wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	return nil, nil
}
