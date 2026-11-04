package record

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "Unspecified"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

// TransactionStatus — статус транзакции
type TransactionStatus string

const (
	TransactionStatusUnexpected TransactionStatus = "неизвестно"
	TransactionStatusPending    TransactionStatus = "в обработке"
	TransactionStatusSucceeded  TransactionStatus = "успешно"
	TransactionStatusFailed     TransactionStatus = "отклонено"
)

// Transaction — «строка» таблицы transactions
type Transaction struct {
	TransactionUUID uuid.UUID
	OrderUUID       uuid.UUID
	PaymentMethod   PaymentMethod
	Status          TransactionStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
