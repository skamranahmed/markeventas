package cmd

import (
	_ "github.com/skamranahmed/twitter-create-gcal-event-api/config"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/api"
	database "github.com/skamranahmed/twitter-create-gcal-event-api/internal/db"
)

// Run : intializes our application
func Run() error {
	// Init database
	log.Info("⏳ connecting to db.....")
	db, err := database.Init()
	if err != nil {
		return err
	}
	log.Info("✅ db connection successful")

	// Migrate db schema
	log.Info("🏃‍♂️ running db migrations")
	err = database.Migrate(db)
	if err != nil {
		return err
	}
	return api.RunServer(db)
}
