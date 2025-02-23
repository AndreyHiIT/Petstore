package storage

import (
	"context"
	"pet-store/internal/db/adapter"
	"pet-store/internal/models"
)

// PetStorage - хранилище животных
type OrderStorage struct {
	adapter *adapter.SQLAdapter
}

// NewUserStorage - конструктор хранилища пользователей
func NewOrderStorage(sqlAdapter *adapter.SQLAdapter) *OrderStorage {
	return &OrderStorage{adapter: sqlAdapter}
}

func (o *OrderStorage) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	return o.adapter.CreateOrder(ctx, order)
}

func (o *OrderStorage) FindOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	return o.adapter.FindOrderByID(ctx, orderID)
}

func (o *OrderStorage) DeleteOrderByID(ctx context.Context, orderID int)  error {
	return o.adapter.DeleteOrderByID(ctx, orderID)
}