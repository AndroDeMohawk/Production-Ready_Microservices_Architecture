package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
)

func ProtoPartToModel(part *inventoryv1.Part) *model.Part {
	if part == nil {
		return nil
	}
	return &model.Part{
		UUID:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      model.PartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt.AsTime(),
	}
}
