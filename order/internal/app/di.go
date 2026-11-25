package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	api "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/api/order/v1"
	inventoryClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/payment/v1"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/config"
	repo "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/order"
	service "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/order"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/closer"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	pgPool              *pgxpool.Pool
	inventoryConn       *grpc.ClientConn
	paymentConn         *grpc.ClientConn
	orderRepo           service.Repository
	orderService        api.Service
	orderHandler        *orderv1.Server
	invClient           service.InventoryClient
	payClient           service.PaymentClient
	txManager           service.TxManager
	paymentGRPCClient   paymentv1.PaymentServiceClient
	inventoryGRPCClient inventoryv1.InventoryServiceClient
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {

		pool, err := pgxpool.New(ctx, config.AppConfig().PGConfig.DSN())
		if err != nil {
			slog.Error("не удалось подключиться к PostgreSQL", "error", err)
			os.Exit(1)
		}
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			pool.Close()
			slog.Error("БД не пингует", "error", err)
			os.Exit(1)
		}

		d.pgPool = pool
		closer.Add("PGpool", func(_ context.Context) error {
			pool.Close()
			return nil
		})
	}
	return d.pgPool

}
func (d *diContainer) TxManager(ctx context.Context) service.TxManager {
	if d.txManager == nil {
		d.txManager, _ = manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		return d.txManager
	}
	return d.txManager
}

func (d *diContainer) Repository(ctx context.Context) service.Repository {
	if d.orderRepo == nil {
		d.orderRepo = repo.New(d.PGPool(ctx))
	}
	return d.orderRepo
}
func (d *diContainer) PaymentClient(ctx context.Context) service.PaymentClient {
	if d.payClient == nil {
		d.payClient = paymentClient.NewClientAdapter(d.PaymentGRPCClient(ctx))

	}
	return d.payClient
}

func (d *diContainer) InventoryClient(ctx context.Context) service.InventoryClient {
	if d.invClient == nil {
		d.invClient = inventoryClient.NewClientAdapter(d.InventoryGRPCClient(ctx))

	}
	return d.invClient
}

func (d *diContainer) Service(ctx context.Context) api.Service {
	if d.orderService == nil {
		d.orderService = service.NewService(d.Repository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx), d.TxManager(ctx))

	}
	return d.orderService
}

func (d *diContainer) OrderHandler(ctx context.Context) *orderv1.Server {
	if d.orderHandler == nil {
		apiHandler := api.NewApi(d.Service(ctx))
		handler, err := orderv1.NewServer(apiHandler)
		if err != nil {
			panic(err)
		}
		d.orderHandler = handler
	}
	return d.orderHandler
}

func (d *diContainer) InventoryGRPCClient(_ context.Context) inventoryv1.InventoryServiceClient {
	if d.inventoryGRPCClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryGRPCConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось создать gRPC-клиент Inventory", "error", err)
			os.Exit(1)
		}

		d.inventoryConn = conn
		closer.Add("Inventory gRPC client", func(_ context.Context) error {
			return d.inventoryConn.Close()
		})
		d.inventoryGRPCClient = inventoryv1.NewInventoryServiceClient(conn)
	}

	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentGRPCClient(_ context.Context) paymentv1.PaymentServiceClient {
	if d.paymentGRPCClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentGRPCConfig.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			slog.Error("не удалось создать gRPC-клиент Payment", "error", err)
			os.Exit(1)
		}

		d.paymentConn = conn
		closer.Add("Payment gRPC client", func(_ context.Context) error {
			return d.paymentConn.Close()
		})
		d.paymentGRPCClient = paymentv1.NewPaymentServiceClient(conn)
	}

	return d.paymentGRPCClient
}
