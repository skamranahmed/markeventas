package handler

type TwitterOAuthCallbackPayload struct {
	OAuthToken    string `json:"oauth_token"`
	OAuthVerifier string `json:"oauth_verifier"`
}

type TwitterOAuthCallbackResponse struct {
	AccessToken string `json:"access_token"`
	TwitterID   string `json:"twitter_id"`
	ScreenName  string `json:"screen_name"`
}