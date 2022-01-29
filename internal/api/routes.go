package api

import (
	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	handler "github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/handlers/userHandler"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"gorm.io/gorm"
)

type repos struct {
	userRepo repo.UserRepository
}

type services struct {
	userService service.UserService
}

type handlers struct {
	userHandler handler.UserHandler
}

func InitRoutes(db *gorm.DB, config *config.Config) *gin.Engine {
	_, _, handlers := setDependencies(db, config)
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Hello": "World",
		})
		return
	})
	router.GET("/api/login", handlers.userHandler.TwitterOAuthLogin)
	router.POST("/api/twitter/callback", handlers.userHandler.HandleTwitterOAuthCallback)
	return router
}

func setDependencies(db *gorm.DB, config *config.Config) (*repos, *services, *handlers) {
	repos := &repos{}
	repos.setDependencies(db)

	services := &services{}
	services.setDependencies(repos, config)

	handlers := &handlers{}
	handlers.setDependencies(services, config)

	return repos, services, handlers
}

func (r *repos) setDependencies(db *gorm.DB) {
	userRepo := repo.NewUserRepository(db)
	r.userRepo = userRepo
}

func (s *services) setDependencies(repos *repos, config *config.Config) {
	userService := service.NewUserService(repos.userRepo, config)
	s.userService = userService
}

func (h *handlers) setDependencies(services *services, config *config.Config) {
	userHandler := handler.NewUserHandler(services.userService, config)
	h.userHandler = userHandler
}
