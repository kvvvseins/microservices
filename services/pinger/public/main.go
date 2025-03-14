package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kvvvseins/mictoservices/services/pinger/config/pinger"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/handler"
	"github.com/kvvvseins/mictoservices/services/pinger/server"
)

func main() {
	cnf := pinger.Config{}
	prefixCfg := ""

	err := pinger.LoadConfig(&cnf, prefixCfg)
	if err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))

	slog.SetDefault(logger)

	router, httpServer := server.NewServer(cnf.HTTP.Port, cnf.App.LogLevel)

	router.Method(
		"GET",
		"/health/",
		handler.NewHelloHandler(&cnf),
	)

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
