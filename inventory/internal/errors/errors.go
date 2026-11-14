package errors

import (
	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
)

var (
	ErrParsingUUID      = errors.New("ошибка парсинга uuid")
	ErrPartNotFoundInDB = pgx.ErrNoRows
	ErrPartNotFound     = errors.New("запчасть не найдена")
	ErrPartsNotFound    = errors.New("одна или несколько запчастей не найдены")
)
