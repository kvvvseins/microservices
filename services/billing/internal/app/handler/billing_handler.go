package handler

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/billing/config"
	"github.com/kvvvseins/mictoservices/services/billing/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/billing/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/billing/internal/app/repository"
	"github.com/kvvvseins/server"
	"gorm.io/gorm"
)

// BillingHandler хендлер создания billing.
type BillingHandler struct {
	config     *config.Config
	repository *repository.BillingRepository
}

// NewBillingHandler создает хендлер crud для billing.
func NewBillingHandler(
	cfg *config.Config,
	repository *repository.BillingRepository,
) http.Handler {
	return &BillingHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterBillingHandlers регистрирует crud ручки счета
func RegisterBillingHandlers(router *chi.Mux, handler http.Handler) {
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
	router.Method(
		http.MethodPut,
		"/",
		handler,
	)
}

func (cu *BillingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, isOk := server.GetUserIDFromRequest(w, r)
	if !isOk {
		return
	}

	// todo вынести в сервис
	switch r.Method {
	case http.MethodGet:
		cu.get(w, r, userID)
	case http.MethodPost:
		cu.create(w, r, userID)
	case http.MethodPut:
		cu.update(w, r, userID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *BillingHandler) get(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing, err := cu.repository.FindByUserID(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить счет пользователя")

		return
	}

	cu.responseBilling(billing, w, r, http.StatusOK)
}

func (cu *BillingHandler) create(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing := &model.Billing{
		UserID:   userID,
		Value:    0,
		Currency: "rub",
	}

	result := cu.config.GetDb().FirstOrCreate(billing, "user_id = ?", userID.String())
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "ошибка создания счета пользователя")

		return
	}

	cu.responseBilling(billing, w, r, http.StatusCreated)
}

// todo прикрутить валидатор
func (cu *BillingHandler) validate(userDto dto.UpdateBilling) error {
	if userDto.Value == 0 {
		return errors.New("нельзя указывать 0")
	}

	return nil
}

func (cu *BillingHandler) update(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing, err := cu.repository.FindByUserID(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить счет пользователя при обновлении")

		return
	}

	var userDto dto.UpdateBilling
	err = json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json для обновленных данных")

		return
	}

	if err = cu.validate(userDto); err != nil {
		server.ErrorResponseOutput(r.Context(), w, nil, err.Error())

		return
	}

	intVal := int(billing.Value) + userDto.Value
	if intVal < 0 {
		server.ErrorResponseOutput(r.Context(), w, nil, "недостаточно средств")

		return
	}

	billing.Value = uint(intVal)

	result := cu.config.GetDb().Model(billing).Update("value", gorm.Expr("value + ?", userDto.Value))
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "не удалось обновить счет пользователя")

		return
	}

	cu.responseBilling(billing, w, r, http.StatusOK)
}

func (cu *BillingHandler) delete(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing, err := cu.repository.FindByUserID(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить счет пользователя при удалении")

		return
	}

	result := cu.config.GetDb().Delete(billing)
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "не удалось удалить счет пользователя")

		return
	}

	w.WriteHeader(204)
}

func (cu *BillingHandler) responseBilling(billing *model.Billing, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewBilling{Value: billing.Value})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка счета пользователя")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
