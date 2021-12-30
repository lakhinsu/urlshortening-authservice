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
	utils.SetupLogger()
	mode := utils.ReadEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	app := gin.New()

	app.Use(gin.Recovery())

	log.Info().Msg("Initializing service")

	// disabling the trusted proxy feature
	app.SetTrustedProxies(nil)

	log.Info().Msg("Adding request id middleware")

	app.Use(middlewares.RequestID(), middlewares.RequestLogger())

	log.Info().Msg("Setting up routers")

	err := routers.SetupRouters(app)

	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while setting up routers")
		panic("Error occured while setting up the routers")
	}

	// https := utils.ReadEnvVar("GIN_HTTPS")
	// if https == "true" {
	// 	m := autocert.Manager{
	// 		Prompt: autocert.AcceptTOS,
	// 		Cache:  autocert.DirCache("certs"),
	// 	}

	// 	autotls.RunWithManager(app, &m)
	// }

	host := utils.ReadEnvVar("GIN_ADDR")
	port := utils.ReadEnvVar("GIN_PORT")

	log.Debug().Msgf("Listening on addr:%s and port:%s", host, port)
	if err := app.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Error().Err(err).Msg("Error occured while setting up the server")
	}
}
