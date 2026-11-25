package part

import (
	"context"
	"strings"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/converter"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/inventory/internal/repository/record"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (*model.Part, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const query = `SELECT uuid,name,description,part_type,price,stock_quantity,created_at,updated_at from parts where uuid = $1`
	var part record.Part
	var partType string
	var createdAt time.Time
	err := querier.QueryRow(ctx, query, id).Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&partType,
		&part.Price,
		&part.StockQuantity,
		&createdAt,
	)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFoundInDB) {
			return nil, errs.ErrPartNotFound
		}
		return nil, err
	}
	part.PartType = strings.ToUpper(partType)
	part.CreatedAt = createdAt

	mPart, err := converter.RecordToModel(&part)
	if err != nil {
		return nil, errs.ErrParsingUUID
	}
	return &mPart, nil

}
