package monitor

import "C"
import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/volwatcher/internal/app"
	"time"
)

func Process(ctx context.Context, config app.Config) {
	log.Info().Msg("Starting monitoring process")

	//	Check immediately:
	CheckWithConfig(ctx, config)

	//	Loop and respond to channels:
	for {
		select {
		case <-time.After(time.Minute * 1):
			CheckWithConfig(ctx, config)
		case <-ctx.Done():
			log.Info().Msg("Stopping monitoring process")
			return
		}
	}
}
