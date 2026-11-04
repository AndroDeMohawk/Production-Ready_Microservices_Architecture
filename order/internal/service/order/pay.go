package order

import (
	"context"
	"fmt"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (s *orderService) Pay(ctx context.Context, orderID uuid.UUID, method model.PaymentMethod) (*model.Order, error) {
	const op = "order.service.pay"

	order, err := s.repository.Get(ctx, orderID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return nil, errs.ErrOrderNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	switch order.Status {
	case model.OrderStatusPaid:
		return nil, errs.ErrOrderAlreadyPaid
	case model.OrderStatusCancelled:
		return nil, errs.ErrOrderCancelled
	}

	txUUID, err := s.paymentClient.Pay(ctx, orderID, method)
	if err != nil {
		return nil, fmt.Errorf("%s: оплата: %w", op, err)
	}

	if err := s.repository.Pay(ctx, orderID, txUUID, method); err != nil {
		return nil, fmt.Errorf("%s: сохранение: %w", op, err)
	}

	stored, err := s.repository.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: чтение после оплаты: %w", op, err)
	}

	return stored, nil
}
