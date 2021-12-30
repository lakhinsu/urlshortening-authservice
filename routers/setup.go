package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/controllers"
	"github.com/rs/zerolog/log"
)

func SetupRouters(app *gin.Engine) error {

	log.Debug().Msg("Setting up V1 routers")
	v1 := app.Group("/v1")
	{
		v1.POST("/login", controllers.Login)
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello world",
			})
		})
	}
	return nil
}
