package converter

import (
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
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
		Method:  ProtoPaymentMethodTModel(r.PaymentMethod),
	}, nil
}

func ProtoPaymentMethodTModel(method paymentv1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
