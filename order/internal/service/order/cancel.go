package order

import (
	"context"
	"fmt"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
)

func (s *orderService) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	const op = "order.service.cancel"
	order, err := s.repository.Get(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	if order.Status == model.OrderStatusPaid {
		return errs.ErrOrderAlreadyPaid
	}
	if order.Status == model.OrderStatusCancelled {
		return errs.ErrOrderCancelled
	}
	err = s.repository.Cancel(ctx, orderUUID)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}
	return nil

}
