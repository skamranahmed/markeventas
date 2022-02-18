package repo

import (
	"github.com/skamranahmed/twitter-create-gcal-event-api/internal/models"
	"gorm.io/gorm"
)

func NewBotLogRepository(db *gorm.DB) BotLogRepository {
	return &botLogRepo{
		db: db,
	}
}

type botLogRepo struct {
	db *gorm.DB
}

func (br *botLogRepo) Create(botLog *models.BotLog) (*models.BotLog, error) {
	return br.createWithDB(br.db, botLog)
}

func (br *botLogRepo) createWithDB(tx *gorm.DB, b *models.BotLog) (*models.BotLog, error) {
	err := tx.Create(b).Error
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (br *botLogRepo) Save(b *models.BotLog) error {
	return br.saveWithDB(br.db, b)
}

func (br *botLogRepo) saveWithDB(tx *gorm.DB, b *models.BotLog) error {
	err := tx.Save(b).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *botLogRepo) GetLatestRepliedToTweetID() (int64, error) {
	type result struct {
		ToTweetID int64 `json:"to_tweet_id"`
	}
	var res result
	err := br.db.Raw(`
		SELECT
			to_tweet_id
		FROM
			bot_logs
		ORDER BY
			id DESC
		LIMIT 1
	`).Scan(&res).Error
	if err != nil {
		return 0, err
	}

	return res.ToTweetID, nil
}
