package handler

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kvvvseins/mictoservices/services/pinger/config/pinger"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/pinger/server"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// ErrorResponse структура ошибки от CrudUser.
type ErrorResponse struct {
	Message string `json:"message"`
}

// CrudUser хендлер создания превью.
type CrudUser struct {
	config *pinger.Config
}

// CrudUserHandler создает хендлер crud для pinger.
func CrudUserHandler(
	cfg *pinger.Config,
) http.Handler {
	return &CrudUser{
		config: cfg,
	}
}

func (cu *CrudUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		cu.get(w, r)
	case http.MethodPost:
		cu.create(w, r)
	case http.MethodPut:
		cu.update(w, r)
	case http.MethodDelete:
		cu.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *CrudUser) get(w http.ResponseWriter, r *http.Request) {
	randomI := randRange(20, 1000)
	time.Sleep(time.Duration(randomI) * time.Millisecond)

	user, err := cu.getUserByRequest(r)
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "не удалось получить пользователя")

		return
	}

	cu.responseUser(user, w, r, http.StatusOK)
}

func (cu *CrudUser) create(w http.ResponseWriter, r *http.Request) {
	var userDto dto.CreateUser
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "не верные json нового юзера")

		return
	}

	user := &model.User{
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	result := cu.config.GetDb().FirstOrCreate(user, "email = ?", userDto.Email)
	if result.Error != nil {
		cu.textErrorResponse(r.Context(), w, result.Error, "ошибка создания юзера")

		return
	}

	cu.responseUser(user, w, r, http.StatusCreated)
}

func (cu *CrudUser) update(w http.ResponseWriter, r *http.Request) {
	user, err := cu.getUserByRequest(r)
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "не удалось получить пользователя при обновлении")

		return
	}

	var userDto dto.UpdateUser
	err = json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "не верные json для обновленных данных")

		return
	}

	if userDto.Email == user.Email {
		cu.responseUser(user, w, r, http.StatusOK)

		return
	}

	var isExists bool
	err = cu.config.GetDb().Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?) AS isExists",
		userDto.Email).Scan(&isExists).Error
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "ошибка проверки существующего пользователя")

		return
	}

	if isExists {
		cu.textErrorResponse(r.Context(), w, nil, "пользователь с такой почтой уже существует")

		return
	}

	result := cu.config.GetDb().Model(user).Updates(&userDto)
	if result.Error != nil {
		cu.textErrorResponse(r.Context(), w, result.Error, "не удалось обновить пользователя")

		return
	}

	cu.responseUser(user, w, r, http.StatusOK)
}

func (cu *CrudUser) delete(w http.ResponseWriter, r *http.Request) {
	user, err := cu.getUserByRequest(r)
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "не удалось получить пользователя при удалении")

		return
	}

	result := cu.config.GetDb().Delete(user)
	if result.Error != nil {
		cu.textErrorResponse(r.Context(), w, result.Error, "не удалось удалить пользователя")

		return
	}

	w.WriteHeader(204)
}

func (cu *CrudUser) textErrorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
	errMsg string,
) {
	if err != nil {
		server.GetLogger(ctx).Error("ошибка user crud", "err", err)
	}

	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: errMsg})
}

func (cu *CrudUser) getUserByRequest(r *http.Request) (*model.User, error) {
	userId, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, strconv.IntSize)
	if err != nil {
		return nil, errors.Wrap(err, "не верны id юзера")
	}

	var user model.User

	result := cu.config.GetDb().First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	}

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, "ошибка получения пользователя")
	}

	return &user, nil
}

func (cu *CrudUser) responseUser(user *model.User, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewUser{Email: user.Email})
	if err != nil {
		cu.textErrorResponse(r.Context(), w, err, "ошибка обновления пользователя")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
