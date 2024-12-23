package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
)

// LoggingHandler добавляет во все запросы свой логер.
func LoggingHandler(logLevel slog.Level) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			})).With(slog.String("request_id", GetRequestID(r.Context())))

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), keyLogger, logger)))
		})
	}
}
