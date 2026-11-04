package v1

import (
	"context"
	"log/slog"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/api/converter"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	slog.Info("api.create конец")
	input := converter.CreateOrderRequestToInput(req)
	order, err := a.service.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	slog.Info("api.create конец")
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil

}
