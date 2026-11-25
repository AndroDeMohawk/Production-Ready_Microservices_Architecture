package model

import "fmt"

type Strength int

const (
	MinHullStrength = 30
	MaxHullStrength = 100
)

type HullProperties struct {
	strength Strength
}

func (h *HullProperties) Strength() Strength { return h.strength }

func NewHullProperties(strength int) (PartProperties, error) {
	st := Strength(strength)
	if st < MinHullStrength || st > MaxHullStrength {
		return PartProperties{}, fmt.Errorf("прочность корпуса должна быть между %d и %d", MinHullStrength, MaxHullStrength)
	}
	return PartProperties{
		hull: &HullProperties{
			strength: st,
		},
	}, nil
}

func (h *HullProperties) CanSupport(e *EngineProperties) bool {
	return int(h.strength) >= int(e.requiredStrength)
}
