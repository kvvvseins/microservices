package handler

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/notify/config"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/repository"
	"github.com/kvvvseins/server"
)

// NotifyHandler хендлер создания notify.
type NotifyHandler struct {
	config     *config.Config
	repository *repository.NotifyRepository
}

// NewNotifyHandler создает хендлер crud для notify.
func NewNotifyHandler(
	cfg *config.Config,
	repository *repository.NotifyRepository,
) http.Handler {
	return &NotifyHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterNotifyHandlers регистрирует crud ручки уведомления
func RegisterNotifyHandlers(router *chi.Mux, handler http.Handler) {
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

func (cu *NotifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *NotifyHandler) get(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	notifies, err := cu.repository.FindByUserID(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить уведомление пользователя")

		return
	}

	w.WriteHeader(http.StatusOK)

	notifiesDto := make([]dto.ViewNotify, 0, len(notifies))
	for _, notify := range notifies {
		notifiesDto = append(notifiesDto, dto.ViewNotify{
			Email:   notify.Email,
			Message: notify.Message,
		})
	}

	err = json.NewEncoder(w).Encode(notifiesDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения уведомлений пользователя")

		return
	}
}

func (cu *NotifyHandler) create(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var notifyDto dto.CreateNotify
	err := json.NewDecoder(r.Body).Decode(&notifyDto)
	if err != nil || "" == notifyDto.Email {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json нового профиля пользователя")

		return
	}

	var notify *model.Notify

	switch notifyDto.Type {
	case "create_order":
		notify = buildCreateOrderNotify(&notifyDto, userID)
	case "change_billing":
		notify = buildPayOrderNotify(&notifyDto, userID)
	}

	if notify == nil {
		server.ErrorResponseOutput(r.Context(), w, err, "неизвестный тип уведомлений или недостаточно данных")

		return
	}

	result := cu.config.GetDb().Create(notify)
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "ошибка создания уведомления пользователя")

		return
	}

	cu.responseNotify(notify, w, r, http.StatusCreated)
}

func buildCreateOrderNotify(notifyDto *dto.CreateNotify, userID uuid.UUID) *model.Notify {
	priceStr, ok := getPriceStrFromMap(notifyDto.Data)
	if !ok {
		return nil
	}

	return &model.Notify{
		UserID:  userID,
		Email:   notifyDto.Email,
		Message: "Заказ создан на сумму — " + priceStr,
	}
}

func buildPayOrderNotify(notifyDto *dto.CreateNotify, userID uuid.UUID) *model.Notify {
	priceStr, ok := getPriceStrFromMap(notifyDto.Data)
	if !ok {
		return nil
	}

	var currencyStr string
	currency, ok := notifyDto.Data["currency"]
	if !ok {
		currencyStr = "rub"
	} else {
		currencyStr = currency.(string)
	}

	var typeNotify string

	if strings.Contains(priceStr, "-") {
		typeNotify = "Списание"
		priceStr = strings.ReplaceAll(priceStr, "-", "")
	} else {
		typeNotify = "Пополнение"
	}

	return &model.Notify{
		UserID:  userID,
		Email:   notifyDto.Email,
		Message: typeNotify + " средств на сумму " + priceStr + " " + currencyStr,
	}
}

func getPriceStrFromMap(data map[string]interface{}) (string, bool) {
	price, ok := data["price"]
	if !ok {
		return "", false
	}

	priceFloat, ok := price.(float64)
	if !ok {
		return "", false
	}

	return strconv.FormatFloat(priceFloat, 'f', 0, 64), true
}

func (cu *NotifyHandler) responseNotify(notify *model.Notify, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewNotify{
		Email:   notify.Email,
		Message: notify.Message,
	})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка уведомления пользователя")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
