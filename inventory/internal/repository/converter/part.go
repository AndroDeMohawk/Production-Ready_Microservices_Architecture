package converter

import (
	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/record"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func RecordToModel(r *record.Part) (*model.Part, error) {
	uuid, err := uuid.Parse(r.UUID)
	if err != nil {
		return nil, errs.ErrParsingUUID
	}
	return &model.Part{
		UUID:          uuid,
		Name:          r.Name,
		Description:   r.Description,
		Price:         r.Price,
		PartType:      RecordPartTypeToModelPartType(r.PartType),
		StockQuantity: r.StockQuantity,
		CreatedAt:     r.CreatedAt.AsTime(),
	}, nil
}
func ModelToRecord(m *model.Part) *record.Part {
	return &record.Part{
		UUID:          m.UUID.String(),
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		PartType:      ModelPartTypeToRecordPartType(m.PartType),
		StockQuantity: m.StockQuantity,
		CreatedAt:     timestamppb.New(m.CreatedAt),
	}
}

func ModelPartTypeToRecordPartType(partType model.PartType) record.PartType {
	switch partType {
	case model.PartTypeHull:
		return record.PartTypeHull

	case model.PartTypeEngine:
		return record.PartTypeEngine

	case model.PartTypeShield:
		return record.PartTypeShield

	case model.PartTypeWeapon:
		return record.PartTypeWeapon
	default:
		return record.PartTypeUnspecified
	}

}

func RecordPartTypeToModelPartType(partType record.PartType) model.PartType {
	switch partType {
	case record.PartTypeHull:
		return model.PartTypeHull
	case record.PartTypeEngine:
		return model.PartTypeEngine
	case record.PartTypeShield:
		return model.PartTypeShield
	case record.PartTypeWeapon:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}
