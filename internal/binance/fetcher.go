package binance

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	"github.com/alextsa22/cryptocurrency-order-book/internal/service"
)

// DepthFetcher implements a service.DepthService.
type DepthFetcher struct {
	symbols     []string
	limit       int
	fetcherRate time.Duration
	// dataChannels connects a symbol and a channel to receive data about this symbol.
	dataChannels map[string]chan *domain.Depth
}

func NewDepthFetcher(serviceCfg *config.ServiceConfig) (service.DepthService, error) {
	symbols := serviceCfg.Symbols
	normalizeSymbolsList(symbols)

	// create channels for receiving data
	dataChannels := make(map[string]chan *domain.Depth)
	for _, symbol := range symbols {
		dataChannels[symbol] = make(chan *domain.Depth, 10)
	}

	depth := DepthFetcher{
		symbols:      symbols,
		limit:        serviceCfg.Limit,
		fetcherRate:  time.Duration(serviceCfg.FetcherRate) * time.Second,
		dataChannels: dataChannels,
	}
	if err := depth.ping(); err != nil {
		return nil, fmt.Errorf("ping: %v", err)
	}
	return &depth, nil
}

func (f *DepthFetcher) ping() error {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	_, err := client.Get(pingUrl)
	if err != nil {
		return fmt.Errorf("unsuccessful get request: %v", err)
	}
	return nil
}

// GetSymbols returns a slice of symbols that fetchers collect.
func (f *DepthFetcher) GetSymbols() []string {
	return f.symbols
}

// GetDataChannel returns the channel into which data about the specified character is written.
func (f *DepthFetcher) GetDataChannel(symbol string) (<-chan *domain.Depth, error) {
	symbol = normalizeSymbol(symbol)
	ch, exists := f.dataChannels[symbol]
	if !exists {
		return nil, fmt.Errorf("symbol %s not found", symbol)
	}
	return ch, nil
}

// RunFetchers launches fetchers for each symbol.
func (f *DepthFetcher) RunFetchers(wg *sync.WaitGroup) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	for symbol, ch := range f.dataChannels {
		wg.Add(1)
		go func(symbol string, dataCh chan *domain.Depth) {
			defer wg.Done()
			f.runFetcher(ctx, symbol, dataCh)
		}(symbol, ch)
	}
	return ctx, cancel
}

// runFetcher needs to be run in a goroutine.
// Writes the results of requests to get the order book into data channels.
func (f *DepthFetcher) runFetcher(ctx context.Context, symbol string, dataCh chan *domain.Depth) {
	ticker := time.NewTicker(f.fetcherRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			depth, err := fetchDepth(symbol, f.limit)
			if err != nil {
				log.Printf("fetchDepth: %v", err)
				return
			}
			select {
			case dataCh <- depth:
			default:
				// This is done if no one is reading from the channel, and it is full.
				// In this case, we read from it ourselves and write a new value.
				// If someone starts reading from the channel, they will receive more recent data.
				<-dataCh
				dataCh <- depth
			}
		case <-ctx.Done():
			return
		}
	}
}

// fetchDepth makes a request to receive the order book.
func fetchDepth(symbol string, limit int) (*domain.Depth, error) {
	reqUrl := fmt.Sprintf(depthUrl, symbol, limit)
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful get request with %d status code: %v", resp.StatusCode, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("response body closing error: %v", err)
		}
	}()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	depth, err := parseJSON(bytes)
	if err != nil {
		return nil, fmt.Errorf("parseJSON: %v", err)
	}
	return depth, nil
}
