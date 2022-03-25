package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	gHandler "github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/handlers/googletokenHandler"
	uHandler "github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/handlers/userHandler"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api/middlewares"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/service"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/token"
	"gorm.io/gorm"
)

var (
	tokenMaker token.Maker
)

type repos struct {
	userRepo   repo.UserRepository
	tokenRepo  repo.TokenRepository
	botLogRepo repo.BotLogRepository
}

type services struct {
	userService   service.UserService
	botLogService service.BotLogService
}

type handlers struct {
	userHandler        uHandler.UserHandler
	googleTokenHandler gHandler.GoogleTokenHandler
}

func InitRoutes(db *gorm.DB) *gin.Engine {
	tokenMaker = token.NewJwtTokenMaker(config.TokenSecretSigningKey)
	_, services, handlers := setDependencies(db)
	router := gin.Default()
	router.GET("/api/login", handlers.userHandler.TwitterOAuthLogin)
	router.POST("/api/twitter/callback", handlers.userHandler.HandleTwitterOAuthCallback)

	// health route
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "ok"})
		return
	})

	authorized := router.Group("/api")
	authorized.Use(middlewares.AuthMiddleware(tokenMaker))
	{
		authorized.POST("/save/google-calendar/tokens", handlers.googleTokenHandler.SaveRefreshToken)
		authorized.GET("/profile/me", handlers.userHandler.GetUserProfile)
	}

	// run the twitter bot in background
	go startTwitterBot(services.userService, services.botLogService)

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

	tokenRepo := repo.NewTokenRepository(db)
	r.tokenRepo = tokenRepo

	botLogRepo := repo.NewBotLogRepository(db)
	r.botLogRepo = botLogRepo
}

func (s *services) setDependencies(repos *repos) {
	userService := service.NewUserService(repos.userRepo, repos.tokenRepo, tokenMaker)
	s.userService = userService

	botLogService := service.NewBotLogService(repos.botLogRepo)
	s.botLogService = botLogService
}

func (h *handlers) setDependencies(services *services) {
	// init user handler
	userHandler := uHandler.NewUserHandler(services.userService)
	h.userHandler = userHandler

	// init googleToken handler
	googleTokenHandler := gHandler.NewGoogleTokenHandler(services.userService)
	h.googleTokenHandler = googleTokenHandler
}
