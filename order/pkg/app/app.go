package app

import (
	"fmt"
	"log/slog"
	"net/http"

	apiv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/api/order/v1"
	inventoryClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/order"
	orderService "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/order"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	ogenmiddleware "github.com/ogen-go/ogen/middleware"
)

type App struct {
	httpHandler http.Handler
}

func NewApp(
	inventoryGRPCClient inventoryv1.InventoryServiceClient,
	paymentGRPCClient paymentv1.PaymentServiceClient,

) (*App, error) {
	repository := orderRepo.NewOrderStore()

	invClient := inventoryClient.NewClientAdapter(inventoryGRPCClient)
	paymClient := paymentClient.NewClientAdapter(paymentGRPCClient)

	Service := orderService.NewService(repository, invClient, paymClient)

	apiHandler := apiv1.NewApi(Service)
	server, err := orderv1.NewServer(
		apiHandler,
		orderv1.WithErrorHandler(apiv1.ErrorHandler),
		orderv1.WithMiddleware(func(req ogenmiddleware.Request, next ogenmiddleware.Next) (ogenmiddleware.Response, error) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("PANIC в обработчике", "panic", r, "op", req.OperationID)
				}
			}()
			return next(req)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания OpenAPI сервера: %w", err)
	}

	return &App{
		httpHandler: server,
	}, nil
}
func (a *App) HTTPHandler() http.Handler {
	return a.httpHandler
}
