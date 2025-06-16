package handler

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/auth/config"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/repository"
	"github.com/kvvvseins/server"
)

// VerifyHandler структура хендлер проверки jwt
type VerifyHandler struct {
	config         *config.Config
	userRepository *repository.UserRepository
}

// NewVerifyHandler создает хендлер login.
func NewVerifyHandler(
	cfg *config.Config,
	userRepository *repository.UserRepository,
) http.Handler {
	return &VerifyHandler{
		config:         cfg,
		userRepository: userRepository,
	}
}

func (cu *VerifyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jwtToken := r.Header.Get("Authorization")
	if "" == jwtToken {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	splitToken := strings.Split(jwtToken, "Bearer ")
	if len(splitToken) != 2 {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	jwtToken = splitToken[1]

	// todo вынести в сервис
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cu.config.App.Jwt.Secret), nil
	},
		jwt.WithIssuer(cu.config.App.Jwt.Issuer),
		jwt.WithAudience(cu.config.App.Jwt.Audience),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	subject, err := claims.GetSubject()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	userID, err := uuid.Parse(subject)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	_, err = cu.userRepository.FindByGuid(userID, true)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	server.SetUserIDToHeader(w.Header(), userID)

	w.WriteHeader(http.StatusOK)
}
