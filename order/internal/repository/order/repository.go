package order

import (
	"sync"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/repository/record"
	"github.com/google/uuid"
)

// orderStore — хранилище заказов (in-memory)
type orderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]record.Order
}

// NewOrderStore создаёт новое пустое хранилище заказов
func NewOrderStore() *orderStore {
	return &orderStore{
		orders: make(map[uuid.UUID]record.Order),
	}
}
