package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/pinger/server"
	"github.com/pkg/errors"
)

// ErrorResponse структура ошибки от ProfileHandler.
type ErrorResponse struct {
	Message string `json:"message"`
}

var guidNotFoundError = errors.New("в заголовке не найден Guid пользователя")
var invalidGuidError = errors.New("в заголовке не верный Guid пользователя")

func textErrorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
	errMsg string,
) {
	if errors.Is(err, guidNotFoundError) || errors.Is(err, invalidGuidError) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})

		return
	}

	if err != nil {
		server.GetLogger(ctx).Error("ошибка profile crud", "err", err)
	}

	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: errMsg})
}

func getGuidFromRequest(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	guid := r.Header.Get("X-User-Id")
	fmt.Println(guid)
	if guid == "" {
		textErrorResponse(r.Context(), w, guidNotFoundError, "")

		return uuid.UUID{}, false
	}

	guidParsed, err := uuid.Parse(guid)
	if err != nil {
		textErrorResponse(r.Context(), w, invalidGuidError, "")

		return uuid.UUID{}, false
	}

	return guidParsed, true
}
