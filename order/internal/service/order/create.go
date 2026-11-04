package order

import (
	"context"
	"fmt"
	"log/slog"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/service/input"
	"github.com/google/uuid"
)

func (s *orderService) Create(ctx context.Context, in input.CreateOrder) (*model.Order, error) {
	const op = "order.service.create"
	stringParts := make([]string, 0, 4)
	slog.Info("service.create начало")
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
	fmt.Println("service до вызова метода")
	parts, err := s.inventoryClient.ListParts(ctx, model.PART_TYPE_UNSPECIFIED, stringParts)
	fmt.Println("service после вызова метода")
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка получения  списка запчастей: %w", op, errs.ErrPartNotFound)
	}
	if len(parts) != len(stringParts) {
		return nil, fmt.Errorf("%s: на складе не хватает запчастей: %w", op, errs.ErrPartNotFound)
	}
	var totalPrice int64
	for _, part := range parts {
		if part.StockQuantity <= 0 {
			return nil, fmt.Errorf("%s: одна из запчастей отсутствует на складе %w", op, errs.ErrOutOfStock)
		}
		totalPrice += part.Price
	}
	order, err := s.repository.Create(ctx, model.CreateOrder{
		HullUUID:   in.HullUUID,
		WeaponUUID: in.WeaponUUID,
		ShieldUUID: in.ShieldUUID,
		EngineUUID: in.EngineUUID,
		TotalPrice: totalPrice,
	})
	slog.Info("service.create конец")
	if err != nil {
		return nil, fmt.Errorf(" %s: ошибка создания заказа: %w", op, errs.ErrPartNotFound)
	}
	return order, nil
}
