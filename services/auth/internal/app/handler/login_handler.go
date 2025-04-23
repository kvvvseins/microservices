package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kvvvseins/mictoservices/services/auth/config"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// LoginHandler структура хендлер логирования
type LoginHandler struct {
	config         *config.Config
	userRepository *repository.UserRepository
}

// NewLoginHandler создает хендлер login.
func NewLoginHandler(
	cfg *config.Config,
	userRepository *repository.UserRepository,
) http.Handler {
	return &LoginHandler{
		config:         cfg,
		userRepository: userRepository,
	}
}

func (cu *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginDto dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginDto)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не верные json логирования")

		return
	}

	// todo вынести в сервис
	user, err := cu.userRepository.FindByEmail(loginDto.Email, true)
	if err != nil {
		mess := "не удалось получить пользователя"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mess = "авторизация не удалась"
		}

		textErrorResponse(r.Context(), w, err, mess)

		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password)); err != nil {
		textErrorResponse(r.Context(), w, err, "не верный пароль")

		return
	}

	jwtResponse(w, r, cu.config, user)
}
