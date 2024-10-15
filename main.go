package main

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
	"github.com/volwatcher/internal/app"
	"github.com/volwatcher/internal/folder"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	ctx := context.Background()

	//	Set log info:
	log.Logger = log.With().Timestamp().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	//	Set log level (default to info)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch strings.ToLower(os.Getenv("LOGGER_LEVEL")) {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		break
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		break
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		break
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		break
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		break
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		break
	}

	//	Set the error stack marshaller
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	//	Set log time format
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log.Debug().Msg("Starting check")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal().Err(err).Msg("could not find home directory")
	}

	//	Set locations to look for the config information
	viper.AddConfigPath(home)         // adding home directory as first search path
	viper.AddConfigPath(".")          // also look in the working directory
	viper.SetConfigName("volwatcher") // name the config file (without extension)
	viper.AutomaticEnv()              // read in environment variables that match

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}

	//	Unmarshal config into config struct
	var C app.Config
	err = viper.Unmarshal(&C)
	if err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal config")
	}

	//	For each volume entry, see if it exists:
	for _, item := range C.Volumes {
		log.Debug().Str("folder", item.Folder).Msg("Checking folder")

		//	If it doesn't exist, run mount it with the command listed
		if folder.DoesNotExist(item.Folder) {
			log.Info().Str("folder", item.Folder).Str("mountscript", item.MountScript).Msg("Folder does not exist.  Running mount command")

			//	Do we need to run this with some retry logic?
			//	Example: osascript -e 'mount volume "smb://username:password@hostname.tld/MySpecialSharename"'
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
	}

	log.Debug().Msg("Shutting down")

}
