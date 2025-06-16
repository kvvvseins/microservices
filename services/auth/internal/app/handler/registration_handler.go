package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/auth/config"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// RegistrationHandler структура хендлер логирования
type RegistrationHandler struct {
	config         *config.Config
	userRepository *repository.UserRepository
	httpClient     *http.Client
}

// NewRegistrationHandler создает хендлер RegistrationHandler.
func NewRegistrationHandler(
	cfg *config.Config,
	userRepository *repository.UserRepository,
	httpClient *http.Client,
) http.Handler {
	return &RegistrationHandler{
		config:         cfg,
		userRepository: userRepository,
		httpClient:     httpClient,
	}
}

func (cu *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginDto dto.RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil || loginDto.Email == "" || loginDto.Password == "" {
		server.ErrorResponseOutput(r.Context(), w, err, "не верные json регистрации")

		return
	}

	user, err := cu.userRepository.FindByEmail(loginDto.Email, false)
	if err == nil {
		server.ErrorResponseOutput(r.Context(), w, nil, "данный пользователь уже существует")

		return
	}

	user, err = cu.createUser(r.Context(), loginDto)
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "транзакция создания юзера прервана")

		return
	}

	go cu.createBillingAccount(r.Context(), user.Guid)

	jwtResponse(w, r, cu.config, user)
}

// todo вынести в сервис
func (cu *RegistrationHandler) createUser(ctx context.Context, loginDto dto.RegistrationRequest) (*model.User, error) {
	user := &model.User{
		Email:    loginDto.Email,
		Password: loginDto.Password,
		IsActive: true,
	}

	err := cu.config.GetDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(user)
		if result.Error != nil {
			return result.Error
		}

		reader := strings.NewReader(`{"name": "add name"}`)
		pingerSchema := cu.config.App.MicroservicesRoutes.Pinger.Schema
		pingerRoute := cu.config.App.MicroservicesRoutes.Pinger.Route
		pingerPort := cu.config.App.MicroservicesRoutes.Pinger.Port
		request, errReq := http.NewRequest(http.MethodPost, pingerSchema+"://"+pingerRoute+":"+pingerPort+"/profile/", reader)
		if errReq != nil {
			return errReq
		}

		server.SetUserIDToHeader(request.Header, user.Guid)
		server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

		var response *http.Response

		response, errReq = cu.httpClient.Do(request)
		if errReq != nil {
			return errors.Wrap(errReq, "не удалось сделать запрос на создание профиля пользователя")
		}

		if response.StatusCode != http.StatusCreated {
			return errors.New("не удалось создать профиль пользователя")
		}

		return nil
	})

	return user, err
}

func (cu *RegistrationHandler) createBillingAccount(ctx context.Context, userID uuid.UUID) {
	billingSchema := cu.config.App.MicroservicesRoutes.Billing.Schema
	billingRoute := cu.config.App.MicroservicesRoutes.Billing.Route
	billingPort := cu.config.App.MicroservicesRoutes.Billing.Port
	errMsgBase := "Не удалось создать платежный аккаунт"
	request, errReq := http.NewRequest(http.MethodPost, billingSchema+"://"+billingRoute+":"+billingPort+"/", nil)
	if errReq != nil {
		server.GetLogger(ctx).Warn(errMsgBase, "msg", errReq.Error())

		return
	}

	server.SetUserIDToHeader(request.Header, userID)
	server.AddRequestIDToRequestHeader(request.Header, server.GetRequestID(ctx))

	var response *http.Response

	response, errReq = cu.httpClient.Do(request)
	if errReq != nil {
		server.GetLogger(ctx).Warn(errMsgBase, "msg", errReq.Error())

		return
	}

	if response.StatusCode != http.StatusCreated {
		server.GetLogger(ctx).Warn(errMsgBase, "msg", "status not created")
	}
}
