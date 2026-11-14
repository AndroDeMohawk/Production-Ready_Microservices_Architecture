package record

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus — статус заказа
type OrderStatus string

type Order struct {
	OrderUUID       uuid.UUID
	Hull            Part
	Engine          Part
	Shield          *Part // опциональный
	Weapon          *Part // опциональный
	TotalPrice      int64 // в копейках
	TransactionUUID *uuid.UUID
	Status          OrderStatus
	CreatedAt       time.Time
	UpdatedAt       *time.Time
}

const (
	OrderStatusPaid      OrderStatus = "оплачен"
	OrderStatusCancelled OrderStatus = "отменен"

	PartTypeHull   = "hull"
	PartTypeEngine = "engine"
	PartTypeShield = "shield"
	PartTypeWeapon = "weapon"
)
