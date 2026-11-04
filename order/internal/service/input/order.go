package input

import (
	"github.com/google/uuid"
)

type CreateOrder struct {
	// UUID корпуса (обязательный, v4).
	HullUUID uuid.UUID
	// UUID двигателя (обязательный, v4).
	EngineUUID uuid.UUID
	// UUID щита (опциональный).
	ShieldUUID *uuid.UUID
	// UUID вооружения (опциональный).
	WeaponUUID *uuid.UUID
}
