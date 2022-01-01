package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/middlewares"
	"github.com/lakhinsu/urlshortening-authservice/routers"
	"github.com/lakhinsu/urlshortening-authservice/utils"
	"github.com/rs/zerolog/log"
)

func init() {
	mode := utils.GetEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	app := gin.New()

	app.Use(gin.Recovery())

	log.Info().Msg("Initializing service")

	// disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	log.Info().Msg("Adding request id middleware")

	app.Use(middlewares.RequestID(), middlewares.RequestLogger(), middlewares.CORSMiddleware())

	log.Info().Msg("Setting up routers")

	err := routers.SetupRouters(app)

	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while setting up routers")
		panic("Error occured while setting up the routers")
	}

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
