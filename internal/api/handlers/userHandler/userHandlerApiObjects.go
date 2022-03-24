package handler

type TwitterOAuthCallbackPayload struct {
	OAuthToken    string `json:"oauth_token"`
	OAuthVerifier string `json:"oauth_verifier"`
}

type TwitterOAuthCallbackResponse struct {
	AccessToken              string `json:"access_token"`
	TwitterID                string `json:"twitter_id"`
	ScreenName               string `json:"screen_name"`
	IsGoogleOauthTokenActive bool   `json:"is_google_oauth_token_active"`
}

type GoogleApiCodePayload struct {
	Code string `json:"code" binding:"required"`
}

type UserProfileResponse struct {
	ID                       uint   `json:"id"`
	TwitterID                string `json:"twitter_id"`
	TwitterScreenName        string `json:"twitter_screen_name"`
	IsGoogleOauthTokenActive bool   `json:"is_google_oauth_token_active"`
}
