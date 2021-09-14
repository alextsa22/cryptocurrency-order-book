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
	log.Infof("order book provider for %s has been successfully launched", symbol)
	// When connecting via web socket to Binance, you need to use a lowercase symbol.
	lowerSymbol := strings.ToLower(symbol)
	streamName := fmt.Sprintf(p.apiConfig.StreamDepthName, lowerSymbol, limit)
	u := url.URL{
		Scheme: p.apiConfig.Scheme,
		Host:   p.apiConfig.Host,
		Path:   fmt.Sprintf(p.apiConfig.Path, streamName),
	}

	// make a request to make sure the connection details are correct.
	_, err := p.getDepth(symbol, limit)
	if err != nil {
		log.Errorf("getDepth: %v", err)
		log.Infof("order book provider for %s stopped", symbol)
		close(dataCh)
		return
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Errorf("websocket connection error with default dialer: %v", err)
		log.Infof("order book provider for %s stopped", symbol)
		close(dataCh)
		return
	}
	defer conn.Close()

	for {
		select {
		case <-ctx.Done():
			log.Infof("order book provider for %s has been successfully completed", symbol)
			close(dataCh)
			return
		default:
			var depth domain.Depth
			err := conn.ReadJSON(&depth)
			if err != nil {
				log.Errorf("json read error: %v", err)
				log.Infof("order book provider for %s stopped", symbol)
				close(dataCh)
				return
			}
			select {
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
}
