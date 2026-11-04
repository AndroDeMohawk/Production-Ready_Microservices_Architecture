package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/service/input"
)

type Service interface {
	Get(ctx context.Context, id string) (*model.Part, error)
	List(ctx context.Context, in input.ListPartsInput) ([]*model.Part, error)
}
