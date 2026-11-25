package model

import (
	"fmt"
)

type EngineClass string
type RequireStrength int

const (
	EngineClassA EngineClass = "A"
	EngineCLassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

type EngineProperties struct {
	class            EngineClass
	requiredStrength RequireStrength
}

func (e *EngineProperties) Class() EngineClass {
	return e.class
}

func NewEngineProperties(class string, requiredStrength int) (PartProperties, error) {
	cl := EngineClass(class)
	reqSt := RequireStrength(requiredStrength)
	switch cl {
	case EngineClassA, EngineCLassB, EngineClassC:
		if requiredStrength < MinHullStrength || requiredStrength > MaxHullStrength {
			return PartProperties{}, fmt.Errorf("недопустимая прочность двигателя")
		}
		return PartProperties{
			engine: &EngineProperties{
				class:            cl,
				requiredStrength: reqSt,
			},
		}, nil
	default:
		return PartProperties{}, fmt.Errorf("некорректный класс двигателя")
	}
}
