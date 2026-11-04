package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/service/input"
)

type Service interface {
	Pay(ctx context.Context, in *input.Pay) (*model.Transaction, error)
}
