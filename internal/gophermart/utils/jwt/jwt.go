package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

const (
	TOKENEXPIRES = 1 * time.Hour
	SECRET       = "secretest key"
)

func GenerateToken(login string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXPIRES)),
		},
		Login: login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(SECRET))
}

func ParseToken(tokenString string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if claims.Valid() != nil {
		return "", claims.Valid()
	}
	return claims.Login, err
}

func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
