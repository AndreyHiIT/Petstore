package storage

import (
	"context"
	"pet-store/internal/db/adapter"
	"pet-store/internal/models"
)

// UserStorage - хранилище пользователей
type UserStorage struct {
	adapter *adapter.SQLAdapter
}

// NewUserStorage - конструктор хранилища пользователей
func NewUserStorage(sqlAdapter *adapter.SQLAdapter) *UserStorage {
	return &UserStorage{adapter: sqlAdapter}
}

// Create - создание пользователя в БД
func (s *UserStorage) Create(ctx context.Context, u models.UserDTO) (int, error) {
	id, err := s.adapter.Create(ctx, &u)

	return id, err
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (models.UserDTO, error) {
	var user models.UserDTO
	err := s.adapter.GetUserByEmail(ctx, &user, email)
	if err != nil {
		return models.UserDTO{}, err
	}
	
	return user, nil
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (models.UserDTO, error) {
	var user models.UserDTO
	err := s.adapter.GetUserByUsername(ctx, &user, username)
	if err != nil {
		return models.UserDTO{}, err
	}
	
	return user, nil
}

func (s *UserStorage) UpdateUser(ctx context.Context,  userdata models.UserDTO) error {
	err:= s.adapter.UpdateUser(ctx, &userdata)
	if err != nil{
		return  err
	}
	return nil
}