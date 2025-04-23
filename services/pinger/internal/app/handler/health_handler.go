package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kvvvseins/mictoservices/services/pinger/config/pinger"
	"github.com/kvvvseins/mictoservices/services/pinger/server"
)

// Hello хендлер создания превью.
type Hello struct {
	config *pinger.Config
}

// HelloResponse ответ ручки
type HelloResponse struct {
	Status string `json:"status"`
	App    string `json:"app"`
}

// NewHelloHandler создает хендлер создания превью.
func NewHelloHandler(
	cfg *pinger.Config,
) http.Handler {
	return &Hello{
		config: cfg,
	}
}

func (rh *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	randomI := randRange(20, 1500)
	time.Sleep(time.Duration(randomI) * time.Millisecond)

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(HelloResponse{Status: "OK", App: rh.config.App.Name})
	if err != nil {
		server.GetLogger(r.Context()).Error("ошибка ответа hello handler", "err", err)

		w.WriteHeader(http.StatusInternalServerError)
	}
}
