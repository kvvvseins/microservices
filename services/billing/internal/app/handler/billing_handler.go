package handler

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

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
	httpClient *http.Client
}

// NewBillingHandler создает хендлер crud для billing.
func NewBillingHandler(
	cfg *config.Config,
	repository *repository.BillingRepository,
	httpClient *http.Client,
) http.Handler {
	return &BillingHandler{
		config:     cfg,
		repository: repository,
		httpClient: httpClient,
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
	billing, err := cu.getOrCreateAccount(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить счет пользователя")

		return
	}

	cu.responseBilling(billing, w, r, http.StatusOK)
}

func (cu *BillingHandler) create(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing, err := cu.getOrCreateAccount(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка создания счета пользователя")

		return
	}

	cu.responseBilling(billing, w, r, http.StatusCreated)
}

// todo прикрутить валидатор
func (cu *BillingHandler) validate(updateBillingDto dto.UpdateBilling) error {
	if updateBillingDto.Value == 0 {
		return errors.New("нельзя указывать 0")
	}

	return nil
}

func (cu *BillingHandler) update(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	billing, err := cu.getOrCreateAccount(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить счет пользователя при обновлении")

		return
	}

	var updateBillingDto dto.UpdateBilling
	err = json.NewDecoder(r.Body).Decode(&updateBillingDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json для обновленных данных")

		return
	}

	if err = cu.validate(updateBillingDto); err != nil {
		server.ErrorResponseOutput(r.Context(), w, nil, err.Error())

		return
	}

	intVal := int(billing.Value) + updateBillingDto.Value
	if intVal < 0 {
		server.ErrorResponseOutput(r.Context(), w, nil, "недостаточно средств")

		return
	}

	billing.Value = uint(intVal)

	result := cu.config.GetDb().Model(billing).Update("value", gorm.Expr("value + ?", updateBillingDto.Value))
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "не удалось обновить счет пользователя")

		return
	}

	go cu.notifyPay(r.Context(), updateBillingDto, userID)

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

func (cu *BillingHandler) notifyPay(
	ctx context.Context,
	updateBillingDto dto.UpdateBilling,
	userID uuid.UUID,
) {
	valueStr := strconv.Itoa(updateBillingDto.Value)
	reader := strings.NewReader(`
{
    "email": "` + updateBillingDto.Email + `",
    "type": "change_billing",
    "data": {
        "price": ` + valueStr + `
    }
}`)
	notifySchema := cu.config.App.MicroservicesRoutes.Notify.Schema
	notifyRoute := cu.config.App.MicroservicesRoutes.Notify.Route
	notifyPort := cu.config.App.MicroservicesRoutes.Notify.Port
	request, errReq := http.NewRequest(http.MethodPost, notifySchema+"://"+notifyRoute+":"+notifyPort+"/", reader)
	if errReq != nil {
		server.GetLogger(ctx).Warn("не удалось создать request изменения по счету")

		return
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, errReq = cu.httpClient.Do(request)
	if errReq != nil {
		server.GetLogger(ctx).Warn("не удалось сделать запрос на отправку изменения по счету")

		return
	}

	if response.StatusCode != http.StatusCreated {
		server.GetLogger(ctx).Warn("не удалось отправить уведомление")
	}
}

func (cu *BillingHandler) getOrCreateAccount(userID uuid.UUID) (*model.Billing, error) {
	billing, err := cu.repository.FindByUserID(userID)
	if err == nil {
		return billing, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		billing = &model.Billing{
			UserID:   userID,
			Value:    0,
			Currency: "rub",
		}

		result := cu.config.GetDb().Create(billing)
		if result.Error != nil {
			return nil, result.Error
		}

		return billing, nil
	}

	return nil, err
}
