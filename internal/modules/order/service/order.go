package service

import (
	"context"
	"database/sql"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/models"
	"pet-store/internal/modules/order/storage"
	"strconv"

	"go.uber.org/zap"
)

type OrderService struct {
	storage storage.Orderer
	logger  *zap.Logger
}

func NewOrderService(storage storage.Orderer, logger *zap.Logger) *OrderService {
	return &OrderService{storage: storage, logger: logger}
}

func (o *OrderService) DeleteOrderByID(ctx context.Context, orderIDstr string) ResponseDeleteOrderByID {
	orderID, err := strconv.Atoi(orderIDstr)
	if err != nil {
		o.logger.Error("Error during conversion:", zap.Error(err))
		return ResponseDeleteOrderByID{
			Status:    false,
			ErrorCode: errors.DeleteOrderByIDErrorDuringConversion,
		}
	}
	if orderID <= 0 {
		o.logger.Error("Error during conversion: orderID <= 0")
		return ResponseDeleteOrderByID{
			Status:    false,
			ErrorCode: errors.DeleteOrderByIDErrorIDLessZero,
		}
	}

	err = o.storage.DeleteOrderByID(ctx, orderID)
	if err != nil {
		o.logger.Error("Delete Order By ID:", zap.Error(err))
		if err.Error() == "no order found with ID" {
			return ResponseDeleteOrderByID{
				Status:    false,
				ErrorCode: errors.OrderServiceDeleteByIDNotFoundID,
			}
		}
		return ResponseDeleteOrderByID{
			Status:    false,
			ErrorCode: errors.OrderServiceDeleteByIDIternalErr,
		}
	}

	return ResponseDeleteOrderByID{
		Status: true,
	}
}

func (o *OrderService) FindOrderByID(ctx context.Context, orderIDstr string) ResponseFindOrderByID {
	orderID, err := strconv.Atoi(orderIDstr)
	if err != nil {
		o.logger.Error("Error during conversion:", zap.Error(err))
		return ResponseFindOrderByID{
			Status:    false,
			ErrorCode: errors.FindOrderByIDErrorDuringConversion,
		}
	}
	if orderID <= 0 {
		o.logger.Error("Error during conversion: orderID <= 0")
		return ResponseFindOrderByID{
			Status:    false,
			ErrorCode: errors.FindOrderByIDErrorIDLessZero,
		}
	}

	order, err := o.storage.FindOrderByID(ctx, orderID)
	if err != nil {
		o.logger.Error("Find Order By ID:", zap.Error(err))
		if err == sql.ErrNoRows {
			return ResponseFindOrderByID{
				Status:    false,
				ErrorCode: errors.OrderServiceFindByIDNotFoundID,
			}
		}
		return ResponseFindOrderByID{
			Status:    false,
			ErrorCode: errors.OrderServiceFindByIDIternalErr,
		}
	}

	return ResponseFindOrderByID{
		Status: true,
		Order:  order,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, order RequestCreateOrder) ResponseCreateOrder {
	orderDto := models.Order{
		PetID:    order.PetId,
		Quantity: order.Quantity,
		ShipDate: order.ShipDate,
		Status:   order.Status,
		Complete: order.Complete,
	}
	ordID, err := o.storage.CreateOrder(ctx, orderDto)
	if err != nil {
		o.logger.Error("Error Create Order:", zap.Error(err))
		return ResponseCreateOrder{
			Status:    false,
			ErrorCode: errors.OrderServiceCreateOrderErr,
		}
	}

	return ResponseCreateOrder{
		Status: true,
		ID:     ordID,
	}
}
