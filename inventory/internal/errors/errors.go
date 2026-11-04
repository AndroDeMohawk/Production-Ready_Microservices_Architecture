package errors

import "github.com/go-faster/errors"

var (

	ErrParsingUUID   = errors.New("ошибка парсинга uuid")
	ErrPartNotFound  = errors.New("запчасть не найдена")
	ErrPartsNotFound = errors.New("одна или несколько запчастей не найдены")
)
