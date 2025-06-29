package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kvvvseins/mictoservices/services/store/config"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/repository"
	"github.com/kvvvseins/server"
)

// ListHandler хендлер списка товаров.
type ListHandler struct {
	config     *config.Config
	repository *repository.StoreRepository
}

// NewListHandler создает хендлер списка товаров.
func NewListHandler(
	cfg *config.Config,
	repository *repository.StoreRepository,
) http.Handler {
	return &ListHandler{
		config:     cfg,
		repository: repository,
	}
}

func (cu *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	stores, err := cu.repository.GetAll(30)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить список товаров")

		return
	}

	w.WriteHeader(http.StatusOK)

	storesDto := dto.BuildViewsFromModels(stores)

	err = json.NewEncoder(w).Encode(storesDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения списка товаров")

		return
	}
}
