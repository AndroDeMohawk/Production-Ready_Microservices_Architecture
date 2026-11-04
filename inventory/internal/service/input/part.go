package input

import "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"

type ListPartsInput struct {
	IDs      []string
	PartType model.PartType
}
