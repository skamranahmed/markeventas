package handler

// RefreshTokenCodeRequest : 
type RefreshTokenCodeRequest struct {
	// we use this code to fetch the refresh token
	// works only for the first time per-account
	Code string `json:"code" binding:"required"`
}
