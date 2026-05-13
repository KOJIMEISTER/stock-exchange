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
		OrderType: matching.TypeBuy,
		Ammount:   10,
		Price:     100,
	})

	matchingEngine.PlaceOrder(matching.Order{
		Id:        1,
		Timestamp: time.Now(),
		OrderType: matching.TypeSell,
		Ammount:   10,
		Price:     100,
	})

	<-ctx.Done()

	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-shutdownCtx.Done()
	log.Println("Graceful shutdown success!")
}
