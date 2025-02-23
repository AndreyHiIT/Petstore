package controller

import (
	"net/http"
	"pet-store/internal/infrastructure/component"
	"pet-store/internal/infrastructure/errors"
	"pet-store/internal/infrastructure/responder"
	"pet-store/internal/modules/order/service"

	"github.com/go-chi/chi"
	"github.com/ptflp/godecoder"
)

type Orderer interface {
	CreateOrder(w http.ResponseWriter, r *http.Request)
	FindOrderByID(w http.ResponseWriter, r *http.Request)
	DeleteOrderByID(w http.ResponseWriter, r *http.Request)
}

type Order struct {
	service service.Orderer
	responder.Responder
	godecoder.Decoder
}

func NewOrder(service service.Orderer, components *component.Components) Orderer {
	return &Order{service: service, Responder: components.Responder, Decoder: components.Decoder}
}

// @Summary Delete Order
// @Tags store
// @Description delete order by ID
// @ID Delete Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "ID of the order to delete"
// @Success 200 {object} OrderDeleteByIDResponse "Seccess"
// @Failure 400 {object} OrderResponseErr "Error"
// @Router /store/order/{orderId} [delete]
func (o *Order) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "orderId")

	out := o.service.DeleteOrderByID(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		msg := "error delete order by ID"
		if out.ErrorCode == errors.OrderServiceDeleteByIDNotFoundID {
			msg = "id not found"
		}
		o.OutputJSON(w, OrderResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	o.OutputJSON(w, OrderDeleteByIDResponse{
		Success: true,
	})
}

// @Summary Find Order
// @Tags store
// @Description find order by ID
// @ID Find Order
// @Accept  json
// @Produce  json
// @Param orderId path string true "ID of the order to delete"
// @Success 200 {object} OrderFindByIDResponse "Successfully retrieved order"
// @Failure 400 {object} OrderResponseErr "Error"
// @Router /store/order/{orderId} [get]
func (o *Order) FindOrderByID(w http.ResponseWriter, r *http.Request) {
	req := chi.URLParam(r, "orderId")

	out := o.service.FindOrderByID(r.Context(), req)
	if out.ErrorCode != errors.NoError {
		msg := "error find order by ID"
		if out.ErrorCode == errors.OrderServiceFindByIDNotFoundID {
			msg = "order id not found"
		}
		o.OutputJSON(w, OrderResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	o.OutputJSON(w, OrderFindByIDResponse{
		Success: true,
		Order:   out.Order,
	})
}

// @Summary Create Order
// @Tags store
// @Description create order 
// @ID Create Order
// @Accept  json
// @Produce  json
// @Param input body service.RequestCreateOrder true "Create order"
// @Success 200 {object} SuccessCreateOrderResponse "Success"
// @Failure 400 {object} OrderResponseErr "Error"
// @Router /store/order [post]
func (o *Order) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req service.RequestCreateOrder

	err := o.Decode(r.Body, &req)
	if err != nil {
		o.ErrorBadRequest(w, err)
		return
	}

	out := o.service.CreateOrder(r.Context(), req)

	if out.ErrorCode != errors.NoError {
		msg := "error create order"
		o.OutputJSON(w, OrderResponseErr{
			Success:   false,
			ErrorCode: out.ErrorCode,
			Data: Data{
				Message: msg,
			},
		})
		return
	}

	o.OutputJSON(w, SuccessCreateOrderResponse{
		Success: true,
		OrderID: out.ID,
	})
}
