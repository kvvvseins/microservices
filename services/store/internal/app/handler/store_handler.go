package handler

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/store/config"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// StoreHandler хендлер создания store.
type StoreHandler struct {
	config     *config.Config
	repository *repository.StoreRepository
}

// NewStoreHandler создает хендлер crud для store.
func NewStoreHandler(
	cfg *config.Config,
	repository *repository.StoreRepository,
) http.Handler {
	return &StoreHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterStoreHandlers регистрирует crud ручки остатка
func RegisterStoreHandlers(router *chi.Mux, handler http.Handler) {
	router.Method(
		http.MethodGet,
		"/",
		handler,
	)
	router.Method(
		http.MethodPost,
		"/",
		handler,
	)
}

func (cu *StoreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// todo вынести в сервис
	switch r.Method {
	case http.MethodGet:
		cu.get(w, r)
	case http.MethodPost:
		cu.create(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *StoreHandler) get(w http.ResponseWriter, r *http.Request) {
	guidsQuery, ok := r.URL.Query()["guid"]
	if !ok || len(guidsQuery) < 1 {
		server.ErrorResponseOutput(r.Context(), w, nil, "не указан guid")

		return
	}

	for _, guid := range guidsQuery {
		_, err := uuid.Parse(guid)
		if err != nil {
			server.ErrorResponseOutput(r.Context(), w, nil, "не валидный guid")

			return
		}
	}

	stores, err := cu.repository.FindByGuids(guidsQuery)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить остаток")

		return
	}

	w.WriteHeader(http.StatusOK)

	storesDto := dto.BuildViewsFromModels(stores)

	err = json.NewEncoder(w).Encode(storesDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения остатка")

		return
	}
}

func (cu *StoreHandler) create(w http.ResponseWriter, r *http.Request) {
	userID, isOk := server.GetUserIDFromRequest(w, r)
	if !isOk {
		return
	}

	var storeDto dto.CreateStore
	err := json.NewDecoder(r.Body).Decode(&storeDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json нового остатка")

		return
	}

	if storeDto.Price <= 0 {
		err = errors.New("не верная цена остатка")
		server.ErrorResponseOutput(r.Context(), w, err, err.Error())

		return
	}

	store, err := cu.repository.FindByName(storeDto.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка проверки существующего товара")

		return
	}

	if store != nil {
		quantity, errQuantity := cu.repository.GetBalance(store.ID)
		if errQuantity != nil {
			server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения баланса")

			return
		}

		cu.responseStore(store, quantity, w, r, http.StatusCreated)

		return
	}

	var reserve *model.Reserve

	store = &model.Store{
		Price: storeDto.Price,
		Name:  storeDto.Name,
	}

	err = cu.config.GetDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(store)
		if result.Error != nil {
			return result.Error
		}

		var errReserve error
		reserve, errReserve = cu.repository.Reserve(tx, int(storeDto.Quantity), store.ID, nil, userID)
		if errReserve != nil {
			return errReserve
		}

		return nil
	})

	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка создания остатка")

		return
	}

	cu.responseStore(store, reserve.Quantity, w, r, http.StatusCreated)
}

func (cu *StoreHandler) responseStore(store *model.Store, quantity int, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewStore{
		Guid:     store.Guid,
		Price:    store.Price,
		Name:     store.Name,
		Quantity: quantity,
	})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка остатка")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
