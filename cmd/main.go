package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/config/binance"
	providers "github.com/alextsa22/cryptocurrency-order-book/internal/delivery/binance"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
	"github.com/alextsa22/cryptocurrency-order-book/internal/service"
	log "github.com/sirupsen/logrus"
)

var (
	trackSymbol = flag.String("track", "", "sets the symbol to be displayed")
)

func main() {
	flag.Parse()

	config, err := config.InitServiceConfig()
	if err != nil {
		log.Fatalf("config initialization error: %v", err)
	}
	log.Infoln("config initialization completed")

	apiConfig, err := binance.InitAPIConfig()
	if err != nil {
		log.Fatalf("binance api config initialization error: %v", err)
	}
	log.Infof("binance api config initialization complete")

	// If you want to use websockets for delivery,
	// just replace provider.RestMethod with provider.WebSocket.
	var method providers.DeliveryMethod = providers.RestMethod
	provider := providers.NewDeliveryProvider(config, apiConfig, method)
	deliveryService := service.NewDeliveryService(config, provider)

	log.Println("we initialize the process of launching providers")
	wg := &sync.WaitGroup{}
	_, cancel := deliveryService.RunProviders(wg)
	if err != nil {
		log.Fatalf("cant run fetchers with error: %v", err)
	}

	if *trackSymbol == "" {
		trackSymbol = &config.Symbols[0]
	}
	*trackSymbol = strings.TrimSpace(*trackSymbol)
	*trackSymbol = strings.ToUpper(*trackSymbol)
	log.Infof("%s symbol successfully set as tracked", *trackSymbol)

	dataCh, err := deliveryService.GetDataChannel(*trackSymbol)
	if err != nil {
		log.Errorf("error receiving data channel for the %s symbol: %v", *trackSymbol, err)
		cancel()
		wg.Wait()
		log.Println("all fetchers stopped")
		return
	}
	log.Infof("data channel for %s successfully received", *trackSymbol)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("start application...")

LOOP:
	for {
		select {
		case <-quit:
			log.Infoln("graceful shutdown")
			cancel()
			wg.Wait()
			log.Println("all providers have successfully completed their work")
			break LOOP
		case depth, opened := <-dataCh:
			if opened {
				printDepth(*trackSymbol, depth)
			} else {
				log.Infoln("stop all providers")
				cancel()
				wg.Wait()
				log.Println("all providers stopped")
				break LOOP
			}
		}
	}
}

func printDepth(symbol string, depth *domain.Depth) {
	fmt.Println()
	fmt.Printf("Symbol: %s    Last Update: %d\n", symbol, depth.LastUpdateId)

	fmt.Println("Bids:")
	printOrders(depth.Bids)
	sum := depth.SumOfBidsQuantity()
	fmt.Printf("\tSum of order quantity: %.5f\n", sum)

	fmt.Println("Asks:")
	printOrders(depth.Asks)
	sum = depth.SumOfAsksQuantity()
	fmt.Printf("\tSum of order quantity: %.5f\n", sum)

	fmt.Println()
}

func printOrders(orders [][]string) {
	for _, order := range orders {
		price, _ := strconv.ParseFloat(order[0], 64)
		quantity, _ := strconv.ParseFloat(order[1], 64)
		fmt.Printf("\tPrice: %.2f    Quantity: %.5f\n", price, quantity)
	}
}
