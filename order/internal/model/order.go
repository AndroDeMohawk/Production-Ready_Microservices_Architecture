package model

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "ожидание оплаты"
	OrderStatusPaid           OrderStatus = "оплачен"
	OrderStatusCancelled      OrderStatus = "отменен"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	Parts           Parts
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
	CanceledAt      *time.Time
}

type Parts struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID // опциональный
	WeaponUUID *uuid.UUID // опциональный
}

type CreateOrder struct {
	// UUID корпуса (обязательный, v4).
	HullUUID uuid.UUID
	// UUID двигателя (обязательный, v4).
	EngineUUID uuid.UUID
	// UUID щита (опциональный).
	ShieldUUID *uuid.UUID
	// UUID вооружения (опциональный).
	WeaponUUID *uuid.UUID
	TotalPrice int64
}
