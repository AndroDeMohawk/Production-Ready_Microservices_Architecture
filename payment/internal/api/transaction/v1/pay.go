package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/api/converter"
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/errors"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req.OrderUuid == "" {
		return nil, status.Errorf(codes.InvalidArgument, "пустой id заказа")
	}
	if req.PaymentMethod == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "пустой метод оплаты заказа")
	}
	input, err := converter.PayOrderRequestToInput(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка парсинга id")
	}
	tx, err := a.service.Pay(ctx, input)
	if err != nil {
		if errors.Is(err, errs.ErrTransactionNotFound) {
			return nil, status.Errorf(codes.NotFound, "транзакция не найдена")

		}
		if errors.Is(err, errs.ErrDuplicatePayOrder) {
			return nil, status.Errorf(codes.AlreadyExists, "дублирование оплаты заказа")

		}
		if errors.Is(err, errs.ErrEmptyOrderID) {
			return nil, status.Errorf(codes.InvalidArgument, "пустой id заказа")
		}
		return nil, status.Errorf(codes.Internal, "внутренняя ошибка сервера")
	}
	return &paymentv1.PayOrderResponse{
		TransactionUuid: tx.TransactionUUID.String(),
	}, nil
}
