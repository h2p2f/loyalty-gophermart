package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

//func init() {
//	err := godotenv.Load("secret.env")
//	if err != nil {
//		panic(err)
//	}
//	Secret = os.Getenv("SECRET")
//}

// Claims - struct for claims
type Claims struct {
	jwt.RegisteredClaims
	Login string
}

//var Secret string = "secretest key"

// const for token
const (
	TOKENEXPIRES = 1 * time.Hour
)

// GenerateToken - generate token
func GenerateToken(login, key string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKENEXPIRES)),
		},
		Login: login,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(key))
}

// ParseToken - parse token
func ParseToken(tokenString, key string) (string, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if claims.Valid() != nil {
		return "", claims.Valid()
	}
	return claims.Login, err
}

// Valid - check if token is valid
func (c *Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
