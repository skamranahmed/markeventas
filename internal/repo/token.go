package repo

import (
	"github.com/skamranahmed/markeventas/internal/models"
	"gorm.io/gorm"
)

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepo{
		db: db,
	}
}

type tokenRepo struct {
	db *gorm.DB
}

func (tr *tokenRepo) Create(u *models.Token) error {
	return tr.createWithDB(tr.db, u)
}

func (tr *tokenRepo) createWithDB(tx *gorm.DB, u *models.Token) error {
	err := tx.Create(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *tokenRepo) Save(u *models.Token) error {
	return tr.saveWithDB(tr.db, u)
}

func (tr *tokenRepo) saveWithDB(tx *gorm.DB, u *models.Token) error {
	err := tx.Save(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (tr *tokenRepo) GetUserToken(userID uint) (*models.Token, error) {
	token := &models.Token{}
	err := tr.db.Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		return nil, err
	}
	return token, nil
}
