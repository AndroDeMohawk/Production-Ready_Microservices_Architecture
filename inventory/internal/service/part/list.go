package part

import (
	"context"
	"fmt"
	"sort"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/service/input"
	"github.com/google/uuid"
)

func (s *Service) List(ctx context.Context, in input.ListPartsInput) ([]*model.Part, error) {
	const op = "inventory.service.list"
	uuids := make([]uuid.UUID, 0, len(in.IDs))
	for _, id := range in.IDs {
		newId, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("%s %w", op, errs.ErrParsingUUID)
		}
		uuids = append(uuids, newId)
	}
	parts, err := s.repository.List(ctx, uuids)

	if in.PartType != "" {
		selectedType := in.PartType

		sort.SliceStable(parts, func(i, j int) bool {
			iSelected := parts[i].PartType == selectedType
			jSelected := parts[j].PartType == selectedType

			return iSelected && !jSelected
		})
	}
	if err != nil {
		return parts, fmt.Errorf("%s %w", op, err)
	}

	return parts, nil
}
