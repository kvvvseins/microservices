package handler

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/delivery/config"
	"github.com/kvvvseins/mictoservices/services/delivery/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/delivery/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/delivery/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var layoutDate = "2006-01-02 15:04:05-07:00"

// DeliveryHandler хендлер создания delivery.
type DeliveryHandler struct {
	config     *config.Config
	repository *repository.DeliveryRepository
}

// NewDeliveryHandler создает хендлер crud для delivery.
func NewDeliveryHandler(
	cfg *config.Config,
	repository *repository.DeliveryRepository,
) http.Handler {
	return &DeliveryHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterDeliveryHandlers регистрирует crud ручки доставки
func RegisterDeliveryHandlers(router *chi.Mux, handler http.Handler) {
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

func (cu *DeliveryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, isOk := server.GetUserIDFromRequest(w, r)
	if !isOk {
		return
	}

	// todo вынести в сервис
	switch r.Method {
	case http.MethodGet:
		cu.get(w, r)
	case http.MethodPost:
		cu.create(w, r, userID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *DeliveryHandler) get(w http.ResponseWriter, r *http.Request) {
	deliveries, err := cu.repository.FindAll()
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить доставку заказа")

		return
	}

	w.WriteHeader(http.StatusOK)

	deliveryDtos := make([]dto.ViewDelivery, 0, len(deliveries))
	for _, delivery := range deliveries {
		deliveryDtos = append(deliveryDtos, dto.ViewDelivery{
			UserId:           delivery.UserID,
			OrderId:          delivery.OrderId,
			PlannedDateStart: delivery.PlannedDateStart.Format(layoutDate),
			PlannedDateEnd:   delivery.PlannedDateEnd.Format(layoutDate),
		})
	}

	err = json.NewEncoder(w).Encode(deliveryDtos)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения доставки")

		return
	}
}

func (cu *DeliveryHandler) create(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var deliveryDto dto.CreateDelivery
	err := json.NewDecoder(r.Body).Decode(&deliveryDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json нового доставки заказа")

		return
	}

	delivery, err := cu.repository.FindByOrderID(deliveryDto.OrderId)
	if err == nil {
		cu.responseDelivery(delivery, w, r, http.StatusCreated)

		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка проверки существующей доставки")

		return
	}

	plannedDateStart, err := time.Parse(layoutDate, deliveryDto.PlannedDateStart)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, nil, "wrong plannedDateStart")

		return
	}

	plannedDateEnd, err := time.Parse(layoutDate, deliveryDto.PlannedDateEnd)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, nil, "wrong plannedDateEnd")

		return
	}

	if plannedDateStart == plannedDateEnd {
		server.ErrorResponseOutput(r.Context(), w, nil, "даты не могут быть равны")

		return
	}

	if plannedDateStart.After(plannedDateEnd) {
		server.ErrorResponseOutput(r.Context(), w, nil, "дата тачала должна быть меньше")

		return
	}

	delivery = &model.Delivery{
		OrderId:          deliveryDto.OrderId,
		UserID:           userID,
		PlannedDateStart: plannedDateStart,
		PlannedDateEnd:   plannedDateEnd,
	}

	result := cu.config.GetDb().Create(delivery)
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "ошибка создания доставки заказа")

		return
	}

	cu.responseDelivery(delivery, w, r, http.StatusCreated)
}

func (cu *DeliveryHandler) responseDelivery(delivery *model.Delivery, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewDelivery{
		UserId:           delivery.UserID,
		OrderId:          delivery.OrderId,
		PlannedDateStart: delivery.PlannedDateStart.Format(layoutDate),
		PlannedDateEnd:   delivery.PlannedDateEnd.Format(layoutDate),
	})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка доставки заказа")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
