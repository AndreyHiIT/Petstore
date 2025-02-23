package service

import (
	"context"
	"net/http"
	"pet-store/config"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/tools/cryptography"
	"pet-store/internal/models"
	uservice "pet-store/internal/modules/user/service"
	"strconv"

	"go.uber.org/zap"
)

type Auth struct {
	conf         config.AppConf
	user         uservice.Userer
	tokenManager cryptography.TokenManager
	hash         cryptography.Hasher
	logger       *zap.Logger
}

func NewAuth(user uservice.Userer, components *component.Components) *Auth {
	return &Auth{
		conf:         components.Conf,
		user:         user,
		tokenManager: components.TokenManager,
		hash:         components.Hash,
		logger:       components.Logger,
	}
}

func (a *Auth) CreateUser(ctx context.Context, in CreateUserIn) CreateUserOut {
	hashPass, err := cryptography.HashPassword(in.Password)
	if err != nil {
		return CreateUserOut{
			Status:    http.StatusInternalServerError,
			ErrorCode: errors.HashPasswordError,
		}
	}

	userCreate := uservice.UserCreateIn{
		Email:     in.Email,
		Password:  hashPass,
		UserName:  in.Username,
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Phone:     in.Phone,
	}

	userOut := a.user.Create(ctx, userCreate)
	if userOut.ErrorCode != errors.NoError {
		if userOut.ErrorCode == errors.UserServiceUserAlreadyExists {
			return CreateUserOut{
				Status:    http.StatusConflict,
				ErrorCode: userOut.ErrorCode,
			}
		}
		return CreateUserOut{
			Status:    http.StatusInternalServerError,
			ErrorCode: userOut.ErrorCode,
		}
	}

	return CreateUserOut{
		Status:    http.StatusOK,
		ErrorCode: errors.NoError,
	}
}

func (a *Auth) AuthorizeEmail(ctx context.Context, in AuthorizeEmailIn) AuthorizeOut {
	// 1. получаем юзера по email
	userOut := a.user.GetByEmail(ctx, uservice.GetByEmailIn{Email: in.Email})
	if userOut.ErrorCode != errors.NoError {
		return AuthorizeOut{
			ErrorCode: userOut.ErrorCode,
		}
	}
	user := userOut.User

	// 2. проверяем пароль
	if !cryptography.CheckPassword(user.Password, in.Password) {
		return AuthorizeOut{
			ErrorCode: errors.AuthServiceWrongPasswordErr,
		}
	}

	// 3. генерируем токены
	accessToken, refreshToken, errorCode := a.generateTokens(user)
	if errorCode != errors.NoError {
		return AuthorizeOut{
			ErrorCode: errorCode,
		}
	}
	// 4. возвращаем токены
	return AuthorizeOut{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (a *Auth) generateTokens(user *models.User) (string, string, int) {
	accessToken, err := a.tokenManager.CreateToken(
		strconv.Itoa(user.ID),
		"",
		a.conf.Token.AccessTTL,
		cryptography.AccessToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", errors.AuthServiceAccessTokenGenerationErr
	}
	refreshToken, err := a.tokenManager.CreateToken(
		strconv.Itoa(user.ID),
		"",
		a.conf.Token.RefreshTTL,
		cryptography.RefreshToken,
	)
	if err != nil {
		a.logger.Error("auth: create access token err", zap.Error(err))
		return "", "", errors.AuthServiceRefreshTokenGenerationErr
	}

	return accessToken, refreshToken, errors.NoError
}
