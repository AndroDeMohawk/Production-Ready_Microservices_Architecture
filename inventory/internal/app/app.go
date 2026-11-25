package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/config"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/closer"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/grpc/health"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/logger"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second
)

type App struct {
	diContainer *diContainer
	listener    net.Listener
	grpc        *grpc.Server
}

func New(ctx context.Context) *App {
	a := &App{}
	_, err := config.Load(config.FetchConfigPath())
	if err != nil {
		slog.Error("ошибка загрузки конфига", "error", err)
		os.Exit(1)
	}
	a.initDeps(ctx)
	slog.Debug("meow")
	return a
}

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runGRPCServer(ctx)
	}()
	gwContext, gwCancel := context.WithCancel(context.Background())
	defer gwCancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := inventoryv1.RegisterInventoryServiceHandlerFromEndpoint(
		gwContext,
		mux,
		config.AppConfig().GRPCConfig.Address(),
		opts,
	)
	if err != nil {
		slog.Error("ошибка регистрации inventory gRPC gateway", "error", err)
		os.Exit(1)
	}
	httpServer := &http.Server{
		Addr:    config.AppConfig().GRPCConfig.GwAddr,
		Handler: mux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка запуска inventory HTTP-Gateway", "error", err)
			os.Exit(1)
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			slog.Error("ошибка gRPC-сервера", "error", err)
		}
		cancel()
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel()

	shutdownCtx, stop := context.WithTimeout(context.Background(), shutdownTimeout)
	defer stop()
	closer.Add("HTTP server", func(ctx context.Context) error {
		return httpServer.Shutdown(ctx)
	})
	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка завершения inventory", "error", err)
		os.Exit(1)
	}
}

func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initGRPCServer,
	}
	for _, f := range inits {
		f(ctx)
	}
}

func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}
func (a *App) initGRPCServer(ctx context.Context) {
	a.grpc = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: false,
		}),
	)

	api := a.diContainer.Api(ctx)
	reflection.Register(a.grpc)
	health.RegisterService(a.grpc)
	closer.Add("gRPC Server", func(_ context.Context) error {
		a.grpc.GracefulStop()
		return nil
	})

	inventoryv1.RegisterInventoryServiceServer(a.grpc, api)
}

func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", config.AppConfig().GRPCConfig.Address())
	if err != nil {
		slog.Error("не удалось создать TCP листенера", "error", err)
		os.Exit(1)
	}
	a.listener = listener
}

func (a *App) initLogger(ctx context.Context) {
	logger.Init(config.AppConfig().LoggerConfig.Level())
}

func (a *App) runGRPCServer(_ context.Context) error {
	slog.Debug("gRPC-сервер запущен", "address", config.AppConfig().GRPCConfig.Address())
	return a.grpc.Serve(a.listener)
}
