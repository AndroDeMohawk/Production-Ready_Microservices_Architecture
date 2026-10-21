package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	inventoryService "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/pkg/service"
	interceptor "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/interceptors"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
)

const (
	grpcAddress               = ":50051"
	grpcGatewayAddress        = ":40041"
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
)

func main() {
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("не удалось создать listener", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
		grpc.ChainUnaryInterceptor(
			interceptor.Logger(),
			interceptor.Recover(),
		),
	)
	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryService.NewServer())
	reflection.Register(grpcServer)

	slog.Info("запуск InventoryService", "адрес", grpcAddress)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("ошибка запуска сервера", fmt.Errorf(err.Error()))
			os.Exit(1)
		}
	}()
	gwContext, gwCancel := context.WithCancel(context.Background())
	defer gwCancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = inventoryv1.RegisterInventoryServiceHandlerFromEndpoint(
		gwContext,
		mux,
		grpcAddress,
		opts,
	)
	if err != nil {
		slog.Error("ошибка регистрации inventory gRPC gateway", "error", err)
		gwCancel()
		return
	}
	httpServer := &http.Server{
		Addr:    grpcGatewayAddress,
		Handler: mux,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			slog.Error("ошибка запуска inventory HTTP-Gateway", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("Остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("Сервер успешно остановлен")

	shutdownCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	if error := httpServer.Shutdown(shutdownCtx); error != nil {
		slog.Error("ошибка остановки inventory HTTP-Gateway", "error", error)
	}
	slog.Info("inventory HTTP-Gateway остановлен")
}
