package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	handler "github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/handlers/userHandler"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"gorm.io/gorm"
)

/*
	############## Order of Call Execution ############
						Handler
						   ↓
						Service
						   ↓
						Repository
	###################################################


	############## Dependency Injection ############
						Repository
						   ↓
						 Service
						   ↓
						 Handler
	###################################################
*/
var (
	db *gorm.DB
)

var (
	userRepo    repo.UserRepository = repo.NewUserRepository(db)
	userService service.UserService = service.NewUserService(userRepo)
	userHandler handler.UserHandler = handler.NewUserHandler(userService)
)

func RunServer(config *config.Config) error {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "World",
		})
		return
	})

	router.GET("/login", userHandler.TwitterOAuthLogin)

	port := fmt.Sprintf(":%s", config.ServerPort)
	return router.Run(port)
}
