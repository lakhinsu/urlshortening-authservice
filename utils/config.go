package utils

import (
	"github.com/lakhinsu/urlshortening-authservice/consts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(consts.ENV_FILE)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Debug().Err(err).
			Msg("Error occurred while reading .env file, might fallback to OS env config")
	}
	// TODO: Implement this to work from os.env as well
	viper.AutomaticEnv()
}

func GetEnvVar(name string) string {
	if !viper.IsSet(name) {
		log.Debug().Msgf("Environment varaible %s is not set", name)
		return ""
	}
	value := viper.GetString(name)
	return value
}
