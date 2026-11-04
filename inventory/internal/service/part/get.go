package part

import (
	"context"
	"fmt"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/google/uuid"
)

func (s *Service) Get(ctx context.Context, id string) (*model.Part, error) {
	const op = "inventory.service.get"
	partUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, errs.ErrParsingUUID)
	}
	part, err := s.repository.Get(ctx, partUUID)
	if err != nil {
		return nil, err
	}
	return part, nil

}
