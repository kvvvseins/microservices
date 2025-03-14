package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

// NewServer создает новый инстанс роутера и сервера.
func NewServer(port int, logLevel slog.Level) (*chi.Mux, *http.Server) {
	router := chi.NewRouter()

	router.Use(
		Recovery,
		LoggingHandler(logLevel),
	)

	accessLog := AccessLog(zerolog.New(os.Stdout))

	server := &http.Server{
		Handler:           accessLog(router),
		Addr:              fmt.Sprintf(":%d", port),
		ReadTimeout:       time.Millisecond * 250,
		ReadHeaderTimeout: time.Millisecond * 200,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Minute * 30,
	}

	return router, server
}
