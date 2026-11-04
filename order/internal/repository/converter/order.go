package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
)

func ModelToRecord(order *model.Order) record.Order {
	var paymentMethodPtr *record.PaymentMethod
	if order.PaymentMethod != nil {
		paymentMethodPtr = new(record.PaymentMethod)
		*paymentMethodPtr = record.PaymentMethod(*order.PaymentMethod)
	}
	return record.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.Parts.HullUUID,
		EngineUUID:      order.Parts.EngineUUID,
		ShieldUUID:      order.Parts.ShieldUUID,
		WeaponUUID:      order.Parts.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethodPtr,
		Status:          record.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		CanceledAt:      order.CanceledAt,
	}
}

func RecordToModel(order *record.Order) model.Order {
	var paymentMethodPtr *model.PaymentMethod
	if order.PaymentMethod != nil {
		paymentMethodPtr = new(model.PaymentMethod)
		*paymentMethodPtr = model.PaymentMethod(*order.PaymentMethod)
	}
	return model.Order{
		OrderUUID: order.OrderUUID,
		Parts: model.Parts{
			HullUUID:   order.HullUUID,
			EngineUUID: order.EngineUUID,
			ShieldUUID: order.ShieldUUID,
			WeaponUUID: order.WeaponUUID,
		},
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethodPtr,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		CanceledAt:      order.CanceledAt,
	}
}
