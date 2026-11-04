package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/input"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, in input.CreateOrder) (*model.Order, error)
	Pay(ctx context.Context, orderID uuid.UUID, method model.PaymentMethod) (*model.Order, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
	Get(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error)
}
