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

type repos struct {
	userRepo repo.UserRepository
}

type services struct {
	userService service.UserService
}

type handlers struct {
	userHandler handler.UserHandler
}

func RunServer(config *config.Config, db *gorm.DB) error {
	router := InitRoutes(db)
	port := fmt.Sprintf(":%s", config.ServerPort)
	return router.Run(port)
}

func InitRoutes(db *gorm.DB) *gin.Engine {
	_, _, handlers := setDependencies(db)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "World",
		})
		return
	})
	router.POST("/login", handlers.userHandler.TwitterOAuthLogin)
	return router
}

func setDependencies(db *gorm.DB) (*repos, *services, *handlers) {
	repos := &repos{}
	repos.setDependencies(db)

	services := &services{}
	services.setDependencies(repos)

	handlers := &handlers{}
	handlers.setDependencies(services)

	return repos, services, handlers
}

func (r *repos) setDependencies(db *gorm.DB) {
	userRepo := repo.NewUserRepository(db)
	r.userRepo = userRepo
}

func (s *services) setDependencies(repos *repos) {
	userService := service.NewUserService(repos.userRepo)
	s.userService = userService
}

func (h *handlers) setDependencies(services *services) {
	userHandler := handler.NewUserHandler(services.userService)
	h.userHandler = userHandler
}
