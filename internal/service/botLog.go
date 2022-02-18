package service

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
)

func NewBotLogService(botLogRepo repo.BotLogRepository) BotLogService {
	return &botLogService{
		repo: botLogRepo,
	}
}

type botLogService struct {
	repo repo.BotLogRepository
}

func (b *botLogService) GetLatestRepliedToTweetID() (int64, error) {
	return b.repo.GetLatestRepliedToTweetID()
}

func (b *botLogService) Insert(botLog *models.BotLog) (*models.BotLog, error) {
	return b.repo.Create(botLog)
}

func (b *botLogService) Save(botLog *models.BotLog) error {
	return b.repo.Save(botLog)
}
