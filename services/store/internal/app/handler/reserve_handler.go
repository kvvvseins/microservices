package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/store/config"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/store/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ReserveHandler хендлер резервирования.
type ReserveHandler struct {
	config     *config.Config
	repository *repository.StoreRepository
}

// NewReserveHandler создает хендлер резервирования.
func NewReserveHandler(
	cfg *config.Config,
	repository *repository.StoreRepository,
) http.Handler {
	return &ReserveHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterReserveHandlers регистрирует ручки резервирования
func RegisterReserveHandlers(router *chi.Mux, handler http.Handler) {
	router.Method(
		http.MethodPost,
		"/reserve",
		handler,
	)
}

func (cu *ReserveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, isOk := server.GetUserIDFromRequest(w, r)
	if !isOk {
		return
	}

	var reserveDto dto.Reserve
	err := json.NewDecoder(r.Body).Decode(&reserveDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json резервирования остатка")

		return
	}

	orderId, err := uuid.Parse(reserveDto.OrderId)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верный гуид")

		return
	}

	if len(reserveDto.Products) == 0 {
		err = errors.New("не указаны товары для резервирования")
		server.ErrorResponseOutput(r.Context(), w, err, err.Error())

		return
	}

	guids := make([]string, 0, len(reserveDto.Products))

	for _, reserveDto := range reserveDto.Products {
		guid, errParse := uuid.Parse(reserveDto.Guid)
		if errParse != nil {
			server.ErrorResponseOutput(r.Context(), w, nil, "не верные guid товара")

			return
		}

		guids = append(guids, guid.String())
	}

	//@todo тут лок нужен на товары
	stores, err := cu.repository.FindByGuids(guids)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения товара")

		return
	}

	if len(stores) == 0 || len(stores) != len(reserveDto.Products) {
		err = errors.New("не найдены все товары")
		server.ErrorResponseOutput(r.Context(), w, err, err.Error())

		return
	}

	err = cu.config.GetDb().Transaction(func(tx *gorm.DB) error {
		for _, store := range stores {
			reserveDtoQuantity, errQ := findQuantity(reserveDto.Products, store.Guid.String())
			if errQ != nil {
				return errQ
			}

			quantity := store.Quantity + reserveDtoQuantity
			if quantity < 0 {
				return errors.Wrap(err, "недостаточно товара на складе")
			}

			_, err = cu.repository.Reserve(reserveDtoQuantity, store.ID, &orderId, userID)
			if err != nil {
				return errors.Wrap(err, "ошибка резервирования")
			}
		}

		return nil
	})

	w.WriteHeader(http.StatusCreated)
}

func findQuantity(reserveDtos []dto.ReserveProduct, guid string) (int, error) {
	for _, reserveDto := range reserveDtos {
		if reserveDto.Guid == guid {
			return reserveDto.Quantity, nil
		}
	}

	return 0, errors.New("ошибка поиска Quantity")
}
