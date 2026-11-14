package order

import (
	"context"
	"errors"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Pay(ctx context.Context, orderID uuid.UUID, txUUID uuid.UUID) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const queryGet = `SELECT status FROM orders WHERE uuid = $1`
	var status model.OrderStatus
	err := querier.QueryRow(ctx, queryGet, orderID).Scan(
		&status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.ErrOrderNotFound
		}
		return err
	}
	if status == model.OrderStatusPaid {
		return errs.ErrOrderAlreadyPaid
	}
	if status == model.OrderStatusCancelled {
		return errs.ErrOrderCancelled
	}
	const queryUpdate = `UPDATE orders SET status = $3,transaction_uuid = $1,updated_at = $5 WHERE uuid = $2 AND status not in ($3,$4)`
	tag, err := querier.Exec(ctx, queryUpdate, txUUID, orderID, model.OrderStatusPaid, model.OrderStatusCancelled, time.Now())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {

		return errs.ErrStatusOrderConflict
	}

	return nil
}
