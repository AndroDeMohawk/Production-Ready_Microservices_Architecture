package model

import "time"

type PartType int32

const (
	PART_TYPE_UNSPECIFIED PartType = 0
	PART_TYPE_HULL        PartType = 1
	PART_TYPE_ENGINE      PartType = 2
	PART_TYPE_SHIELD      PartType = 3
	PART_TYPE_WEAPON      PartType = 4
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}
