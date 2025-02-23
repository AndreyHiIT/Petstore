package storage

import (
	"context"
	"pet-store/internal/models"
)

type Orderer interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	FindOrderByID(ctx context.Context, orderID int) (models.Order, error)
	DeleteOrderByID(ctx context.Context, orderID int)  error
}
