package interceptors

import (
	"context"
	"log/slog"
	"path"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		method := path.Base(info.FullMethod)
		slog.Info("Был вызван gRPC метод %s ", method)
		startTime := time.Now()
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		if err != nil {
			st, _ := status.FromError(err)
			slog.Error("gRPC метод завершился с ошибкой %s %v", method, st.Code())
		} else {
			slog.Info("gRPC метод успешно завершил работу за %v секунд %s ", duration, method)
		}
		return resp, err
	}
}
