package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secretKey []byte
}

type JWTService interface {
	Generate(userID string) (string, error)
	Parse(token string) (*jwt.Token, error)
}

func NewJWTService(secretKey []byte) JWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (j *jwtService) Generate(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * 7 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtService) Parse(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
}
