package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/record"
)

func RecordToModel(r *record.Transaction) *model.Transaction {
	return &model.Transaction{
		TransactionUUID: r.TransactionUUID,
		OrderUUID:       r.OrderUUID,
		Status:          RecordStatusToModelStatus(r.Status),
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func ModelToRecord(m *model.Transaction) *record.Transaction {
	return &record.Transaction{
		TransactionUUID: m.TransactionUUID,
		OrderUUID:       m.OrderUUID,
		Status:          ModelStatusToRecordStatus(m.Status),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func RecordStatusToModelStatus(status record.TransactionStatus) model.TransactionStatus {
	switch status {
	case record.TransactionStatusPending:
		return model.TransactionStatusPending
	case record.TransactionStatusFailed:
		return model.TransactionStatusFailed
	case record.TransactionStatusSucceeded:
		return model.TransactionStatusSucceeded
	default:
		return model.TransactionStatusUnexpected
	}
}

func ModelStatusToRecordStatus(status model.TransactionStatus) record.TransactionStatus {
	switch status {
	case model.TransactionStatusPending:
		return record.TransactionStatusPending
	case model.TransactionStatusFailed:
		return record.TransactionStatusFailed
	case model.TransactionStatusSucceeded:
		return record.TransactionStatusSucceeded
	default:
		return ""
	}
}
