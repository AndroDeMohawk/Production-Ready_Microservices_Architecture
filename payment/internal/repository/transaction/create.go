package transaction

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, tx *model.Transaction) error {
	recordTx := converter.ModelToRecord(tx)
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const queryInsert = `INSERT INTO transactions (uuid,order_uuid,status,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (uuid) DO NOTHING`
	tag, err := querier.Exec(ctx, queryInsert, recordTx.TransactionUUID, recordTx.OrderUUID, recordTx.Status, recordTx.CreatedAt, recordTx.UpdatedAt)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrTransactionNotFound
	}
	return nil
}
