package handler

import (
	"context"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"slices"
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
		ordersDto = append(ordersDto, dto.ViewOrder{
			Number: order.Guid.String(),
			Price:  order.Price,
		})
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

	if len(orderDto.Products) == 0 {
		server.ErrorResponseOutput(r.Context(), w, nil, "не указаны товары")

		return
	}

	price, err := cu.calculateOrderPrice(r.Context(), orderDto.Products, userID)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка расчета цены заказа")

		return
	}

	if price <= 0 {
		server.ErrorResponseOutput(r.Context(), w, err, "не верная цена заказа")

		return
	}

	if err = cu.writeOffMoney(r.Context(), userID, int(price), false); err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка списания средств, возможно недостаточно денег на счете")

		return
	}

	orderGuid := uuid.New()

	reserveProducts := slices.Clone(orderDto.Products)
	for i := range reserveProducts {
		reserveProducts[i].Quantity = reserveProducts[i].Quantity * -1
	}

	err = cu.reserveProducts(r.Context(), reserveProducts, userID, orderGuid)
	if err != nil {
		go cu.backWriteOffMoney(r.Context(), price, userID)

		server.ErrorResponseOutput(r.Context(), w, err, "недостаточно товара")

		return
	}

	order := &model.Order{
		UserID: userID,
		Price:  price,
		Guid:   orderGuid,
	}

	result := cu.config.GetDb().Create(order)
	if result.Error != nil {
		go cu.backReserveProducts(r.Context(), orderDto.Products, userID, orderGuid)
		go cu.backWriteOffMoney(r.Context(), price, userID)

		server.ErrorResponseOutput(r.Context(), w, result.Error, "ошибка создания заказа пользователя")

		return
	}

	go cu.notifyCreateOrder(r.Context(), order, orderDto.Email, userID)

	cu.responseOrder(order, w, r, http.StatusCreated)
}

func (cu *OrderHandler) backWriteOffMoney(ctx context.Context, price uint, userID uuid.UUID) {
	go func() {
		err := cu.writeOffMoney(ctx, userID, int(price), true)
		if err != nil {
			//@todo сделать retry или через кафку
		}
	}()
}

func (cu *OrderHandler) backReserveProducts(ctx context.Context, products []dto.ProductDto, userID, orderGuid uuid.UUID) {
	go func() {
		err := cu.reserveProducts(ctx, products, userID, orderGuid)
		if err != nil {
			//@todo сделать retry или через кафку
		}
	}()
}

func (cu *OrderHandler) responseOrder(order *model.Order, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewOrder{
		Number: order.Guid.String(),
		Price:  order.Price,
	})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка заказа пользователя")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}

func (cu *OrderHandler) writeOffMoney(ctx context.Context, userID uuid.UUID, price int, isBack bool) error {
	priceStr := strconv.Itoa(price)

	since := "-"
	if isBack {
		since = ""
	}

	reader := strings.NewReader(`{"value": ` + since + priceStr + `}`)
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

func (cu *OrderHandler) notifyCreateOrder(
	ctx context.Context,
	order *model.Order,
	email string,
	userID uuid.UUID,
) {
	priceStr := strconv.Itoa(int(order.Price))
	reader := strings.NewReader(`
{
    "email": "` + email + `",
    "type": "create_order",
    "data": {
        "price": ` + priceStr + `
    }
}`)
	notifySchema := cu.config.App.MicroservicesRoutes.Notify.Schema
	notifyRoute := cu.config.App.MicroservicesRoutes.Notify.Route
	notifyPort := cu.config.App.MicroservicesRoutes.Notify.Port
	request, errReq := http.NewRequest(http.MethodPost, notifySchema+"://"+notifyRoute+":"+notifyPort+"/", reader)
	if errReq != nil {
		server.GetLogger(ctx).Warn("не удалось создать request для Notify", "msg", errReq.Error())

		return
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, errReq = cu.httpClient.Do(request)
	if errReq != nil {
		server.GetLogger(ctx).Warn("не удалось сделать запрос на отправку уведомления", "msg", errReq.Error())

		return
	}

	if response.StatusCode != http.StatusCreated {
		server.GetLogger(ctx).Warn("не удалось отправить уведомление")
	}
}

func (cu *OrderHandler) calculateOrderPrice(
	ctx context.Context,
	products []dto.ProductDto,
	userID uuid.UUID,
) (uint, error) {
	queryGuids := "?"
	for _, product := range products {
		queryGuids += "guid=" + product.Guid.String() + "&"
	}

	//@todo отдельный мс цен
	storeSchema := cu.config.App.MicroservicesRoutes.Store.Schema
	storeRoute := cu.config.App.MicroservicesRoutes.Store.Route
	storePort := cu.config.App.MicroservicesRoutes.Store.Port
	request, err := http.NewRequest(http.MethodGet, storeSchema+"://"+storeRoute+":"+storePort+"/"+queryGuids, nil)
	if err != nil {
		return 0, errors.Wrap(err, "не удалось создать request для получения цен")
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, err = cu.httpClient.Do(request)
	if err != nil {
		return 0, errors.Wrap(err, "не удалось сделать запрос на получения цен")
	}

	if response.StatusCode != http.StatusOK {
		return 0, errors.New("не удалось получить цены")
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, errors.Wrap(err, "Ошибка при чтении тела ответа получения цен")
	}

	var priceDtos []dto.Price
	err = json.Unmarshal(body, &priceDtos)
	if err != nil {
		return 0, errors.Wrap(err, "Ошибка парсинга цен")
	}

	var prices uint

	for _, price := range priceDtos {
		quantity, errQ := findQuantity(products, price.Guid)
		if errQ != nil {
			return 0, errQ
		}

		prices = prices + (price.Price * quantity)
	}

	return prices, nil
}

func findQuantity(products []dto.ProductDto, guid uuid.UUID) (uint, error) {
	for _, product := range products {
		if product.Guid == guid {
			if product.Quantity <= 0 {
				return 0, errors.New("количество товара не может быть <= 0")
			}

			return uint(product.Quantity), nil
		}
	}

	return 0, errors.New("ошибка поиска Quantity")
}

func (cu *OrderHandler) reserveProducts(
	ctx context.Context,
	products []dto.ProductDto,
	userID uuid.UUID,
	orderGuid uuid.UUID,
) error {
	jsonBody, err := json.Marshal(products)
	if err != nil {
		return errors.Wrap(err, "не удалось получить json")
	}

	reader := strings.NewReader(`
{
    "order_id": "` + orderGuid.String() + `",
    "products": ` + string(jsonBody) + `
}`)

	//@todo отдельный мс цен
	storeSchema := cu.config.App.MicroservicesRoutes.Store.Schema
	storeRoute := cu.config.App.MicroservicesRoutes.Store.Route
	storePort := cu.config.App.MicroservicesRoutes.Store.Port
	request, err := http.NewRequest(http.MethodPost, storeSchema+"://"+storeRoute+":"+storePort+"/reserve", reader)
	if err != nil {
		return errors.Wrap(err, "не удалось создать request для резервирования")
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, err = cu.httpClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "не удалось сделать запрос на резервирование")
	}

	if response.StatusCode != http.StatusCreated {
		return errors.New("не удалось зарезервировать")
	}

	return nil
}
