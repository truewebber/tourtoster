package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"

	"github.com/truewebber/tourtoster/user"
)

type (
	UserProfilePage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleUserProfilePath         = "/user/profile"
	ConsoleUserProfileTemplateName = "console-user-profile"
)

func (h *Handlers) ConsoleUserProfilePage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	data := UserProfilePage{
		Menu: menu{},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleUserProfileTemplateName].Execute(w, data); err != nil {
		h.logger.Error("Error execute template", "template", ConsoleUserProfileTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
