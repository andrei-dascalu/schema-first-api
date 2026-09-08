// Package jwttoken implements app.TokenIssuer with HS256 JWTs.
package jwttoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"spec-first-backend/internal/app"
)

const TokenTTL = 24 * time.Hour

var ErrInvalidToken = errors.New("invalid token")

var _ app.TokenIssuer = (*Issuer)(nil)

type Claims struct {
	UserID uuid.UUID `json:"sub"`
	jwt.RegisteredClaims
}

type Issuer struct {
	secret []byte
}

func NewIssuer(secret string) *Issuer {
	return &Issuer{secret: []byte(secret)}
}

func (i *Issuer) Issue(userID uuid.UUID) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(TokenTTL)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(i.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

func (i *Issuer) Parse(token string) (uuid.UUID, error) {
	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return i.secret, nil
	})
	if err != nil || !parsed.Valid {
		return uuid.UUID{}, ErrInvalidToken
	}
	return claims.UserID, nil
}
