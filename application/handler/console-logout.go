package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

const (
	ConsoleSignoutPath = "/signout"
)

func (h *Handlers) ConsoleSignoutPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if err := h.token.Delete(u.Token.Token); err != nil {
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

	http.Redirect(w, r, ConsolePathPrefix+ConsoleIndexPath, http.StatusFound)
}
