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
		u, err := m.authUser(r)
		if err != nil {
			log.Error("auth error", "error", err.Error())
		}

		if u == nil && r.URL.Path != handler.ConsolePathPrefix+handler.MainPageAuthorizationPath {
			http.Redirect(w, r, handler.ConsolePathPrefix+handler.MainPageAuthorizationPath, http.StatusFound)

			return
		}

		if u != nil && r.URL.Path == handler.ConsolePathPrefix+handler.MainPageAuthorizationPath {
			http.Redirect(w, r, handler.ConsolePathPrefix+handler.MainPageIndexPath, http.StatusFound)

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

	if u.Hotel.ID != 0 {
		h, err := m.hotel.Hotel(u.Hotel.ID)
		if err != nil || h != nil {
			if err != nil {
				return nil, errors.Wrapf(userErr, "hotel repo error")
			}

			return nil, nil
		}
		u.Hotel = h
	}

	u.Token = t

	return u, nil
}
