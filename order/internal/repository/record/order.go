package record

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus — статус заказа
type OrderStatus string

// PaymentMethod — способ оплаты заказа
type PaymentMethod string
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	CanceledAt      *time.Time
}

const (
	OrderStatusPaid      OrderStatus = "оплачен"
	OrderStatusCancelled OrderStatus = "отменен"
)
