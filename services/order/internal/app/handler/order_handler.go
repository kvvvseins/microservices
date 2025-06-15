package handler

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/order/config"
	"github.com/kvvvseins/mictoservices/services/order/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/order/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/order/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/pkg/errors"
)

// OrderHandler хендлер создания order.
type OrderHandler struct {
	config     *config.Config
	repository *repository.OrderRepository
	httpClient *http.Client
}

// NewOrderHandler создает хендлер crud для order.
func NewOrderHandler(
	cfg *config.Config,
	repository *repository.OrderRepository,
	httpClient *http.Client,
) http.Handler {
	return &OrderHandler{
		config:     cfg,
		repository: repository,
		httpClient: httpClient,
	}
}

// RegisterOrderHandlers регистрирует crud ручки заказа
func RegisterOrderHandlers(router *chi.Mux, handler http.Handler) {
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

func (cu *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (cu *OrderHandler) get(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	orders, err := cu.repository.FindByUserID(userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не удалось получить заказ пользователя")

		return
	}

	w.WriteHeader(http.StatusOK)

	ordersDto := make([]dto.ViewOrder, 0, len(orders))
	for _, order := range orders {
		ordersDto = append(ordersDto, dto.ViewOrder{Number: order.Guid.String()})
	}

	err = json.NewEncoder(w).Encode(ordersDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка получения заказа пользователя")

		return
	}
}

func (cu *OrderHandler) create(w http.ResponseWriter, r *http.Request, userID uuid.UUID) {
	var orderDto dto.CreateOrder
	err := json.NewDecoder(r.Body).Decode(&orderDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json нового заказа пользователя")

		return
	}

	var price int
	for _, productDto := range orderDto.Products {
		price += productDto.Price
	}

	if price <= 0 {
		server.ErrorResponseOutput(r.Context(), w, err, "не верная цена заказа")

		return
	}

	if err = cu.writeOffMoney(r.Context(), userID, price); err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка списания средств, возможно недостаточно денег на счете")

		return
	}

	order := &model.Order{
		UserID: userID,
	}

	result := cu.config.GetDb().Create(order)
	if result.Error != nil {
		server.ErrorResponseOutput(r.Context(), w, result.Error, "ошибка создания заказа пользователя")

		return
	}

	cu.responseOrder(order, w, r, http.StatusCreated)
}

func (cu *OrderHandler) responseOrder(order *model.Order, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewOrder{Number: order.Guid.String()})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка заказа пользователя")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (cu *OrderHandler) writeOffMoney(ctx context.Context, userID uuid.UUID, price int) error {
	priceStr := strconv.Itoa(price)
	reader := strings.NewReader(`{"value": -` + priceStr + `}`)
	billingSchema := cu.config.App.MicroservicesRoutes.Billing.Schema
	billingRoute := cu.config.App.MicroservicesRoutes.Billing.Route
	billingPort := cu.config.App.MicroservicesRoutes.Billing.Port
	request, errReq := http.NewRequest(http.MethodPut, billingSchema+"://"+billingRoute+":"+billingPort+"/", reader)
	if errReq != nil {
		return errReq
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, errReq = cu.httpClient.Do(request)
	if errReq != nil {
		return errors.Wrap(errReq, "не удалось сделать запрос на списание средств")
	}

	if response.StatusCode != http.StatusOK {
		return errors.New("не удалось списать средства")
	}

	return nil
}
