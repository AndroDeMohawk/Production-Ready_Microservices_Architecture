package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/service/input"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToProto(part *model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      ModelPartTypeToProtoPartType(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func ModelPartTypeToProtoPartType(partType model.PartType) inventoryv1.PartType {
	switch partType {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL

	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE

	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD

	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
func ProtoPartTypeToModelPartType(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull

	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine

	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield

	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}

}

func ListPartsRequestToInput(r *inventoryv1.ListPartsRequest) input.ListPartsInput {
	return input.ListPartsInput{
		IDs:      r.Uuids,
		PartType: ProtoPartTypeToModelPartType(r.PartType),
	}
}
