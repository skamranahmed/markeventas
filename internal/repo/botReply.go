package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"gorm.io/gorm"
)

func NewBotReplyRepository(db *gorm.DB) BotReplyRepository {
	return &botReplyRepo{
		db: db,
	}
}

type botReplyRepo struct {
	db *gorm.DB
}

func (br *botReplyRepo) Create(botReply *models.BotReply) (*models.BotReply, error) {
	return br.createWithDB(br.db, botReply)
}

func (br *botReplyRepo) createWithDB(tx *gorm.DB, b *models.BotReply) (*models.BotReply, error) {
	err := tx.Create(b).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (br *botReplyRepo) Save(b *models.BotReply) error {
	return br.saveWithDB(br.db, b)
}

func (br *botReplyRepo) saveWithDB(tx *gorm.DB, b *models.BotReply) error {
	err := tx.Save(b).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *botReplyRepo) GetLatestRepliedToTweetID() (int64, error) {
	type result struct {
		ToTweetID int64 `json:"to_tweet_id"`
	}
	var res result
	err := br.db.Raw(`
		SELECT
			to_tweet_id
		FROM
			bot_replies
		ORDER BY
			id DESC
		LIMIT 1
	`).Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res.ToTweetID, nil
}
