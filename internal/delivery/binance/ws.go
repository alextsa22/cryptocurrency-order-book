package binance

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// wsDepthDelivery delivers via WebSocket.
func (p *DeliveryProvider) wsDepthDelivery(ctx context.Context, symbol string, limit int, dataCh chan *domain.Depth) {
	// When connecting via web socket to Binance, you need to use a lowercase symbol.
	lowerSymbol := strings.ToLower(symbol)
	streamName := fmt.Sprintf(p.apiConfig.StreamDepthName, lowerSymbol, limit)
	u := url.URL{
		Scheme: p.apiConfig.Scheme,
		Host:   p.apiConfig.Host,
		Path:   fmt.Sprintf(p.apiConfig.Path, streamName),
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("connection error: %v", err)
		close(dataCh)
		return
	}
	defer conn.Close()

	var depth domain.Depth
	for {
		err := conn.ReadJSON(&depth)
		if err != nil {
			log.Printf("read json error: %v", err)
			close(dataCh)
			return
		}
		select {
		case <-ctx.Done():
			log.Infof("order book provider for %s has been successfully completed", symbol)
			return
		case dataCh <- &depth:
		default:
			// This is done if no one is reading from the channel, and it is full.
			// In this case, we read from it ourselves and write a new value.
			// If someone starts reading from the channel, they will receive more recent data.
			<-dataCh
			dataCh <- &depth
		}
	}
}
