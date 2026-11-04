package errors

import "github.com/go-faster/errors"

var (
	ErrParsingUUID         = errors.New("ошибкаа парсинга uuid")
	ErrEmptyOrderID        = errors.New("пустой номер заказа")
	ErrDuplicatePayOrder   = errors.New("дублирование транзакции")
	ErrTransactionNotFound = errors.New("транзакция не найдена")
)
