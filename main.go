package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	<-shutdownCtx.Done()
	log.Println("Graceful shutdown success!")
}
