package domain

// Depth is depth of market with bids and asks.
type Depth struct {
	LastUpdateId int      `json:"lastUpdateId"`
	Bids         []*Order `json:"bids"`
	Asks         []*Order `json:"asks"`
}

// Order is an order from depth of market.
type Order struct {
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

func (d *Depth) SumOfBidsQuantity() float64 {
	return sumOfOrderQuantities(d.Bids)
}

func (d *Depth) SumOfAsksQuantity() float64 {
	return sumOfOrderQuantities(d.Asks)
}

func sumOfOrderQuantities(orders []*Order) float64 {
	var sum float64
	for _, order := range orders {
		sum += order.Quantity
	}
	return sum
}
