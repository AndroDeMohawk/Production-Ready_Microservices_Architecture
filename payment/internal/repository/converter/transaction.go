package converter

import (
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/model"
	"github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/payment/internal/repository/record"
)

func RecordToModel(r *record.Transaction) *model.Transaction {
	return &model.Transaction{
		TransactionUUID: r.TransactionUUID,
		OrderUUID:       r.OrderUUID,
		Method:          RecordMethodToModelMethod(r.PaymentMethod),
		Status:          RecordStatusToModelStatus(r.Status),
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func ModelToRecord(m *model.Transaction) *record.Transaction {
	return &record.Transaction{
		TransactionUUID: m.TransactionUUID,
		OrderUUID:       m.OrderUUID,
		PaymentMethod:   ModelMethodToRecordMethod(m.Method),
		Status:          ModelStatusToRecordStatus(m.Status),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func RecordMethodToModelMethod(method record.PaymentMethod) model.PaymentMethod {
	switch method {
	case record.PaymentMethodCard:
		return model.PaymentMethodCard
	case record.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case record.PaymentMethodCreditCard:
		return model.PaymentMethodCreditCard
	case record.PaymentMethodInvestorMoney:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

func ModelMethodToRecordMethod(method model.PaymentMethod) record.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return record.PaymentMethodCard
	case model.PaymentMethodSBP:
		return record.PaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return record.PaymentMethodCreditCard
	case model.PaymentMethodInvestorMoney:
		return record.PaymentMethodInvestorMoney
	default:
		return record.PaymentMethodUnspecified
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
		return record.TransactionStatusUnexpected
	}
}
