package models

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TwitterID         string `gorm:"uniqueIndex" json:"twitter_id"`
	TwitterScreenName string `gorm:"uniqueIndex" json:"twitter_screen_name"`
}

type UserRepo struct {
	db *gorm.DB
}

func (ur *UserRepo) Create(u *User) error {
	return ur.CreateWithDB(ur.db, u)
}

func (ur *UserRepo) CreateWithDB(tx *gorm.DB, u *User) error {
	err := tx.Create(u).Error
	if err != nil {
		log.Errorf("unable to create user record in db, error: %v", err)
		return err
	}
	return nil
}