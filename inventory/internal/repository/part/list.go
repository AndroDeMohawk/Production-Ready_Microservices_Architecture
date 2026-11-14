package part

import (
	"context"
	"strings"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/converter"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/record"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *repository) List(ctx context.Context, uuids []uuid.UUID) ([]*model.Part, error) {
	parts := make([]*model.Part, 0)
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const query = `SELECT uuid,name,description,part_type,price,stock_quantity,created_at,updated_at from parts where uuid = $1`
	for _, id := range uuids {
		var part record.Part
		var partType string
		var createdAt time.Time
		var updatedAt *time.Time

		err := querier.QueryRow(ctx, query, id).Scan(
			&part.UUID,
			&part.Name,
			&part.Description,
			&partType,
			&part.Price,
			&part.StockQuantity,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		part.PartType = record.PartTypeFromString(strings.ToUpper(partType))
		part.CreatedAt = timestamppb.New(createdAt)
		if updatedAt != nil {
			part.UpdatedAt = timestamppb.New(*updatedAt)
		}
		convPart, err := converter.RecordToModel(&part)
		if err != nil {
			continue
		}
		parts = append(parts, convPart)
	}
	if len(parts) < len(uuids) {
		return nil, errs.ErrPartsNotFound
	}
	return parts, nil

}
