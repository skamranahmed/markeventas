package api

import (
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

func InitRoutes(db *gorm.DB, config *config.Config) *gin.Engine {
	tokenMaker = token.NewJwtTokenMaker(config.TokenSecretSigningKey)
	_, services, handlers := setDependencies(db, config)
	router := gin.Default()
	router.GET("/api/login", handlers.userHandler.TwitterOAuthLogin)
	router.POST("/api/twitter/callback", handlers.userHandler.HandleTwitterOAuthCallback)

	authorized := router.Group("/api/google/calendar")
	authorized.Use(middlewares.AuthMiddleware(tokenMaker))
	{
		// authorized.POST("/auth-code", handlers.userHandler.SaveGoogleCalendarRefreshToken)
		authorized.POST("/auth-code", handlers.googleTokenHandler.SaveRefreshToken)
	}

	// run the twitter bot in background
	go startTwitterBot(config, services.userService, services.botLogService)

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

	tokenRepo := repo.NewTokenRepository(db)
	r.tokenRepo = tokenRepo

	botLogRepo := repo.NewBotLogRepository(db)
	r.botLogRepo = botLogRepo
}

func (s *services) setDependencies(repos *repos, config *config.Config) {
	userService := service.NewUserService(repos.userRepo, repos.tokenRepo, config, tokenMaker)
	s.userService = userService

	botLogService := service.NewBotLogService(repos.botLogRepo)
	s.botLogService = botLogService
}

func (h *handlers) setDependencies(services *services, config *config.Config) {
	// init user handler
	userHandler := uHandler.NewUserHandler(services.userService, config)
	h.userHandler = userHandler

	// init googleToken handler
	googleTokenHandler := gHandler.NewGoogleTokenHandler(services.userService)
	h.googleTokenHandler = googleTokenHandler
}
