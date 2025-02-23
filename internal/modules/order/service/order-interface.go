package service

import (
	"context"
	"pet-store/internal/models"
	"time"
)

type Orderer interface {
	CreateOrder(ctx context.Context, order RequestCreateOrder) ResponseCreateOrder
	FindOrderByID(ctx context.Context, orderIDstr string) ResponseFindOrderByID
	DeleteOrderByID(ctx context.Context, orderIDstr string) ResponseDeleteOrderByID
}

type RequestCreateOrder struct {
	PetId    int       `json:"petid"`
	Quantity int       `json:"quantity"`
	ShipDate time.Time `json:"shipdate"`
	Status   string    `json:"status"`
	Complete bool      `json:"complete"`
}

type ResponseCreateOrder struct {
	Status    bool
	ErrorCode int
	ID        int
}

type ResponseFindOrderByID struct {
	Status    bool
	ErrorCode int
	Order     models.Order
}

type ResponseDeleteOrderByID struct {
	Status    bool
	ErrorCode int
}