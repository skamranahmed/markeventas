package cmd

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/db"
)

// Run : intializes our application
func Run() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Errorf("unable to load config file, error: %v\n", err)
		return err
	}

	// Init database
	log.Info("connecting to db.....")
	_, err = db.Init(cfg)
	if err != nil {
		return err
	}
	log.Info("✅ db connection successful")
	
	return api.RunServer(cfg)
}
