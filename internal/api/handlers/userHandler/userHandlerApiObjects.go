package handler

type UserTwitterOAuthLoginInput struct {
	TwitterID         string `json:"twitter_id"`
	TwitterScreenName string `json:"twitter_screen_name"`
}
