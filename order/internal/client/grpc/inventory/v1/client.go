package v1

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/client/grpc/inventory/v1/converter"
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type clientAdapter struct {
	grpcClient inventoryv1.InventoryServiceClient
}

func NewClientAdapter(grpcClient inventoryv1.InventoryServiceClient) *clientAdapter {
	return &clientAdapter{
		grpcClient: grpcClient,
	}
}

func (a *clientAdapter) GetPart(ctx context.Context, partID string) (*model.Part, error) {
	const op = "order.client.inventory.getpart"
	resp, err := a.grpcClient.GetPart(ctx, &inventoryv1.GetPartRequest{
		Uuid: partID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return &model.Part{}, errs.ErrPartNotFound
			}
		}
		return &model.Part{}, fmt.Errorf("%s: %w", op, err)

	}
	return converter.ProtoPartToModel(resp.GetPart()), nil

}
func (a *clientAdapter) ListParts(ctx context.Context, partType model.PartType, uuids []string) ([]*model.Part, error) {
	const op = "order.client.inventory.listparts"
	slog.Info("adapter.inventory нач")
	resp, err := a.grpcClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		PartType: inventoryv1.PartType(partType),
		Uuids:    uuids,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.NotFound {
				return []*model.Part{}, errs.ErrPartNotFound
			}
		}
		return []*model.Part{}, fmt.Errorf("%s: %w", op, err)
	}
	protoParts := resp.GetParts()
	parts := make([]*model.Part, 0, len(resp.GetParts()))
	for _, part := range protoParts {

		modelPart := &model.Part{
			UUID:          part.GetUuid(),
			Name:          part.GetName(),
			Description:   part.GetDescription(),
			Price:         part.GetPrice(),
			PartType:      partType,
			StockQuantity: part.GetStockQuantity(),
			CreatedAt:     part.GetCreatedAt().AsTime(),
		}
		parts = append(parts, modelPart)
	}
	slog.Info("adapter.inventory кон ")
	return parts, nil
}
