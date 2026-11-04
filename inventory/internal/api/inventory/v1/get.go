package v1

import (
	"context"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/api/converter"
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	"github.com/go-faster/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "id запчасти не может быть пустым")
	}
	part, err := a.service.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "деталь не найдена")
		}
		if errors.Is(err, errs.ErrParsingUUID) {
			return nil, status.Errorf(codes.InvalidArgument, "ошибка парсинга uuid")
		}
		return nil, status.Errorf(codes.Internal, "внутренняя ошибка сервера %v", err)
	}
	return &inventoryv1.GetPartResponse{
		Part: converter.ModelToProto(part),
	}, nil
}
