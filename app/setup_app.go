package app

import (
	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/middlewares"
	"github.com/lakhinsu/urlshortening-authservice/routers"
	"github.com/rs/zerolog/log"
)

func SetupApp() *gin.Engine {
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
		log.Fatal().Err(err).Msg("Error occurred while setting up routers")
		panic("Error occurred while setting up the routers")
	}

	return app
}
