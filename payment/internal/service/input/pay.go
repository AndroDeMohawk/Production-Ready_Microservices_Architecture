package input

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/google/uuid"
)

type Pay struct {
	OrderId uuid.UUID
	Method  model.PaymentMethod
}
