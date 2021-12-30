package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	controllers "github.com/lakhinsu/urlshortening-authservice/controllers"
)

func SetupRouters(app *gin.Engine) error {
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
