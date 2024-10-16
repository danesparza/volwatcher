package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/volwatcher/internal/app"
	"github.com/volwatcher/internal/monitor"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	app.InitLogging()
	appConfig := app.InitConfig()

	//	Trap program exit appropriately
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 2)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(ctx, sigs, cancel)

	//	Start the background process
	go monitor.Process(ctx, *appConfig)

	//	Wait for our signal and shutdown gracefully
	<-ctx.Done()
}

func handleSignals(ctx context.Context, sigs <-chan os.Signal, cancel context.CancelFunc) {
	select {
	case <-ctx.Done():
	case sig := <-sigs:
		switch sig {
		case os.Interrupt:
			log.Info().Str("signal", "SIGINT").Msg("Shutting down")
		case syscall.SIGTERM:
			log.Info().Str("signal", "SIGTERM").Msg("Shutting down")
		}

		cancel()
		os.Exit(0)
	}
}
