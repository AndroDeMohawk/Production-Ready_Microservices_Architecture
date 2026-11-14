package model

import (
	"time"

	"github.com/google/uuid"
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
	Status          TransactionStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
