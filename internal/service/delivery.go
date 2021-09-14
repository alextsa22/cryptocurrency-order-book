package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/delivery"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
)

// DeliveryService implement DeliveryManager
type DeliveryService struct {
	symbols []string
	limit   int
	rate    time.Duration
	// dataChannels connects a symbol and a channel to receive data about this symbol.
	dataChannels     map[string]chan *domain.Depth
	deliveryProvider delivery.DepthDeliveryProvider
}

func NewDeliveryService(serviceConfig *config.ServiceConfig, deliveryProvider delivery.DepthDeliveryProvider) DeliveryManager {
	symbols := serviceConfig.Symbols
	normalizeSymbolsList(symbols)

	// create channels for receiving data
	dataChannels := make(map[string]chan *domain.Depth)
	for _, symbol := range symbols {
		dataChannels[symbol] = make(chan *domain.Depth, 10)
	}

	service := DeliveryService{
		symbols:          symbols,
		limit:            serviceConfig.Limit,
		rate:             time.Duration(serviceConfig.Rate) * time.Second,
		dataChannels:     dataChannels,
		deliveryProvider: deliveryProvider,
	}
	return &service
}

// GetSymbols returns the set of available symbols.
func (s *DeliveryService) GetSymbols() []string {
	return s.symbols
}

// GetDataChannel returns the channel into which data about the specified character is written.
func (s *DeliveryService) GetDataChannel(symbol string) (<-chan *domain.Depth, error) {
	symbol = normalizeSymbol(symbol)
	ch, exists := s.dataChannels[symbol]
	if !exists {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}
	return ch, nil
}

// RunProviders launches providers.
func (s *DeliveryService) RunProviders(wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	for symbol, ch := range s.dataChannels {
		wg.Add(1)
		go func(symbol string, dataCh chan *domain.Depth) {
			defer wg.Done()
			s.deliveryProvider.DepthDelivery(ctx, symbol, s.limit, dataCh)
		}(symbol, ch)
	}
	time.Sleep(time.Second) // give a little time to initialize the goroutines
	return ctx, cancel
}

// normalizeSymbolsList normalizes symbols in a slice using normalizeSymbol.
func normalizeSymbolsList(symbols []string) {
	for i, symbol := range symbols {
		symbols[i] = normalizeSymbol(symbol)
	}
}

// normalizeSymbol converts the string with the symbol to the desired form.
func normalizeSymbol(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	symbol = strings.ToUpper(symbol)
	return symbol
}
