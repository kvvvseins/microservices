package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kvvvseins/mictoservices/services/auth/config"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/handler"
	"github.com/kvvvseins/mictoservices/services/auth/internal/app/repository"
	"github.com/kvvvseins/server"
)

func main() {
	cnf := config.Config{}
	prefixCfg := ""

	err := config.LoadConfig(&cnf, prefixCfg)
	if err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cnf.App.LogLevel,
	}))

	httpClient := &http.Client{
		Timeout: time.Second * 2,
	}

	slog.SetDefault(logger)

	router, httpServer := server.NewServer(cnf.HTTP.Port, cnf.App.LogLevel)

	userRepository := repository.NewUserRepository(cnf.GetDb())

	router.Method(
		http.MethodGet,
		"/health/",
		handler.NewHelloHandler(&cnf),
	)

	router.Method(
		http.MethodPost,
		"/login/",
		handler.NewLoginHandler(&cnf, userRepository),
	)

	router.Method(
		http.MethodPost,
		"/registration/",
		handler.NewRegistrationHandler(&cnf, userRepository, httpClient),
	)

	router.Method(
		http.MethodGet,
		"/verify",
		handler.NewVerifyHandler(&cnf, userRepository),
	)

	router.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	if err = chi.Walk(router, walkRouters); err != nil {
		slog.Debug("Logging err", "error", err.Error())
	}

	slog.Debug("public server started", slog.Int("port", cnf.HTTP.Port))

	if err = httpServer.ListenAndServe(); err != nil {
		slog.Error("finished app with error", "error", err)
	}
}

func walkRouters(
	method string,
	route string,
	_ http.Handler,
	_ ...func(http.Handler) http.Handler,
) error {
	route = strings.ReplaceAll(route, "/*/", "/")

	slog.Debug(fmt.Sprintf("%s %s\n", method, route))

	return nil
}
