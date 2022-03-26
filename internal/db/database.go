package db

import (
	"fmt"

	"github.com/skamranahmed/markeventas/config"
	"github.com/skamranahmed/markeventas/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Calcutta",
		config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("unable to open db connection, error: %s", err)
		return nil, err
	}
	return db, nil
}
