package token

import (
	"errors"
	"time"
)

var (
	ErrExpiredToken = errors.New("token has expired")

	ErrInvalidToken = errors.New("invalid token")
)

type Payload struct {
	UserID    uint      `json:"user_id"`
	TwitterID string    `json:"twitter_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPayload(userID uint, twitterID string) *Payload {
	return &Payload{
		UserID:    userID,
		TwitterID: twitterID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
