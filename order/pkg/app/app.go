package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	apiv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/api/order/v1"
	inventoryClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/payment/v1"
	orderRepo "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/order"
	orderService "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/order"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	manager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	ogenmiddleware "github.com/ogen-go/ogen/middleware"
)

type App struct {
	httpHandler http.Handler
}

func NewApp(
	inventoryGRPCClient inventoryv1.InventoryServiceClient,
	paymentGRPCClient paymentv1.PaymentServiceClient,

) (*App, error) {
	ctx := context.Background()
	_ = godotenv.Load("../order.env")
	DSN := os.Getenv("DB_URI")
	pool, err := pgxpool.New(ctx, DSN)
	if err != nil {
		slog.Error("не удалось создать пул подключений к БД", "error", err)
		os.Exit(1)
	}
	err = pool.Ping(ctx)
	if err != nil {
		slog.Error("проверка соединения с БД", "error", err)
		os.Exit(1)
	}

	slog.Info("подключение к PostgreSQL установлено")
	manager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		slog.Error("создание transaction manager", "error", err)
		os.Exit(1)
	}
	Repository := orderRepo.New(pool)
	invClient := inventoryClient.NewClientAdapter(inventoryGRPCClient)
	paymClient := paymentClient.NewClientAdapter(paymentGRPCClient)

	Service := orderService.NewService(Repository, invClient, paymClient, manager)

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
