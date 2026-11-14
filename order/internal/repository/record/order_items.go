package record

import "github.com/google/uuid"

type PartType string
type Part struct {
	PartUUID   uuid.UUID
	PartType   PartType
	TotalPrice int64
}
