package controller

import "pet-store/internal/models"

type OrderResponseErr struct {
	Success   bool
	ErrorCode int
	Data      Data
}

type Data struct {
	Message string
}

type SuccessCreateOrderResponse struct {
	Success bool
	OrderID int
}

type OrderFindByIDResponse struct {
	Success bool
	Order   models.Order
}

type OrderDeleteByIDResponse struct {
	Success bool
}
