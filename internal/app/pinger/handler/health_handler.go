package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kvvvseins/mictoservices/config"
	"github.com/kvvvseins/mictoservices/server"
)

// Hello хендлер создания превью.
type Hello struct {
	config *config.Config
}

// HelloResponse ответ ручки
type HelloResponse struct {
	Status string `json:"status"`
}

// NewHelloHandler создает хендлер создания превью.
func NewHelloHandler(
	cfg *config.Config,
) http.Handler {
	return &Hello{
		config: cfg,
	}
}

func (rh *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK"})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа hello handler", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}
