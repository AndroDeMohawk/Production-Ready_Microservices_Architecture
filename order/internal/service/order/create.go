package order

import (
	"context"
	"fmt"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/input"
	"github.com/google/uuid"
)

func (s *orderService) Create(ctx context.Context, in input.CreateOrder) (*model.Order, error) {
	const op = "order.service.create"
	stringParts := make([]string, 0, 4)
	if in.HullUUID == uuid.Nil || in.EngineUUID == uuid.Nil {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrInvalidUUID)
	}
	stringParts = append(stringParts, in.HullUUID.String(), in.EngineUUID.String())
	weaponUUID := in.WeaponUUID
	ShieldUUID := in.ShieldUUID
	if weaponUUID != nil {
		stringParts = append(stringParts, weaponUUID.String())
	}

	if ShieldUUID != nil {
		stringParts = append(stringParts, ShieldUUID.String())
	}
	parts, err := s.inventoryClient.ListParts(ctx, model.PART_TYPE_UNSPECIFIED, stringParts)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка получения  списка запчастей: %w", op, errs.ErrPartNotFound)
	}
	if len(parts) != len(stringParts) {
		return nil, fmt.Errorf("%s: на складе не хватает запчастей: %w", op, errs.ErrPartNotFound)
	}

	orderParts, total, err := buildOrderParts(in, parts)
	if err != nil {
		return nil, fmt.Errorf("%s: проверка деталей: %w", op, err)
	}
	var order *model.Order
	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		order, err = s.repository.Create(ctx, model.CreateOrder{
			Hull:       orderParts.Hull,
			Engine:     orderParts.Engine,
			Shield:     orderParts.Shield,
			Weapon:     orderParts.Weapon,
			TotalPrice: total,
		})
		return err
	})
	if err != nil {
		return nil, fmt.Errorf(" %s: ошибка создания заказа: %w", op, err)
	}
	return order, nil
}

func buildOrderParts(in input.CreateOrder, parts []*model.Part) (model.Parts, int64, error) {
	var (
		result model.Parts
		total  int64
	)

	for _, p := range parts {
		if p.StockQuantity <= 0 {
			return model.Parts{}, 0, errs.ErrOutOfStock
		}

		partUUID, err := uuid.Parse(p.UUID)
		if err != nil {
			return model.Parts{}, 0, fmt.Errorf("%w: %s", errs.ErrInvalidUUID, p.UUID)
		}

		part := model.OrderPart{
			UUID:  partUUID,
			Type:  p.PartType,
			Price: p.Price,
		}
		total += p.Price

		switch p.PartType {
		case model.PART_TYPE_HULL:
			result.Hull = part
		case model.PART_TYPE_ENGINE:
			result.Engine = part
		case model.PART_TYPE_SHIELD:
			result.Shield = &part
		case model.PART_TYPE_WEAPON:
			result.Weapon = &part
		}
	}

	if (in.ShieldUUID != nil) != (result.Shield != nil) ||
		(in.WeaponUUID != nil) != (result.Weapon != nil) {
		return model.Parts{}, 0, errs.ErrPartNotFound
	}

	return result, total, nil
}
