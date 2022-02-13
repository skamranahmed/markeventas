package service

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/repo"
)

func NewBotReplyService(botReplyRepo repo.BotReplyRepository) BotReplyService {
	return &botReplyService{
		repo: botReplyRepo,
	}
}

type botReplyService struct {
	repo repo.BotReplyRepository
}

func (b *botReplyService) GetLatestRepliedToTweetID() (int64, error) {
	return b.repo.GetLatestRepliedToTweetID()
}

func (b *botReplyService) Insert(botReply *models.BotReply) (*models.BotReply, error) {
	return b.repo.Create(botReply)
}

func (b *botReplyService) Save(botReply *models.BotReply) error {
	return b.repo.Save(botReply)
}
