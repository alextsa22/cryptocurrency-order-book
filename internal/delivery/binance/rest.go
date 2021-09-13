package binance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	log "github.com/sirupsen/logrus"
)

// restDepthDelivery delivers via REST.
func (p *DeliveryProvider) restDepthDelivery(ctx context.Context, symbol string, limit int, dataCh chan *domain.Depth) {
	log.Infof("order book fetcher for %s has been successfully launched", symbol)

	ticker := time.NewTicker(p.rate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			depth, err := p.getDepth(symbol, limit)
			if err != nil {
				log.Errorf("getDepth: %v", err)
				log.Infof("order book fetcher for %s stopped", symbol)
				close(dataCh)
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
			log.Infof("order book fetcher for %s has been successfully completed", symbol)
			return
		}
	}
}

// getDepth makes a request to receive the order book.
func (p *DeliveryProvider) getDepth(symbol string, limit int) (*domain.Depth, error) {
	reqUrl := fmt.Sprintf(p.apiConfig.GetDepthUrl, symbol, limit)
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("unsuccessful get request with %d status code: %v", resp.StatusCode, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Errorf("response body closing error: %v", err)
		}
	}()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("bad reguest with this '%s' symbol and this limit %d", symbol, limit)
	}

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
