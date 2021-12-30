package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	logLevel := ReadEnvVar("LOG_LEVEL")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	host, err := os.Hostname()
	if err != nil {
		log.Logger = log.With().Str("host", "unknown").Logger()
	} else {
		log.Logger = log.With().Str("host", host).Logger()
	}

	log.Logger = log.With().Str("service", "urlshortening").Logger()

	log.Logger = log.With().Caller().Logger()
}
