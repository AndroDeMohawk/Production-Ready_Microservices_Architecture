package transaction

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	GetByOrder(ctx context.Context, orderId uuid.UUID) (*model.Transaction, error)
}
