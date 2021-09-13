package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/alextsa22/cryptocurrency-order-book/internal/binance"
	"github.com/alextsa22/cryptocurrency-order-book/internal/config"
	"github.com/alextsa22/cryptocurrency-order-book/internal/domain"
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

	depthService, err := binance.NewDepthFetcher(config)
	if err != nil {
		log.Fatalf("error creating binance depthService: %v", err)
	}
	log.Infoln("fetcher service has completed initialization")

	log.Println("we initialize the process of launching fetchers")
	wg := &sync.WaitGroup{}
	_, cancel := depthService.RunFetchers(wg)
	if err != nil {
		log.Fatalf("cant run fetchers with error: %v", err)
	}

	if *trackSymbol == "" {
		trackSymbol = &config.Symbols[0]
	}
	*trackSymbol = strings.TrimSpace(*trackSymbol)
	*trackSymbol = strings.ToUpper(*trackSymbol)
	log.Infof("%s symbol successfully set as tracked", *trackSymbol)

	dataCh, err := depthService.GetDataChannel(*trackSymbol)
	if err != nil {
		cancel()
		log.Fatalf("error receiving data channel for the %s symbol: %v", *trackSymbol, err)
	}
	log.Infof("data channel for %s successfully received", *trackSymbol)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	log.Infof("start application...")
	fmt.Println()
LOOP:
	for {
		select {
		case <-quit:
			cancel()
			log.Infoln("quit the program")
			break LOOP
		case depth := <-dataCh:
			printDepth(*trackSymbol, depth)
		}
	}
	wg.Wait()
	log.Println("all fetchers have successfully completed their work")
}

func printDepth(symbol string, depth *domain.Depth) {
	fmt.Printf("Symbol: %s    Last Update: %d\n", symbol, depth.LastUpdateId)
	fmt.Println("Bids:")
	printOrders(depth.Bids)
	fmt.Printf("\tSum of order quantity: %.5f\n", depth.SumOfBidsQuantity())
	fmt.Println("Asks:")
	printOrders(depth.Asks)
	fmt.Printf("\tSum of order quantity: %.5f\n", depth.SumOfAsksQuantity())
}

func printOrders(orders []*domain.Order) {
	for _, order := range orders {
		fmt.Printf("\tPrice: %.2f    Quantity: %.5f\n", order.Price, order.Quantity)
	}
}
