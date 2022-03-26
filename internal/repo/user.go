package repo

import (
	"github.com/skamranahmed/markeventas/internal/models"
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
	return ur.createWithDB(ur.db, u)
}

func (ur *userRepo) createWithDB(tx *gorm.DB, u *models.User) error {
	err := tx.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) Save(u *models.User) error {
	return ur.saveWithDB(ur.db, u)
}

func (ur *userRepo) saveWithDB(tx *gorm.DB, u *models.User) error {
	err := tx.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *userRepo) FindByTwitterID(twitterID string) (*models.User, error) {
	user := &models.User{}
	err := ur.db.Where("twitter_id = ?", twitterID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
