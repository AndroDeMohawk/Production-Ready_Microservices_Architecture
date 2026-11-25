package model

import "fmt"

type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

type ShieldProperties struct {
	shieldType ShieldType
}

func (s *ShieldProperties) ShieldType() ShieldType { return s.shieldType }

func NewShieldProperties(shieldType string) (PartProperties, error) {
	st := ShieldType(shieldType)
	switch st {
	case ShieldTypeEnergy, ShieldTypePlasma:
		return PartProperties{
			shield: &ShieldProperties{
				shieldType: st,
			},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("некорректный тип щита")
	}
}

func (s *ShieldProperties) ConflictsWith(w *WeaponProperties) bool {
	return s.shieldType == ShieldTypePlasma && w.weaponType == WeaponTypeLaser
}
