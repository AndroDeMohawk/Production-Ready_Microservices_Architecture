package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderHandler "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/pkg/handler"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
)

const (
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	inventoryConn, err := grpc.NewClient(inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		os.Exit(1)
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		os.Exit(1)
	}
	defer paymentConn.Close()
	// Создаём хранилище и обработчик
	store := orderHandler.NewOrderStore()
	h := orderHandler.NewHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	orderServer, err := orderHandler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	httpServer := http.Server{
		Addr:              ":8080",
		Handler:           orderServer,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      60 * time.Second,
		ReadTimeout:       60 * time.Second,
	}

	go func() {
		slog.Info("запуск OrderService", "port", 8080)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка запуска сервера", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("получен сигнал завершения, начинаем graceful shutdown...")
	shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown is bad", err.Error())
		os.Exit(1)
	}
	slog.Info("OrderService успешно остановлен")
}
