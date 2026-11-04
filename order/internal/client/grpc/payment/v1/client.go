package v1

import (
	"context"
	"fmt"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/payment/v1/converter"
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientAdapter struct {
	grpcClient paymentv1.PaymentServiceClient
}

func NewClientAdapter(grpcClient paymentv1.PaymentServiceClient) *clientAdapter {
	return &clientAdapter{
		grpcClient: grpcClient,
	}
}

func (a *clientAdapter) Pay(ctx context.Context, orderID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	const op = "order.client.payment.pay"

	protoPaymentMethod, err := converter.PaymentMethodToProto(method)
	if err != nil {
		return uuid.Nil, errs.ErrInvalidPaymentMethod
	}

	resp, err := a.grpcClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     orderID.String(),
		PaymentMethod: protoPaymentMethod,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidPaymentMethod)
			case codes.AlreadyExists:
				return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrOrderAlreadyPaid)
			}
		}

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	txUUID, err := uuid.Parse(resp.GetTransactionUuid())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: разбор transaction_uuid: %w", op, err)
	}

	return txUUID, nil
}
