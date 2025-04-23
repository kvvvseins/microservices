package handler

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kvvvseins/mictoservices/services/pinger/config"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/dto"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/model"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/repository"
)

// ProfileHandler хендлер создания превью.
type ProfileHandler struct {
	config     *config.Config
	repository *repository.ProfileRepository
}

// NewProfileHandler создает хендлер crud для pinger.
func NewProfileHandler(
	cfg *config.Config,
	repository *repository.ProfileRepository,
) http.Handler {
	return &ProfileHandler{
		config:     cfg,
		repository: repository,
	}
}

// RegisterProfileHandlers регистрирует crud ручки профиля
func RegisterProfileHandlers(router *chi.Mux, handler http.Handler) {
	router.Method(
		http.MethodGet,
		"/profile/",
		handler,
	)
	router.Method(
		http.MethodPost,
		"/profile/",
		handler,
	)
	router.Method(
		http.MethodPut,
		"/profile/",
		handler,
	)
}

func (cu *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	guid, isOk := getGuidFromRequest(w, r)
	if !isOk {
		return
	}

	// todo вынести в сервис
	switch r.Method {
	case http.MethodGet:
		cu.get(w, r, guid)
	case http.MethodPost:
		cu.create(w, r, guid)
	case http.MethodPut:
		cu.update(w, r, guid)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (cu *ProfileHandler) get(w http.ResponseWriter, r *http.Request, guid uuid.UUID) {
	profile, err := cu.repository.FindByGuid(guid)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не удалось получить пользователя")

		return
	}

	cu.responseProfile(profile, w, r, http.StatusOK)
}

func (cu *ProfileHandler) create(w http.ResponseWriter, r *http.Request, guid uuid.UUID) {
	var userDto dto.CreateProfile
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не верные json нового юзера")

		return
	}

	profile := &model.Profile{
		Name: userDto.Name,
		Guid: &guid,
	}

	result := cu.config.GetDb().FirstOrCreate(profile, "guid = ?", guid.String())
	if result.Error != nil {
		textErrorResponse(r.Context(), w, result.Error, "ошибка создания юзера")

		return
	}

	cu.responseProfile(profile, w, r, http.StatusCreated)
}

// todo прикрутить валидатор
func (cu *ProfileHandler) validate(userDto dto.UpdateProfiel) error {
	if userDto.Name == "" {
		return errors.New("необходимо указать имя")
	}

	return nil
}

func (cu *ProfileHandler) update(w http.ResponseWriter, r *http.Request, guid uuid.UUID) {
	profile, err := cu.repository.FindByGuid(guid)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не удалось получить профиль пользователя при обновлении")

		return
	}

	var userDto dto.UpdateProfiel
	err = json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не верные json для обновленных данных")

		return
	}

	if err = cu.validate(userDto); err != nil {
		textErrorResponse(r.Context(), w, nil, err.Error())

		return
	}

	if userDto.Name == profile.Name {
		cu.responseProfile(profile, w, r, http.StatusOK)

		return
	}

	result := cu.config.GetDb().Model(profile).Updates(&userDto)
	if result.Error != nil {
		textErrorResponse(r.Context(), w, result.Error, "не удалось обновить пользователя")

		return
	}

	cu.responseProfile(profile, w, r, http.StatusOK)
}

func (cu *ProfileHandler) delete(w http.ResponseWriter, r *http.Request, guid uuid.UUID) {
	profile, err := cu.repository.FindByGuid(guid)
	if err != nil {
		textErrorResponse(r.Context(), w, err, "не удалось получить профиль пользователя при удалении")

		return
	}

	result := cu.config.GetDb().Delete(profile)
	if result.Error != nil {
		textErrorResponse(r.Context(), w, result.Error, "не удалось удалить профиль пользователя")

		return
	}

	w.WriteHeader(204)
}

func (cu *ProfileHandler) responseProfile(profile *model.Profile, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(dto.ViewProfile{Name: profile.Name})
	if err != nil {
		textErrorResponse(r.Context(), w, err, "ошибка обновления профиля")

		return
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
