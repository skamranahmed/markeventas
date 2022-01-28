package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
	"gorm.io/gorm"
)

// NewUserRepository : returns a userRepo struct that implements the UserRepository interface
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

type userRepo struct {
	db *gorm.DB
}

func (ur *userRepo) Create(u *models.User) error {
	return ur.CreateWithDB(ur.db, u)
}

func (ur *userRepo) CreateWithDB(tx *gorm.DB, u *models.User) error {
	err := tx.Create(u).Error
	if err != nil {
		log.Errorf("unable to create user record in db, error: %v", err)
		return err
	}
	return nil
}
