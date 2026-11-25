package record

import (
	"strings"
	"time"
)

type PartType int32
type Part struct {
	UUID          string
	Name          string
	Description   string
	PartType      string
	Price         int64
	StockQuantity int
	Reserved      int
	Properties    []byte // JSONB из PostgreSQL
	CreatedAt     time.Time
}

const (
	PartTypeUnspecified PartType = 0
	PartTypeHull        PartType = 1
	PartTypeEngine      PartType = 2
	PartTypeShield      PartType = 3
	PartTypeWeapon      PartType = 4
)

func PartTypeFromString(partType string) PartType {
	switch strings.ToUpper(partType) {
	case "HULL":
		return PartTypeHull
	case "ENGINE":
		return PartTypeEngine
	case "SHIELD":
		return PartTypeShield
	case "WEAPON":
		return PartTypeWeapon
	default:
		return PartTypeUnspecified
	}
}
