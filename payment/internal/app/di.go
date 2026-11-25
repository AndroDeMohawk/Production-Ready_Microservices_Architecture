package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	api "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/api/transaction/v1"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/config"
	repository "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/transaction"
	service "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/service/transaction"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/closer"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	pgPool     *pgxpool.Pool
	repository service.Repository
	service    api.Service
	handler    paymentv1.PaymentServiceServer
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PGConfig.DSN())
		if err != nil {
			slog.Error("ошибка инициализации пула соединений к БД", "error", err)
			os.Exit(1)
		}
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := pool.Ping(pingCtx); err != nil {
			pool.Close()
			slog.Error("БД не пингует", "error", err)
			os.Exit(1)
		}
		closer.Add("PGpool", func(_ context.Context) error {
			pool.Close()
			return nil
		})
		d.pgPool = pool

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
		d.service = service.New(d.Repository(ctx))
	}
	return d.service
}

func (d *diContainer) Api(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.handler == nil {
		d.handler = api.NewApi(d.Service(ctx))
	}
	return d.handler
}
