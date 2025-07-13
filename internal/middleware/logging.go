package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Debug("incoming request details",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("query", r.URL.RawQuery),
				slog.String("user_agent", r.UserAgent()),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("host", r.Host),
				slog.Any("headers", r.Header),
			)
			logger.Info("request started",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("user_agent", r.UserAgent()),
				slog.String("remote_addr", r.RemoteAddr),
			)

			next.ServeHTTP(w, r)
			duration := time.Since(start)

			logger.Debug("request performance",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", duration),
				slog.Int64("duration_ms", duration.Milliseconds()),
			)

			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", duration),
			)
		})
	}
}
