package cmd

import (
	"flag"

	_ "github.com/skamranahmed/markeventas/config"
	"github.com/skamranahmed/markeventas/pkg/log"

	"github.com/skamranahmed/markeventas/internal/api"
	database "github.com/skamranahmed/markeventas/internal/db"
)

// Run : intializes our application
func Run() error {
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")

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
