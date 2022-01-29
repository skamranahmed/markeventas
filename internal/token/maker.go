package token

type Maker interface {
	CreateToken(userID uint, twitterID string) (string, error)
	VerifyToken(token string) (*Payload, error)
}