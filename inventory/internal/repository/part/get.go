package part

import (
	"context"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/converter"
	"github.com/google/uuid"
)

func (p *partStore) Get(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	order, exist := p.parts[id]
	if !exist {
		return nil, errs.ErrPartNotFound
	}
	mOrder, err := converter.RecordToModel(&order)
	if err != nil {
		return nil, errs.ErrParsingUUID
	}
	return mOrder, nil

}
