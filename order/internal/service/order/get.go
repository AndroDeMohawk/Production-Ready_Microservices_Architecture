package order

import (
	"context"
	"errors"
	"fmt"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
)

func (s *orderService) Get(ctx context.Context, orderUUID uuid.UUID) (*model.Order, error) {
	const op = "order.service.get"
	order, err := s.repository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return order, nil
}
