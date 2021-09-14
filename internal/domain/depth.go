package domain

import (
	"strconv"

	log "github.com/sirupsen/logrus"
)

// Depth is depth of market with bids and asks.
type Depth struct {
	LastUpdateId int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

func (d *Depth) SumOfBidsQuantity() float64 {
	sum, err := sumOfOrderQuantities(d.Asks)
	if err != nil {
		log.Errorf("sumOfOrderQuantities: %v", err)
	}
	return sum
}

func (d *Depth) SumOfAsksQuantity() float64 {
	sum, err := sumOfOrderQuantities(d.Asks)
	if err != nil {
		log.Errorf("sumOfOrderQuantities: %v", err)
	}
	return sum
}

func sumOfOrderQuantities(orders [][]string) (float64, error) {
	var sum float64
	for _, order := range orders {
		quantity, err := strconv.ParseFloat(order[1], 64)
		if err != nil {
			return 0, err
		}
		sum += quantity
	}
	return sum, nil
}
