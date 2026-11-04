package v1

import (
	"context"

	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, req orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	err := a.service.Cancel(ctx, req.OrderUUID)
	if err != nil {
		return nil, err
	}
	return &orderv1.CancelOrderResponse{}, nil
}
