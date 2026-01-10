package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (mw *MiddlewareCustom) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// before
		start := time.Now()

		next.ServeHTTP(w, r)

		// after
		duration := time.Since(start)
		mw.Log.Info("Activity route", zap.String("Method", r.Method), zap.String("URL", r.URL.String()), zap.Duration("Duration", duration))
	})
}
