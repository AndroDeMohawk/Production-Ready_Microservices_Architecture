package order

import (
	"context"
	"time"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, in model.CreateOrder) (*model.Order, error) {
	querier := r.getter.DefaultTrOrDB(ctx, r.pool)

	order := &model.Order{
		OrderUUID: uuid.New(),
		Parts: model.Parts{
			Hull:   in.Hull,
			Engine: in.Engine,
			Weapon: in.Weapon,
			Shield: in.Shield,
		},
		TotalPrice: in.TotalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}

	const insertQueryOrders = `INSERT INTO orders (uuid, status, created_at) VALUES ($1, $2, $3)`
	tag, err := querier.Exec(ctx, insertQueryOrders, order.OrderUUID, order.Status, order.CreatedAt)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, errs.ErrOrderNotFound
	}
	const insertQueryOrderItems = `INSERT INTO order_items (order_uuid, part_uuid, part_type, price) VALUES ($1, $2, $3, $4)`
	tag, err = querier.Exec(ctx, insertQueryOrderItems, order.OrderUUID, in.Hull.UUID, partTypeName(in.Hull.Type), in.Hull.Price)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, errs.ErrOrderNotFound
	}
	tag, err = querier.Exec(ctx, insertQueryOrderItems, order.OrderUUID, in.Engine.UUID, partTypeName(in.Engine.Type), in.Engine.Price)
	if err != nil {
		return nil, err
	}
	if in.Shield != nil {
		tag, err := querier.Exec(ctx, insertQueryOrderItems, order.OrderUUID, in.Shield.UUID, partTypeName(in.Shield.Type), in.Shield.Price)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			return nil, errs.ErrOrderNotFound
		}
	}
	if in.Weapon != nil {
		tag, err := querier.Exec(ctx, insertQueryOrderItems, order.OrderUUID, in.Weapon.UUID, partTypeName(in.Weapon.Type), in.Weapon.Price)
		if err != nil {
			return nil, err
		}
		if tag.RowsAffected() == 0 {
			return nil, errs.ErrOrderNotFound
		}
	}

	return order, nil
}

func partTypeName(partType model.PartType) string {
	switch partType {
	case model.PART_TYPE_HULL:
		return "hull"
	case model.PART_TYPE_ENGINE:
		return "engine"
	case model.PART_TYPE_SHIELD:
		return "shield"
	case model.PART_TYPE_WEAPON:
		return "weapon"
	default:
		return ""
	}
}
