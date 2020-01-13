package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/handler"
	"tourtoster/user"
)

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/authorization" {
			next.ServeHTTP(w, r)

			return
		}

		u, err := m.authUser(r)
		if err != nil {
			log.Error("auth error", "error", err.Error())
		}

		if u == nil && r.URL.Path != handler.MainAuthorizationPagePath {
			http.Redirect(w, r, handler.MainAuthorizationPagePath, http.StatusFound)

			return
		}

		if u != nil && r.URL.Path == handler.MainAuthorizationPagePath {
			http.Redirect(w, r, handler.MainPagePath, http.StatusFound)

			return
		}

		context.Set(r, "user", u)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) authUser(r *http.Request) (*user.User, error) {
	tokenCookie, cookieErr := r.Cookie("token")
	if cookieErr != nil || tokenCookie == nil {
		if cookieErr != nil && cookieErr != http.ErrNoCookie {
			return nil, errors.Wrapf(cookieErr, "get cookie error")
		}

		return nil, nil
	}

	t, tokenErr := m.token.Token(tokenCookie.Value)
	if tokenErr != nil || t == nil {
		if tokenErr != nil {
			return nil, errors.Wrapf(tokenErr, "token repo error")
		}

		return nil, nil
	}

	u, userErr := m.user.User(t.UserID)
	if userErr != nil || u == nil {
		if userErr != nil {
			return nil, errors.Wrapf(userErr, "user repo error")
		}

		return nil, nil
	}

	u.Token = t

	return u, nil
}
