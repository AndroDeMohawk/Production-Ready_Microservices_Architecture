package app

import (
	"context"
	"log/slog"
	"os"

	api "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/api/inventory/v1"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/config"
	repository "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/part"
	service "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/service/part"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/closer"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"

	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	pgPool     *pgxpool.Pool
	repository service.Repository
	service    api.Service
	handler    inventoryv1.InventoryServiceServer
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PGConfig.DSN())

		if err != nil {
			slog.Error("ошибка создания пула соединений к БД ", "error", err)
			os.Exit(1)
		}

		if err := pool.Ping(ctx); err != nil {
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

func (d *diContainer) Repository(ctx context.Context) service.Repository {
	if d.repository == nil {
		d.repository = repository.New(d.PGPool(ctx))
	}
	return d.repository
}

func (d *diContainer) Service(ctx context.Context) api.Service {
	if d.service == nil {
		d.service = service.NewService(d.Repository(ctx))
	}
	return d.service
}
func (d *diContainer) Api(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.handler == nil {
		d.handler = api.NewApi(d.Service(ctx))
	}
	return d.handler
}
