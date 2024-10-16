package app

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Volumes []VolumeInfo
}

type VolumeInfo struct {
	Folder      string
	MountScript string
	AfterMount  string
}

func InitConfig() *Config {
	retval := &Config{}

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
	err = viper.Unmarshal(retval)
	if err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal config")
	}

	log.Info().Any("config", retval).Msg("config loaded")

	return retval
}
