package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kvvvseins/mictoservices/config"
	"github.com/kvvvseins/mictoservices/server"
)

// CrudPinger хендлер создания превью.
type CrudPinger struct {
	config *config.Config
}

// CrudPingerHandler создает хендлер crud для pinger.
func CrudPingerHandler(
	cfg *config.Config,
) http.Handler {
	return &Hello{
		config: cfg,
	}
}

func (rh *CrudPinger) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rh *CrudPinger) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (rh *CrudPinger) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа crud handler get", "err", err)

		w.WriteHeader(http.StatusNoContent)
	}
}

func (rh *CrudPinger) deleye(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}
