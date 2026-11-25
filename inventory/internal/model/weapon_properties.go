package model

import "fmt"

type WeaponType string

const (
	WeaponTypeLaser   WeaponType = "laser"
	WeaponTypeMissile WeaponType = "missile"
)

type WeaponProperties struct {
	weaponType WeaponType
}

func (w *WeaponProperties) WeaponType() WeaponType { return w.weaponType }

func NewWeaponProperties(weaponType string) (PartProperties, error) {
	wt := WeaponType(weaponType)
	switch wt {
	case WeaponTypeLaser, WeaponTypeMissile:
		return PartProperties{
			weapon: &WeaponProperties{
				weaponType: wt,
			},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("некорректный тип оружия")
	}
}
