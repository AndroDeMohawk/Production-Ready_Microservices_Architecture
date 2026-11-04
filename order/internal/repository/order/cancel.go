package order

import (
	"context"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
	"github.com/google/uuid"
)

func (r *orderStore) Cancel(ctx context.Context, orderID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	order, exist := r.orders[orderID]
	if !exist {
		return errs.ErrOrderNotFound
	}
	now := time.Now()
	order.CanceledAt = &now
	order.Status = record.OrderStatus(model.OrderStatusCancelled)
	r.orders[orderID] = order

	return nil

}
