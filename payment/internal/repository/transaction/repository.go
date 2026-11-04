package transaction

import (
	"sync"

	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/record"
	"github.com/google/uuid"
)

type repository struct {
	mu      sync.RWMutex
	txStore map[uuid.UUID]record.Transaction
}

func New() *repository {
	return &repository{
		txStore: make(map[uuid.UUID]record.Transaction),
	}
}
