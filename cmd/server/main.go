package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/maxrasky/faraway/internal/pow"
	"github.com/maxrasky/faraway/internal/quotes"
	"github.com/maxrasky/faraway/internal/server"
)

func main() {
	cfg := server.NewConfig()
	powServer := server.NewServer(cfg.Port, pow.NewService(), quotes.New())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sigint
		log.Info().Msg("interrupt signal received")
		cancel()
	}()

	if err := powServer.Run(ctx); err != nil {
		log.Err(err).Send()
	}
}
