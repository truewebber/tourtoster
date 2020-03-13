package middleware

import (
	"net/http"

	"github.com/mgutz/logxi/v1"
)

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("REQUEST", "method", r.Method, "path", r.URL.Path, "query", r.URL.Query(),
			"headers", r.Header)
		next.ServeHTTP(w, r)
	})
}
