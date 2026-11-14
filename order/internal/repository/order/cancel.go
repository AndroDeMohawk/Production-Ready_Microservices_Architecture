package order

import (
	"context"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Cancel(ctx context.Context, orderID uuid.UUID) error {
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const query = `UPDATE orders SET status = $1,updated_at = $2 WHERE uuid = $3 AND status != $1`
	tag, err := querier.Exec(ctx, query, model.OrderStatusCancelled, time.Now(), orderID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil

}
