package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/maxrasky/faraway/internal/client"
	"github.com/maxrasky/faraway/internal/pow"
)

func main() {
	cfg := client.NewConfig()
	powClient := client.NewClient(cfg.ServerAddr, pow.NewService())

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigint
		cancel()
	}()

	if err := powClient.Start(ctx); err != nil {
		log.Err(err).Send()
	}
}
