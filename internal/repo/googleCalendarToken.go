package repo

import (
	"github.com/skamranahmed/markeventas/internal/models"
	"gorm.io/gorm"
)

func NewGoogleCalendarTokenRepository(db *gorm.DB) GoogleCalendarTokenRepository {
	return &googleCalendartokenRepo{
		db: db,
	}
}

type googleCalendartokenRepo struct {
	db *gorm.DB
}

func (tr *googleCalendartokenRepo) Create(u *models.GoogleCalendarToken) error {
	return tr.createWithDB(tr.db, u)
}

func (tr *googleCalendartokenRepo) createWithDB(tx *gorm.DB, u *models.GoogleCalendarToken) error {
	err := tx.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *googleCalendartokenRepo) Save(u *models.GoogleCalendarToken) error {
	return tr.saveWithDB(tr.db, u)
}

func (tr *googleCalendartokenRepo) saveWithDB(tx *gorm.DB, u *models.GoogleCalendarToken) error {
	err := tx.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *googleCalendartokenRepo) GetUserToken(userID uint) (*models.GoogleCalendarToken, error) {
	token := &models.GoogleCalendarToken{}
	err := tr.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}
