package monitor

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/volwatcher/internal/app"
	"github.com/volwatcher/internal/folder"
	"os/exec"
	"strings"
)

// CheckWithConfig performs a check of all configured volumes based on the
// passed configuration information
func CheckWithConfig(ctx context.Context, config app.Config) {
	//	For each volume entry, see if it exists:
	for _, item := range config.Volumes {
		log.Debug().Str("folder", item.Folder).Msg("Checking folder")

		//	If it doesn't exist, run mount it with the command listed
		if folder.DoesNotExist(item.Folder) {
			log.Info().Str("folder", item.Folder).Str("mountscript", item.MountScript).Msg("Folder does not exist.  Running mount command")

			cmdMount := exec.CommandContext(ctx, "/bin/bash", "-c", item.MountScript)
			err := cmdMount.Run()
			if err != nil {
				log.Err(err).Str("folder", item.Folder).Msg("Mount script failed")
				continue
			}

			//	Check it again
			if folder.DoesExist(item.Folder) {
				log.Info().Str("folder", item.Folder).Msg("Folder exists now")

				//	If it exists now and if we have an AfterMount script, run it
				//	Example: docker compose -f /path/to/your/docker-compose.yml restart
				if len(strings.TrimSpace(item.AfterMount)) > 0 {
					log.Info().Str("folder", item.Folder).Msg("Running after mount script")
					cmdAfterMount := exec.CommandContext(ctx, "/bin/bash", "-c", item.AfterMount)
					err := cmdAfterMount.Run()
					if err != nil {
						log.Err(err).Str("folder", item.Folder).Msg("After mount script failed")
						continue
					}
				}
			}
		}

		log.Debug().Str("folder", item.Folder).Msg("Finished checking folder")
	}
}
