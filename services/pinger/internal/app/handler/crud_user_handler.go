package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kvvvseins/mictoservices/services/pinger/config/pinger"
	"github.com/kvvvseins/mictoservices/services/pinger/server"
)

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

func (rh *CrudUser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		rh.get(w, r)
	case http.MethodPost:
		rh.create(w, r)
	case http.MethodPut:
		rh.update(w, r)
	case http.MethodDelete:
		rh.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rh *CrudUser) get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rh *CrudUser) create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rh *CrudUser) update(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusNoContent)
	}
}

func (rh *CrudUser) delete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
