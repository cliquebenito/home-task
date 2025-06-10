package pkg

import (
	"log/slog"
	"net/http"
)

// MethodCheckMiddleware is a middleware that checks if the request method is the expected method
func MethodCheckMiddleware(expectedMethod string, logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			logger.Warn("Method Not Allowed",
				slog.String("method", r.Method),
				slog.String("expected", expectedMethod),
				slog.String("url", r.URL.Path),
			)
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
