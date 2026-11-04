package order

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
	"github.com/google/uuid"
)

func (r *orderStore) Pay(ctx context.Context, orderID uuid.UUID, txUUID uuid.UUID, method model.PaymentMethod) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exist := r.orders[orderID]
	if !exist {
		return errs.ErrOrderNotFound
	}

	switch order.Status {
	case record.OrderStatusPaid:
		return errs.ErrOrderAlreadyPaid
	case record.OrderStatusCancelled:
		return errs.ErrOrderCancelled
	}

	order.Status = record.OrderStatusPaid
	order.TransactionUUID = &txUUID
	paymentMethod := record.PaymentMethod(method)
	order.PaymentMethod = &paymentMethod
	r.orders[orderID] = order

	return nil
}
