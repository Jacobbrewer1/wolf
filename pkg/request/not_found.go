package request

import (
	"encoding/json"
	"net/http"

	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"golang.org/x/exp/slog"
)

// NotFoundHandler returns a handler that returns a 404 response.
func NotFoundHandler(l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := NewMessage("Not found")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			l.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
		}
	}
}

// MethodNotAllowedHandler returns a handler that returns a 405 response.
func MethodNotAllowedHandler(l *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := NewMessage("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			l.Error("Error encoding response", slog.String(logging.KeyError, err.Error()))
		}
	}
}
