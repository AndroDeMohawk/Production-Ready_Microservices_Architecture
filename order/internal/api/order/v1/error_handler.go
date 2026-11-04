package v1

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	errs "github.com/AndroDeMohawk/Production-Ready_Microservices_Architecture/order/internal/errors"
	"github.com/ogen-go/ogen/ogenerrors"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
	code, message := mapError(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if encErr := json.NewEncoder(w).Encode(errorResponse{Code: code, Message: message}); encErr != nil {
		slog.ErrorContext(ctx, "ошибка кодирования ответа", "error", encErr)
	}
}

func mapError(err error) (int, string) {
	var decodeParams *ogenerrors.DecodeParamsError
	var decodeRequest *ogenerrors.DecodeRequestError

	switch {
	case errors.As(err, &decodeParams), errors.As(err, &decodeRequest):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, errs.ErrPartNotFound),
		errors.Is(err, errs.ErrOrderNotFound):

		return http.StatusNotFound, err.Error()
	case errors.Is(err, errs.ErrOutOfStock),
		errors.Is(err, errs.ErrOrderCancelled),
		errors.Is(err, errs.ErrOrderAlreadyPaid):
		return http.StatusConflict, err.Error()
	case errors.Is(err, errs.ErrInvalidPaymentMethod),
		errors.Is(err, errs.ErrInvalidUUID):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
