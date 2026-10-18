package handler

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	orderv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/shared/pkg/proto/payment/v1"
	uuid "github.com/google/uuid"
)

// OrderStatus — статус заказа
type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

// PaymentMethod — способ оплаты заказа
type PaymentMethod string

const (
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// Order представляет заказ на постройку космического корабля
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

// orderStore — хранилище заказов (in-memory)
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// handler реализует интерфейс orderv1.Handler, сгенерированный ogen
type handler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт новый обработчик заказов
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *handler {
	return &handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика
func SetupServer(h *handler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации)
// GET /api/v1/orders/{order_uuid}.
func (h *handler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CreateOrder реализует операцию createOrder
func (h *handler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	hullId, err := uuid.Parse(req.GetHullUUID().String())
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга айди корпуса: %w", err)
	}
	engineId, err := uuid.Parse(req.GetEngineUUID().String())
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга айди двигателя : %w", err)
	}
	partIDs := make([]string, 0, 4)
	partIDs = append(partIDs, hullId.String(), engineId.String())
	var shieldID *uuid.UUID
	var weaponID *uuid.UUID
	if req.GetShieldUUID().IsSet() {
		val := req.GetShieldUUID().Value
		shieldID = &val
		partIDs = append(partIDs, val.String())
	}
	if req.GetWeaponUUID().IsSet() {
		val := req.GetWeaponUUID().Value
		weaponID = &val
		partIDs = append(partIDs, val.String())
	}
	parts, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{
		Uuids: partIDs,
	})
	if err != nil {
		return &orderv1.CreateOrderInternalServerError{}, nil

	}
	if len(parts.Parts) < len(partIDs) {
		return &orderv1.CreateOrderNotFound{}, nil
	}
	var totalPrice int64
	for _, part := range parts.Parts {
		if part.StockQuantity <= 0 {
			return &orderv1.CreateOrderConflict{}, nil
		}
		totalPrice += part.Price
	}
	orderId := uuid.New()
	status := OrderStatusPendingPayment
	order := Order{
		OrderUUID:  orderId,
		HullUUID:   hullId,
		EngineUUID: engineId,
		WeaponUUID: weaponID,
		ShieldUUID: shieldID,
		Status:     status,
		TotalPrice: totalPrice,
		CreatedAt:  time.Now(),
	}
	h.store.mu.Lock()
	h.store.orders[order.OrderUUID] = order
	h.store.mu.Unlock()
	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil

}

// PayOrder реализует операцию payOrder
// POST /api/v1/orders/{order_uuid}/pay
func (h *handler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	h.store.mu.RLock()
	orderId := params.OrderUUID
	order, exists := h.store.orders[orderId]
	h.store.mu.RUnlock()
	if !exists {
		return &orderv1.PayOrderNotFound{}, nil
	}
	if order.Status != OrderStatusPendingPayment {
		return &orderv1.PayOrderBadRequest{}, nil
	}
	pay, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{})
	if err != nil {
		return &orderv1.PayOrderInternalServerError{}, nil
	}
	id := uuid.MustParse(pay.TransactionUuid)
	localMethod := PaymentMethod(req.PaymentMethod)
	updatedOrder := Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		Status:          OrderStatusPaid,
		CreatedAt:       order.CreatedAt,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: &id,
		PaymentMethod:   &localMethod,
	}
	h.store.mu.Lock()
	h.store.orders[orderId] = updatedOrder
	h.store.mu.Unlock()
	return &orderv1.PayOrderResponse{
		TransactionUUID: id,
	}, nil

}

// CancelOrder реализует операцию cancelOrder
// POST /api/v1/orders/{order_uuid}/cancel
func (h *handler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.RLock()
	orderId := params.OrderUUID
	order, exists := h.store.orders[orderId]
	h.store.mu.RUnlock()
	if !exists {
		return &orderv1.CancelOrderNotFound{}, nil
	}
	if order.Status != OrderStatusPendingPayment {
		return &orderv1.CancelOrderBadRequest{}, nil
	}
	updateOrder := Order{
		OrderUUID:       order.OrderUUID,
		Status:          OrderStatusCancelled,
		HullUUID:        order.HullUUID,
		WeaponUUID:      order.WeaponUUID,
		ShieldUUID:      order.ShieldUUID,
		EngineUUID:      order.EngineUUID,
		TotalPrice:      order.TotalPrice,
		CreatedAt:       order.CreatedAt,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   order.PaymentMethod,
	}
	h.store.mu.Lock()
	h.store.orders[orderId] = updateOrder
	h.store.mu.Unlock()
	return &orderv1.CancelOrderResponse{}, nil
}
