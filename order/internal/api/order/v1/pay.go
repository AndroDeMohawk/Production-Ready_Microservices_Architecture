package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {

	orderUUID := params.OrderUUID
	paymentMethod := model.PaymentMethod(req.PaymentMethod)
	order, err := a.service.Pay(ctx, orderUUID, paymentMethod)
	if err != nil {
		return nil, err
	}
	return &orderv1.PayOrderResponse{
		TransactionUUID: *order.TransactionUUID,
	}, nil
}
