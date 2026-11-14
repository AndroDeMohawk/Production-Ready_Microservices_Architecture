package record

import (
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PartType int32
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
	UpdatedAt     *timestamppb.Timestamp
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
