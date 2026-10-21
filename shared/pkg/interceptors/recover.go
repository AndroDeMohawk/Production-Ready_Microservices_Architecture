package interceptors

import (
	"context"
	"log/slog"
	"path"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Recover() grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		defer func() {
			if err := recover(); err != nil {
				method := path.Base(info.FullMethod)
				slog.Error("паника в методе",
					"method", method,
					"panic", err,
					"stack", string(debug.Stack()),
				)
				err = status.Errorf(codes.Internal, "внутренняя ошибка сервера")

			}

		}()
		return handler(ctx, req)
	}
}
