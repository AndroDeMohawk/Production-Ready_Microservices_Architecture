package converter

import (
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/service/input"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func PayOrderRequestToInput(r *paymentv1.PayOrderRequest) (*input.Pay, error) {
	id, err := uuid.Parse(r.OrderUuid)
	if err != nil {
		return nil, errs.ErrParsingUUID
	}
	return &input.Pay{
		OrderId: id,
	}, nil
}
