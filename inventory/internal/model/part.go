package model

import (
	"time"

	"github.com/google/uuid"
)

type PartType string
type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

const (
	PartTypeUnspecified PartType = "Неопределено"
	PartTypeHull        PartType = "Корпус"
	PartTypeEngine      PartType = "Двигатель"
	PartTypeShield      PartType = "Щит"
	PartTypeWeapon      PartType = "Орудие"
)
