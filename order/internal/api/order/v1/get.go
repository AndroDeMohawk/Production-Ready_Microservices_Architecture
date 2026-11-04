package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/api/converter"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, req orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.service.Get(ctx, req.OrderUUID)
	if err != nil {
		return nil, err
	}
	dto := converter.OrderToDTO(order)
	return &dto, nil
}
