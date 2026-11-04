package order

import (
	"context"
	"time"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *orderStore) Create(ctx context.Context, in model.CreateOrder) (*model.Order, error) {
	order := &model.Order{
		OrderUUID: uuid.New(),
		Parts: model.Parts{
			HullUUID:   in.HullUUID,
			EngineUUID: in.EngineUUID,
			WeaponUUID: in.WeaponUUID,
			ShieldUUID: in.ShieldUUID,
		},
		TotalPrice: in.TotalPrice,
		Status:     model.OrderStatusPendingPayment,
		CreatedAt:  time.Now(),
	}
	createdOrder := converter.ModelToRecord(order)
	r.mu.Lock()
	r.orders[order.OrderUUID] = createdOrder
	r.mu.Unlock()
	return order, nil
}
