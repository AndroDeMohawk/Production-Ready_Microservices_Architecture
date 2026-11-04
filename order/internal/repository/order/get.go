package order

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *orderStore) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	r.mu.RLock()
	order, exist := r.orders[id]
	r.mu.RUnlock()
	if !exist {
		return nil, errs.ErrOrderNotFound
	}
	orderModel := converter.RecordToModel(&order)
	return &orderModel, nil

}
