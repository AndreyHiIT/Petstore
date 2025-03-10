package component

import (
	"pet-store/config"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/infrastructure/tools/cryptography"

	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Conf         config.AppConf
	TokenManager cryptography.TokenManager
	Responder    responder.Responder
	Decoder      godecoder.Decoder
	Logger       *zap.Logger
	Hash         cryptography.Hasher
}

func NewComponents(conf config.AppConf, tokenManager cryptography.TokenManager, responder responder.Responder, decoder godecoder.Decoder, hash cryptography.Hasher, logger *zap.Logger) *Components {
	return &Components{Conf: conf, TokenManager: tokenManager, Responder: responder, Decoder: decoder, Hash: hash, Logger: logger}
}
