package db

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Info("creating User table")
	err := db.AutoMigrate(models.User{})
	if err != nil {
		log.Errorf("unable to migrate User table, error: %v", err)
		return err
	}

	log.Info("creating Token table")
	err = db.AutoMigrate(models.Token{})
	if err != nil {
		log.Errorf("unable to migrate Token table, error: %v", err)
		return err
	}

	log.Info("creating BotLog table")
	err = db.AutoMigrate(models.BotLog{})
	if err != nil {
		log.Errorf("unable to migrate BotLog table, error: %v", err)
		return err
	}
	return nil
}