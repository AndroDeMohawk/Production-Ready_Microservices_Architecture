package transaction

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, tx *model.Transaction) error {
	recordTx := converter.ModelToRecord(tx)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.txStore[tx.OrderUUID]; exist {
		return errs.ErrDuplicatePayOrder
	}
	r.txStore[tx.OrderUUID] = *recordTx
	return nil
}
