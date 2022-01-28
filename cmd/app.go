package cmd

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api"
)

// Run : intializes our application
func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("unable to load config file, error: %v\n", err)
		return err
	}
	return api.RunServer(cfg)
}
