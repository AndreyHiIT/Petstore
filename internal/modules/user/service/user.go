package service

import (
	"context"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/tools/cryptography"
	"pet-store/internal/models"
	"pet-store/internal/modules/user/storage"

	"go.uber.org/zap"
)

type UserService struct {
	storage storage.Userer
	logger  *zap.Logger
}

func NewUserService(storage storage.Userer, logger *zap.Logger) *UserService {
	return &UserService{storage: storage, logger: logger}
}

func (u *UserService) Create(ctx context.Context, in UserCreateIn) UserCreateOut {
	var dto models.UserDTO
	dto.SetUserName(in.UserName).
		SetPhone(in.Phone).
		SetEmail(in.Email).
		SetPassword(in.Password).
		SetFirstName(in.FirstName).
		SetLastName(in.LastName)

	userID, err := u.storage.Create(ctx, dto)
	if err != nil {
		return UserCreateOut{
			ErrorCode: errors.UserServiceCreateUserErr,
		}
	}

	return UserCreateOut{
		UserID: userID,
	}
}

func (u *UserService) GetByEmail(ctx context.Context, in GetByEmailIn) UserOut {
	userDTO, err := u.storage.GetByEmail(ctx, in.Email)
	if err != nil {
		u.logger.Error("user: GetByEmail err", zap.Error(err))
		return UserOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}

	return UserOut{
		User: &models.User{
			ID:        userDTO.GetID(),
			Username:  userDTO.GetUserName(),
			Phone:     userDTO.GetPhone(),
			Email:     userDTO.GetEmail(),
			Password:  userDTO.GetPassword(),
			FirstName: userDTO.GetFirstName(),
			LastName:  userDTO.GetLastName(),
		},
	}
}

func (u *UserService) GetByUsername(ctx context.Context, username string) UserOut {
	userDTO, err := u.storage.GetByUsername(ctx, username)
	if err != nil {
		u.logger.Error("user: GetByUsername err", zap.Error(err))
		return UserOut{
			ErrorCode: errors.UserServiceRetrieveUserErr,
		}
	}

	return UserOut{
		User: &models.User{
			ID:        userDTO.GetID(),
			Username:  userDTO.GetUserName(),
			Phone:     userDTO.GetPhone(),
			Email:     userDTO.GetEmail(),
			Password:  userDTO.GetPassword(),
			FirstName: userDTO.GetFirstName(),
			LastName:  userDTO.GetLastName(),
		},
	}
}

func (u *UserService) UpdateUser(ctx context.Context, userdata UpdateUserRequest) UpdateUserResponse {
	hashPass, err := cryptography.HashPassword(userdata.Password)
	if err != nil {
		return UpdateUserResponse{
			Success:   false,
			ErrorCode: errors.HashPasswordError,
		}
	}

	var dto models.UserDTO
	dto.SetUserName(userdata.UserName).
		SetPhone(userdata.Phone).
		SetPassword(hashPass).
		SetFirstName(userdata.FirstName).
		SetLastName(userdata.LastName)

	err = u.storage.UpdateUser(ctx, dto)

	if err != nil {
		u.logger.Error("user: UpdateUser err", zap.Error(err))
		return UpdateUserResponse{
			Success:   false,
			ErrorCode: errors.UserServiceUpdateErr,
		}
	}

	return UpdateUserResponse{
		Success: true,
	}
}
