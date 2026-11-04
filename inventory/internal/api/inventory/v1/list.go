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

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	if len(req.GetUuids()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "список id не может быть пустым")

	}
	input := converter.ListPartsRequestToInput(req)
	parts, err := a.service.List(ctx, input)
	protoParts := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, converter.ModelToProto(part))
	}
	if err != nil {
		if errors.Is(err, errs.ErrPartsNotFound) {
			return &inventoryv1.ListPartsResponse{
				Parts: protoParts,
			}, status.Errorf(codes.OK, "Предупреждение.Одна или несколько деталей не были получены.")

		}
		return nil, status.Errorf(codes.Internal, "внутренняя ошибка сервера.")

	}
	return &inventoryv1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
