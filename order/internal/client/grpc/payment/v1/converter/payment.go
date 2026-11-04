package converter

import (
	"fmt"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(method model.PaymentMethod) (paymentv1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, fmt.Errorf("неизвестный метод оплаты %s", method)
	}
}

func PaymentMethodFromProto(method paymentv1.PaymentMethod) (model.PaymentMethod, error) {
	switch method {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard, nil
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney, nil
	default:
		return model.PaymentMethodUnspecified, fmt.Errorf("неизвестный метод оплаты(proto) %s", method)
	}
}
