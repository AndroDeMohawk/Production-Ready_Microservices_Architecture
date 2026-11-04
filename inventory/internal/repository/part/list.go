package part

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/converter"
	"github.com/google/uuid"
)

func (p *partStore) List(ctx context.Context, uuids []uuid.UUID) ([]*model.Part, error) {
	parts := make([]*model.Part, 0)
	for _, id := range uuids {
		part, exist := p.parts[id]
		if !exist {
			continue
		}
		convPart, err := converter.RecordToModel(&part)
		if err != nil {
			continue
		}
		parts = append(parts, convPart)
	}
	if len(parts) < len(uuids) {
		return parts, errs.ErrPartsNotFound
	}
	return parts, nil

}
