package service

import (
	"fmt"

	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
)

// NewUserService : returns a userService struct that implements the UserService interface
func NewUserService(userRepo repo.UserRepository) UserService {
	return &userService{
		repo: userRepo,
	}
}

type userService struct {
	repo repo.UserRepository
}

func (us *userService) Create(u *models.User) error {
	return fmt.Errorf("not implemented")
}
