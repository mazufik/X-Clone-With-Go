package main

import (
	"GO-SOCMED/config"
	"GO-SOCMED/router"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.LoadDB()

	app := gin.Default()
	api := app.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// router
	router.AuthRouter(api)

	app.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
