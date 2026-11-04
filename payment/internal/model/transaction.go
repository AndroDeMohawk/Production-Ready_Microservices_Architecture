package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod int32

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type TransactionStatus int32

const (
	TransactionStatusUnexpected TransactionStatus = iota
	TransactionStatusPending
	TransactionStatusSucceeded
	TransactionStatusFailed
)

type Transaction struct {
	TransactionUUID uuid.UUID
	OrderUUID       uuid.UUID
	Method          PaymentMethod
	Status          TransactionStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
