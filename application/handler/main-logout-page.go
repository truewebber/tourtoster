package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

const (
	LogoutPagePath = "/logout"
)

func (h *Handlers) LogoutAdminPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	err := h.token.Delete(u.Token.Token)
	if err != nil {
		log.Error("Error delete token", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, "internal error")

		return

	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	http.Redirect(w, r, LogoutPagePath, http.StatusFound)
}
