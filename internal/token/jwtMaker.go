package token

import (
	"github.com/dgrijalva/jwt-go"
)

type jwtMaker struct {
	secretSigningKey string
}

func NewJwtTokenMaker(secretSigningKey string) Maker {
	return &jwtMaker{
		secretSigningKey: secretSigningKey,
	}
}

func (maker *jwtMaker) CreateToken(userID uint, twitterID string) (string, error) {
	payload := NewPayload(userID, twitterID)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretSigningKey))
}

func (maker *jwtMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			// this means the signing algorithm of the token that we have got doesn't match with our signing algorithm
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretSigningKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && (verr.Inner.Error() == ErrExpiredToken.Error()) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
