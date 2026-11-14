package transaction

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/converter"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/record"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (r *repository) GetByOrder(ctx context.Context, orderId uuid.UUID) (*model.Transaction, error) {
	querirer := r.getter.DefaultTrOrDB(ctx, r.pool)
	const query = `SELECT uuid,order_uuid,status,created_at,updated_at FROM transactions where order_uuid = $1`
	var tx record.Transaction
	err := querirer.QueryRow(ctx, query, orderId).Scan(
		&tx.TransactionUUID,
		&tx.OrderUUID,
		&tx.Status,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, errs.ErrTransactionNotFoundInDB) {
			return nil, errs.ErrTransactionNotFound
		}
		return nil, err
	}
	modelTx := converter.RecordToModel(&tx)
	return modelTx, nil
}
