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
