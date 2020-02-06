package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	UserPage struct {
		Menu menu
		User *htmlUser
		Year int
	}
)

const (
	ConsoleUserPath         = "/users"
	ConsoleUserTemplateName = "console-user"
)

func (h *Handlers) ConsoleUserPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	data := UserPage{
		Menu: menu{Users: true},
		User: templateUser(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleUserTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleUserTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
