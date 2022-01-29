package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"gorm.io/gorm"
)

// UserRepository
type UserRepository interface {
	Create(u *models.User) error
	CreateWithDB(tx *gorm.DB, u *models.User) error
}