package record

import (
	"time"

	"github.com/google/uuid"
)

// TransactionStatus — статус транзакции
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "в обработке"
	TransactionStatusSucceeded TransactionStatus = "успешно"
	TransactionStatusFailed    TransactionStatus = "отклонено"
)

// Transaction — «строка» таблицы transactions
type Transaction struct {
	TransactionUUID uuid.UUID
	OrderUUID       uuid.UUID
	Status          TransactionStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
