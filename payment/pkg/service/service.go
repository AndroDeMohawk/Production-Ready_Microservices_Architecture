package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
)

// server реализует gRPC сервис оплаты
type server struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// NewServer создаёт новый экземпляр сервера оплаты
func NewServer() *server {
	return &server{}
}

// PayOrder обрабатывает оплату заказа
func (s *server) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	id := req.GetOrderUuid()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "order uuid required")
	}
	if req.GetPaymentMethod().String() == "UNSPECIFIED" {
		return nil, status.Error(codes.InvalidArgument, "unspecified payment method")
	}
	idStr, err := uuid.Parse(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}
	transactionId := uuid.New()
	slog.Info("Оплата прошла успешно",
		"order_uuid: ", idStr,
		"transaction_id: ", transactionId,
	)
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionId.String(),
	}, nil

}
