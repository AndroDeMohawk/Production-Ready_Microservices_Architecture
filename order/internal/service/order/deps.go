package order

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Order, error)
	Create(ctx context.Context, in model.CreateOrder) (*model.Order, error)
	Pay(ctx context.Context, orderID uuid.UUID, txUUID uuid.UUID) error
	Cancel(ctx context.Context, orderID uuid.UUID) error
}
type InventoryClient interface {
	GetPart(ctx context.Context, partID string) (*model.Part, error)
	ListParts(ctx context.Context, partType model.PartType, uuids []string) ([]*model.Part, error)
}
type PaymentClient interface {
	Pay(ctx context.Context, orderID uuid.UUID) (uuid.UUID, error)
}
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
