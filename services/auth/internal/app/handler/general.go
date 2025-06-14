package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kvvvseins/mictoservices/services/auth/config"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/model"
	"github.com/kvvvseins/server"
)

const userGuidHeader = "X-User-Id"

func jwtResponse(
	w http.ResponseWriter,
	r *http.Request,
	cnf *config.Config,
	user *model.User,
) {
	// Создаем JWT токен
	claims := jwt.MapClaims{
		"sub": user.Guid,
		"iss": cnf.App.Jwt.Issuer,
		"aud": cnf.App.Jwt.Audience,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cnf.App.Jwt.Secret))
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(dto.LoginResponse{Token: tokenString})
	if err != nil {
		server.ErrorResponseOutput(r.Context(), w, err, "ошибка логирования")

		return
	}
}
