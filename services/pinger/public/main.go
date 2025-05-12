package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kvvvseins/mictoservices/services/pinger/config"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/handler"
	"github.com/kvvvseins/mictoservices/services/pinger/internal/app/repository"
	"github.com/kvvvseins/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
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

	metricMiddleware := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			DurationBuckets: []float64{.5, 1, 2.5, 5, 10, 20, 40, 80, 160, 320, 640, 1000},
		}),
		Service: cnf.App.Name,
	})
	router, httpServer := server.NewServer(cnf.HTTP.Port, cnf.App.LogLevel)
	router.Use(std.HandlerProvider("/health/", metricMiddleware))
	router.Use(std.HandlerProvider("/user/{:id}/", metricMiddleware))

	profileRepository := repository.NewProfileRepository(cnf.GetDb())

	router.Method(
		http.MethodGet,
		"/health/",
		handler.NewHelloHandler(&cnf),
	)

	profileHandler := handler.NewProfileHandler(&cnf, profileRepository)

	handler.RegisterProfileHandlers(router, profileHandler)

	go func() {
		err = http.ListenAndServe(":9070", promhttp.Handler())
		if err != nil {
			slog.Error("pinger metric closed", slog.Int("port", 9070))
		}
	}()

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
