package handler

type TwitterOAuthCallbackPayload struct {
	OAuthToken    string `json:"oauth_token"`
	OAuthVerifier string `json:"oauth_verifier"`
}
