package cryptography

import (
	"pet-store/config"

	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessToken = iota
	RefreshToken
)

type TokenManager interface {
	GetAccessSecret() *jwtauth.JWTAuth
	CreateToken(userID, groups string, ttl time.Duration, kind int) (string, error)
}

type TokenJWT struct {
	AccessSecret  *jwtauth.JWTAuth
	RefreshSecret []byte
}

func (t *TokenJWT) GetAccessSecret() *jwtauth.JWTAuth {
	return t.AccessSecret
}

func NewTokenJWT(token config.Token) TokenManager {
	accSec := jwtauth.New("HS256", []byte(token.AccessSecret), nil)
	return &TokenJWT{AccessSecret: accSec, RefreshSecret: []byte(token.RefreshSecret)}
}

// UserClaims include custom claims on jwt.
type UserClaims struct {
	ID     string `json:"uid"`
	Role   string `json:"role"`
	Groups string `json:"groups"`
	Layers string `json:"layers"`
	jwt.RegisteredClaims
}

type UserFromClaims struct {
	ID     int
	Role   int
	Groups []int
	Layers []int
}

// CreateToken create new token with parameters.
func (o *TokenJWT) CreateToken(userID, groups string, ttl time.Duration, kind int) (string, error) {
	// Создание данных для токена
	claims := map[string]interface{}{
		"user_id": userID,
		"groups":  groups,
		"kind":    kind,
		"exp":     time.Now().Add(ttl).Unix(), // Устанавливаем время истечения токена
	}

	// Генерация токена с указанными данными
	_, token, err := o.AccessSecret.Encode(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
