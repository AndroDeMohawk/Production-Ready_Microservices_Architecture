package order

import (
	"context"
	"errors"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/converter"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) Get(ctx context.Context, id uuid.UUID) (*model.Order, error) {

	querier := r.getter.DefaultTrOrDB(ctx, r.pool)
	const orderQuery = `
		SELECT uuid, status, transaction_uuid, created_at, updated_at
		FROM orders
		WHERE uuid = $1;`
	var rec record.Order
	err := querier.QueryRow(ctx, orderQuery, id).Scan(
		&rec.OrderUUID,
		&rec.Status,
		&rec.TransactionUUID,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrOrderNotFound
		}
		return nil, err
	}

	const partsQuery = `SELECT part_type,part_uuid,price FROM order_items WHERE order_uuid = $1;`
	rows, err := querier.Query(ctx, partsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var totalPrice int64
	for rows.Next() {
		var part record.Part
		if err := rows.Scan(&part.PartType, &part.PartUUID, &part.TotalPrice); err != nil {
			return nil, err
		}
		totalPrice += part.TotalPrice
		switch part.PartType {
		case record.PartTypeHull:
			rec.Hull = record.Part{
				PartUUID:   part.PartUUID,
				PartType:   part.PartType,
				TotalPrice: part.TotalPrice,
			}

		case record.PartTypeEngine:
			rec.Engine = record.Part{
				PartUUID:   part.PartUUID,
				PartType:   part.PartType,
				TotalPrice: part.TotalPrice,
			}
		case record.PartTypeShield:
			rec.Shield = &record.Part{
				PartUUID:   part.PartUUID,
				PartType:   part.PartType,
				TotalPrice: part.TotalPrice,
			}
		case record.PartTypeWeapon:
			rec.Weapon = &record.Part{
				PartUUID:   part.PartUUID,
				PartType:   part.PartType,
				TotalPrice: part.TotalPrice,
			}

		}

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	rec.TotalPrice = totalPrice

	orderModel := converter.RecordToModel(&rec)
	return &orderModel, nil

}
