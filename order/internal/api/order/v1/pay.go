package v1

import (
	"context"

	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {

	orderUUID := params.OrderUUID
	order, err := a.service.Pay(ctx, orderUUID)
	if err != nil {
		return nil, err
	}
	return &orderv1.PayOrderResponse{
		TransactionUUID: *order.TransactionUUID,
	}, nil
}
