package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
)

func RunServer(config *config.Config) error {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "World",
		})
		return
	})

	port := fmt.Sprintf(":%s", config.ServerPort)
	return router.Run(port)
}
