package db

import (
	"github.com/skamranahmed/markeventas/internal/models"
	"github.com/skamranahmed/markeventas/pkg/log"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Info("creating User table")
	err := db.AutoMigrate(models.User{})
	if err != nil {
		log.Errorf("unable to migrate User table, error: %v", err)
		return err
	}

	log.Info("creating GoogleCalendarToken table")
	err = db.AutoMigrate(models.GoogleCalendarToken{})
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

	log.Info("creating BotReplyType table")
	err = db.AutoMigrate(models.BotReplyType{})
	if err != nil {
		log.Errorf("unable to migrate BotReplyType table, error: %v", err)
		return err
	}

	err = SaveBotReplyTypes(db)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func SaveBotReplyTypes(db *gorm.DB) error {
	for _, botReplyType := range models.BotReplyTypes {
		toFind := models.BotReplyType{Code: botReplyType.Code}
		err := db.Where(toFind).FirstOrCreate(&botReplyType).Error
		if err != nil {
			log.Errorf("unable to save botReplyType: %+v in db, err: %v", botReplyType, err)
			return err
		}
	}
	return nil
}
