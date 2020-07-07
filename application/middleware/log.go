package middleware

import (
	"net/http"
)

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("REQUEST", "method", r.Method, "path", r.URL.Path, "query", r.URL.Query(),
			"headers", r.Header)
		next.ServeHTTP(w, r)
	})
}
