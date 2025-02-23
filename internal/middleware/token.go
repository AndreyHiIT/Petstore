package middleware

import (
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/infrastructure/tools/cryptography"
)

type Token struct {
	responder.Responder
	jwt cryptography.TokenManager
}

func NewTokenManager(responder responder.Responder, jwt cryptography.TokenManager) *Token {
	return &Token{
		Responder: responder,
		jwt:       jwt,
	}
}
