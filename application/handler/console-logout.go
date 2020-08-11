package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"

	"github.com/truewebber/tourtoster/user"
)

const (
	ConsoleSignoutPath = "/signout"
)

func (h *Handlers) ConsoleSignoutPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	if err := h.token.Delete(u.Token.Token); err != nil {
		h.logger.Error("Error delete token", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, "internal error")

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
