package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/app"
	"github.com/lakhinsu/urlshortening-authservice/utils"
	"github.com/rs/zerolog/log"
)

func init() {
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {

	app := app.SetupApp()

	host := utils.GetEnvVar("GIN_ADDR")
	port := utils.GetEnvVar("GIN_PORT")

	https := utils.GetEnvVar("GIN_HTTPS")
	if https == "true" {
		certFile := utils.GetEnvVar("GIN_CERT")
		certKey := utils.GetEnvVar("GIN_CERT_KEY")

		if err := app.RunTLS(fmt.Sprintf("%s:%s", host, port), certFile, certKey); err != nil {
			log.Error().Err(err).Msg("Error occured while setting up the server in HTTPS mode")
		}
	}

	log.Debug().Msgf("Listening on addr:%s and port:%s", host, port)
	if err := app.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Error().Err(err).Msg("Error occured while setting up the server")
	}
}
