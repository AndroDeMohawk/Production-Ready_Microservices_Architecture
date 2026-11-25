package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/config"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/closer"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/platform/pkg/logger"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
}

func New(ctx context.Context) *App {
	a := &App{}
	_, err := config.Load(config.FetchConfigPath())
	if err != nil {
		slog.Error("ошибка загрузки конфига", "error", err)
		os.Exit(1)
	}
	a.initDeps(ctx)
	return a
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	slog.Info("HTTP-сервер запускается", "address", config.AppConfig().OrderHTTPConfig.Address())
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка запуска http сервера", "error", err) //TODO: потом заменить на свой библиотечный  логер
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("получен сигнал завершения, начинаем graceful shutdown...")
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	closer.Add("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
	if err := closer.CloseAll(shutDownCtx); err != nil {
		slog.Info("ошибка запуска плавного завершения сервера", "error", err)
		os.Exit(1)
	}
	slog.Info("сервер успешно остановлен")
}

func (a *App) intiDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

func (a *App) initLogger(ctx context.Context) {
	logger.Init(config.AppConfig().LoggerConfig.Level())
}
func (a *App) initHTTPServer(ctx context.Context) {
	handler := a.diContainer.OrderHandler(ctx)
	httpServer := http.Server{
		Addr:              config.AppConfig().OrderHTTPConfig.Address(),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadTimeout:       60 * time.Second,
	}
	a.httpServer = &httpServer

}
func (a *App) initDeps(ctx context.Context) {
	inits := []func(ctx context.Context){
		a.intiDI,
		a.initLogger,
		a.initHTTPServer,
	}
	for _, f := range inits {
		f(ctx)
	}
}
