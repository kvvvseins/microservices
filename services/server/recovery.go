package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/handlers"
)

// Recovery middleware для восстановления после паники в http handler.
func Recovery(handler http.Handler) http.Handler {
	return handlers.RecoveryHandler(handlers.RecoveryLogger(recoveryLogger{}))(handler)
}

type recoveryLogger struct{}

func (r recoveryLogger) Println(i ...interface{}) {
	var msg string

	for _, val := range i {
		switch err := val.(type) {
		case error:
			msg += err.Error()
		case string:
			msg += err
		default:
			slog.Any("", err)
		}
	}

	slog.Error(msg)
}
