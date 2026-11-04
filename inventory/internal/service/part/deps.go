package part

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/google/uuid"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (*model.Part, error)
	List(ctx context.Context, uuids []uuid.UUID) ([]*model.Part, error)
}
