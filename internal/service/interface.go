package service

import "github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"

type UserService interface {
	Create(u *models.User) error
}
