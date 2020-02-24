package handler

import (
	"net/http"
	"time"

	"github.com/mgutz/logxi/v1"
	"golang.org/x/crypto/bcrypt"

	"tourtoster/token"
	"tourtoster/user"
)

type (
	authResp struct {
		Token string `json:"token"`
	}
)

var (
	authError = respError{
		Error: "login or password is invalid",
	}

	unconfirmedError = respError{
		Error: "user is unconfirmed yet",
	}

	disabledError = respError{
		Error: "user is disabled",
	}
)

const (
	AuthorizationApiPath = "/authorization"
)

func (h *Handlers) AuthorizationApi(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	values := r.Form

	u, uErr := h.user.UserWithEmail(values.Get("email"))
	if uErr != nil {
		log.Error("Error find user", "error", uErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	if u == nil || !checkPasswordHash(values.Get("password"), u.PasswordHash) {
		w.WriteHeader(http.StatusUnauthorized)
		write(w, authError)

		return
	}

	switch u.Status {
	case user.StatusNew:
		w.WriteHeader(http.StatusUnauthorized)
		write(w, unconfirmedError)

		return
	case user.StatusDisabled:
		w.WriteHeader(http.StatusUnauthorized)
		write(w, disabledError)

		return
	}

	newToken, hashErr := HashPassword(u.Email + u.PasswordHash + time.Now().String())
	if hashErr != nil {
		log.Error("Error generate auth token", "error", hashErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	if err := h.token.Save(&token.Token{UserID: u.ID, Token: newToken}); err != nil {
		log.Error("Error insert new token", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	resp := &authResp{
		Token: newToken,
	}

	write(w, resp)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
