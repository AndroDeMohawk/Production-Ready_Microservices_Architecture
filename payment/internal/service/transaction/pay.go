package transaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/service/input"
	"github.com/google/uuid"
)

func (s *service) Pay(ctx context.Context, in *input.Pay) (*model.Transaction, error) {
	const op = "payment.service.pay"
	if in.OrderId == uuid.Nil {
		return nil, fmt.Errorf("%s %w", op, errs.ErrEmptyOrderID)
	}

	tx, err := s.repo.GetByOrder(ctx, in.OrderId)
	switch {
	case err == nil:
		return tx, nil
	case !errors.Is(err, errs.ErrTransactionNotFound):
		return nil, fmt.Errorf("%s: чтение: %w", op, err)
	}
	now := time.Now()
	tx = &model.Transaction{
		TransactionUUID: uuid.New(),
		OrderUUID:       in.OrderId,
		Method:          in.Method,
		Status:          model.TransactionStatusSucceeded,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		if errors.Is(err, errs.ErrDuplicatePayOrder) {
			existing, getErr := s.repo.GetByOrder(ctx, in.OrderId)
			if getErr == nil {
				return existing, nil
			}
		}

		return nil, fmt.Errorf("%s: запись: %w", op, err)
	}

	return tx, nil
}
