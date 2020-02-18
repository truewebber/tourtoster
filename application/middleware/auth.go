package middleware

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/handler"
	"tourtoster/user"
)

func (m *Middleware) PageAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := m.authUser(r)
		if err != nil {
			log.Error("auth error", "error", err.Error())
		}

		path1 := handler.ConsolePathPrefix + handler.ConsoleAuthorizationPath
		path2 := handler.ConsolePathPrefix + handler.ConsoleRegistrationPath

		// unauthorized and page except auth -> redirect to auth
		if u == nil && !(r.URL.Path == path1 || r.URL.Path == path2) {
			http.Redirect(w, r, handler.ConsolePathPrefix+handler.ConsoleAuthorizationPath, http.StatusFound)

			return
		}

		// authorized and page auth -> redirect to dashboard
		if u != nil && (r.URL.Path == path1 || r.URL.Path == path2) {
			http.Redirect(w, r, handler.ConsolePathPrefix+handler.ConsoleIndexPath, http.StatusFound)

			return
		}

		context.Set(r, "user", u)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) APIAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := m.authUser(r)
		if err != nil {
			log.Error("auth error", "error", err.Error())
		}

		path := handler.ApiPathPrefix + handler.AuthorizationApiPath
		if u == nil && r.URL.Path != path {
			w.WriteHeader(http.StatusUnauthorized)
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

	if u.Status != user.StatusEnabled {
		return nil, nil
	}

	if u.Hotel.ID != 0 {
		h, err := m.hotel.Hotel(u.Hotel.ID)
		if err != nil {
			return nil, errors.Wrapf(userErr, "hotel repo error")
		}
		if h == nil {
			return nil, nil
		}
		u.Hotel = h
	}

	u.Token = t

	return u, nil
}
