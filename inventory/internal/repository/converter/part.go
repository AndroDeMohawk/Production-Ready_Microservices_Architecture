package converter

import (
	"fmt"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/record"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/json"
)

func RecordToModel(r *record.Part) (model.Part, error) {
	var partProps record.PartPropertiesRecord
	if err := json.Unmarshal(r.Properties, &partProps); err != nil {
		return model.Part{}, fmt.Errorf("десериализовать свойства: %w", err)
	}
	props, err := partPropertiesFromRecord(partProps)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать свойства: %w", err)
	}

	partType, err := model.NewPartType(r.PartType)
	if err != nil {
		return model.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	return model.RestorePart(
		uuid.MustParse(r.UUID),
		r.Name,
		r.Description,
		partType,
		r.Price,
		r.StockQuantity,
		r.Reserved,
		props,
		r.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartPropertiesRecord) (model.PartProperties, error) {
	switch {
	case rec.Hull != nil:
		return model.NewHullProperties(rec.Hull.Strength)
	case rec.Engine != nil:
		return model.NewEngineProperties(rec.Engine.Class, rec.Engine.RequiredStrength)
	case rec.Shield != nil:
		return model.NewShieldProperties(rec.Shield.ShieldType)
	case rec.Weapon != nil:
		return model.NewWeaponProperties(rec.Weapon.WeaponType)
	default:
		return model.PartProperties{}, nil
	}
}

//func ModelToRecord(m *model.Part) *record.Part {
//	return &record.Part{
//		UUID:          m.UUID.String(),
//		Name:          m.Name,
//		Description:   m.Description,
//		Price:         m.Price,
//		PartType:      ModelPartTypeToRecordPartType(m.PartType),
//		StockQuantity: m.StockQuantity,
//		CreatedAt:     timestamppb.New(m.CreatedAt),
//	}
//}
//
//func ModelPartTypeToRecordPartType(partType model.PartType) record.PartType {
//	switch partType {
//	case model.PartTypeHull:
//		return record.PartTypeHull
//
//	case model.PartTypeEngine:
//		return record.PartTypeEngine
//
//	case model.PartTypeShield:
//		return record.PartTypeShield
//
//	case model.PartTypeWeapon:
//		return record.PartTypeWeapon
//	default:
//		return record.PartTypeUnspecified
//	}
//
//}
//
//func RecordPartTypeToModelPartType(partType record.PartType) model.PartType {
//	switch partType {
//	case record.PartTypeHull:
//		return model.PartTypeHull
//	case record.PartTypeEngine:
//		return model.PartTypeEngine
//	case record.PartTypeShield:
//		return model.PartTypeShield
//	case record.PartTypeWeapon:
//		return model.PartTypeWeapon
//	default:
//		return model.PartTypeUnspecified
//	}
//}
