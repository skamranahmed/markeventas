package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
)

// UserRepository
type UserRepository interface {
	Create(u *models.User) error
	Save(u *models.User) error
	FindByTwitterID(twitterID string) (*models.User, error)
}
