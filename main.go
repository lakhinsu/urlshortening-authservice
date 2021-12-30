package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lakhinsu/urlshortening-authservice/middlewares"
	"github.com/lakhinsu/urlshortening-authservice/routers"
	"github.com/lakhinsu/urlshortening-authservice/utils"
)

func init() {
	mode := utils.ReadEnvVar("GIN_MODE")
	gin.SetMode(mode)
}

func main() {
	app := gin.Default()
	app.SetTrustedProxies(nil)

	app.Use(middlewares.RequestID())

	err := routers.SetupRouters(app)

	if err != nil {
		panic("Error occured while setup up the routers")
	}

	// https := utils.ReadEnvVar("GIN_HTTPS")
	// if https == "True" {
	// 	m := autocert.Manager{
	// 		Prompt: autocert.AcceptTOS,
	// 		Cache:  autocert.DirCache("certs"),
	// 	}

	// 	autotls.RunWithManager(app, &m)
	// }

	host := utils.ReadEnvVar("GIN_HOST")
	port := utils.ReadEnvVar("GIN_PORT")
	app.Run(fmt.Sprintf("%s:%s", host, port))
}
