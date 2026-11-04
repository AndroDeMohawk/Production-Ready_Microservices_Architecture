package transaction

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) GetByOrder(ctx context.Context, orderId uuid.UUID) (*model.Transaction, error) {
	r.mu.RLock()
	tx, exist := r.txStore[orderId]
	r.mu.RUnlock()
	if !exist {
		return nil, errs.ErrTransactionNotFound
	}
	modelTx := converter.RecordToModel(&tx)
	return modelTx, nil
}
