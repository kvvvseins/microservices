package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kvvvseins/mictoservices/services/notify/config"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/handler"
	"github.com/kvvvseins/mictoservices/services/notify/internal/app/repository"
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

	slog.SetDefault(logger)

	router, httpServer := server.NewServer(cnf.HTTP.Port, cnf.App.LogLevel)

	notifyRepository := repository.NewNotifyRepository(cnf.GetDb())

	router.Method(
		http.MethodGet,
		"/health/",
		handler.NewHelloHandler(&cnf),
	)

	notifyHandler := handler.NewNotifyHandler(&cnf, notifyRepository)

	handler.RegisterNotifyHandlers(router, notifyHandler)

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
