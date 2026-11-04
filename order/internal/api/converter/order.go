package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/input"
	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func CreateOrderRequestToInput(request *orderv1.CreateOrderRequest) input.CreateOrder {
	var weaponUUID, shieldUUID *uuid.UUID
	if request.WeaponUUID.Set && !request.WeaponUUID.Null {
		id := request.WeaponUUID.Value
		weaponUUID = &id
	}
	if request.ShieldUUID.Set && !request.ShieldUUID.Null {
		id := request.ShieldUUID.Value
		shieldUUID = &id
	}

	return input.CreateOrder{
		HullUUID:   request.HullUUID,
		EngineUUID: request.EngineUUID,
		WeaponUUID: weaponUUID,
		ShieldUUID: shieldUUID,
	}
}

func OrderToDTO(o *model.Order) orderv1.OrderDto {
	dto := orderv1.OrderDto{
		OrderUUID:       o.OrderUUID,
		HullUUID:        o.Parts.HullUUID,
		EngineUUID:      o.Parts.EngineUUID,
		TotalPrice:      o.TotalPrice,
		Status:          orderv1.OrderStatus(o.Status),
		CreatedAt:       o.CreatedAt,
		WeaponUUID:      orderv1.OptNilUUID{Set: true, Null: true},
		ShieldUUID:      orderv1.OptNilUUID{Set: true, Null: true},
		TransactionUUID: orderv1.OptNilUUID{Set: true, Null: true},
		PaymentMethod:   orderv1.OptNilPaymentMethod{Set: true, Null: true},
	}

	if o.Parts.WeaponUUID != nil {
		dto.WeaponUUID = orderv1.OptNilUUID{Value: *o.Parts.WeaponUUID, Set: true}
	}
	if o.Parts.ShieldUUID != nil {
		dto.ShieldUUID = orderv1.OptNilUUID{Value: *o.Parts.ShieldUUID, Set: true}
	}
	if o.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.OptNilUUID{Value: *o.TransactionUUID, Set: true}
	}
	if o.PaymentMethod != nil {
		dto.PaymentMethod = orderv1.OptNilPaymentMethod{
			Value: orderv1.PaymentMethod(*o.PaymentMethod),
			Set:   true,
		}
	}
	return dto
}
