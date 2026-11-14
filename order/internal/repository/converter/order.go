package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
)

func ModelToRecord(order *model.Order) record.Order {
	return record.Order{
		OrderUUID: order.OrderUUID,
		Hull: record.Part{
			PartUUID:   order.Parts.Hull.UUID,
			PartType:   record.PartTypeHull,
			TotalPrice: order.Parts.Hull.Price,
		},
		Engine: record.Part{
			PartUUID:   order.Parts.Engine.UUID,
			PartType:   record.PartTypeEngine,
			TotalPrice: order.Parts.Engine.Price,
		},
		Shield: &record.Part{
			PartUUID:   order.Parts.Shield.UUID,
			PartType:   record.PartTypeShield,
			TotalPrice: order.Parts.Shield.Price,
		},
		Weapon: &record.Part{
			PartUUID:   order.Parts.Weapon.UUID,
			PartType:   record.PartTypeWeapon,
			TotalPrice: order.Parts.Weapon.Price,
		},
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		Status:          record.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func RecordToModel(order *record.Order) model.Order {
	result := model.Order{
		OrderUUID: order.OrderUUID,
		Parts: model.Parts{
			Hull: model.OrderPart{
				UUID:  order.Hull.PartUUID,
				Type:  model.PART_TYPE_HULL,
				Price: order.Hull.TotalPrice,
			},
			Engine: model.OrderPart{
				UUID:  order.Engine.PartUUID,
				Type:  model.PART_TYPE_ENGINE,
				Price: order.Engine.TotalPrice,
			},
		},
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
	if order.Shield != nil {
		result.Parts.Shield = &model.OrderPart{
			UUID:  order.Shield.PartUUID,
			Type:  model.PART_TYPE_SHIELD,
			Price: order.Shield.TotalPrice,
		}
	}
	if order.Weapon != nil {
		result.Parts.Weapon = &model.OrderPart{
			UUID:  order.Weapon.PartUUID,
			Type:  model.PART_TYPE_WEAPON,
			Price: order.Weapon.TotalPrice,
		}
	}
	return result
}
