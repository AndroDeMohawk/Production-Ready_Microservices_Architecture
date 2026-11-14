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
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}
type OrderPart struct {
	UUID  uuid.UUID
	Type  PartType
	Price int64
}
type Parts struct {
	Hull   OrderPart
	Engine OrderPart
	Shield *OrderPart // опциональный
	Weapon *OrderPart // опциональный
}

type CreateOrder struct {
	// UUID корпуса (обязательный, v4).
	Hull OrderPart
	// UUID двигателя (обязательный, v4).
	Engine OrderPart
	// UUID щита (опциональный).
	Shield *OrderPart
	// UUID вооружения (опциональный).
	Weapon     *OrderPart
	TotalPrice int64
}
