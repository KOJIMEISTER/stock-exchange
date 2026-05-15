package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kojimeister/stock-exchange/matching"
)

func main() {
	log.Println("Started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	matchingEngine := matching.NewMatchingEngine()

	matchingEngine.PlaceOrder(matching.Order{
		Id:        0,
		Timestamp: time.Now(),
		OrderType: matching.Limit,
		OrderSide: matching.Buy,
		Ammount:   10,
		Price:     100,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        0,
		Timestamp: time.Now(),
		OrderType: matching.Limit,
		OrderSide: matching.Buy,
		Ammount:   10,
		Price:     105,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        0,
		Timestamp: time.Now(),
		OrderType: matching.Limit,
		OrderSide: matching.Buy,
		Ammount:   10,
		Price:     95,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        1,
		Timestamp: time.Now(),
		OrderType: matching.Market,
		OrderSide: matching.Sell,
		Ammount:   10,
		Price:     100,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        1,
		Timestamp: time.Now(),
		OrderType: matching.Market,
		OrderSide: matching.Sell,
		Ammount:   10,
		Price:     105,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        1,
		Timestamp: time.Now(),
		OrderType: matching.Market,
		OrderSide: matching.Sell,
		Ammount:   10,
		Price:     95,
	})

	<-ctx.Done()

	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-shutdownCtx.Done()
	log.Println("Graceful shutdown success!")
}
