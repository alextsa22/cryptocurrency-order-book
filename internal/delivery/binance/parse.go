package binance

import (
	"fmt"
	"strconv"

	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	"github.com/bitly/go-simplejson"
)

func parseJSON(data []byte) (*domain.Depth, error) {
	j, err := simplejson.NewJson(data)
	if err != nil {
		return nil, fmt.Errorf("error creating simplejson.Json: %v", err)
	}

	var depth domain.Depth
	depth.LastUpdateId = j.Get("lastUpdateId").MustInt()
	bids, err := parseOrders(j, "bids")
	if err != nil {
		return nil, fmt.Errorf("bids parsing error: %v", err)
	}
	asks, err := parseOrders(j, "asks")
	if err != nil {
		return nil, fmt.Errorf("asks parsing error: %v", err)
	}

	depth.Bids = bids
	depth.Asks = asks
	return &depth, nil
}

func parseOrders(j *simplejson.Json, key string) ([]*domain.Order, error) {
	ordersLen := len(j.Get(key).MustArray())
	orders := make([]*domain.Order, ordersLen)
	for i := 0; i < ordersLen; i++ {
		order := j.Get(key).GetIndex(i)

		orderPrice := order.GetIndex(0).MustString()
		price, err := strconv.ParseFloat(orderPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("price parsing error: %v", err)
		}

		orderQuantity := order.GetIndex(1).MustString()
		quantity, err := strconv.ParseFloat(orderQuantity, 64)
		if err != nil {
			return nil, fmt.Errorf("quantity parsing error: %v", err)
		}

		orders[i] = &domain.Order{
			Price:    price,
			Quantity: quantity,
		}
	}
	return orders, nil
}
