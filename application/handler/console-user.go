package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	UserPage struct {
		Menu menu
		Me   *me
		Year int
		//
		Users    []user.User
		EditUser *user.User
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

	var editUser *user.User
	val := r.URL.Query().Get("edit_user")
	if val != "" {
		editUserID, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Warn("Error parse edit user ID", "value", val, "error", err.Error())
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		var errUser error
		editUser, errUser = h.user.User(editUserID)
		if errUser != nil {
			log.Error("Error get edit user", "error", errUser.Error())
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if editUser == nil {
			log.Warn("no user with such ID", "value", val)
			w.WriteHeader(http.StatusBadRequest)

			return
		}
	}

	users, errUsers := h.user.List()
	if errUsers != nil {
		log.Error("Error get user list", "error", errUsers.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data := UserPage{
		Menu:     menu{Users: true},
		Me:       templateMe(u),
		Year:     time.Now().Year(),
		Users:    users,
		EditUser: editUser,
	}

	if err := h.templates[ConsoleUserTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleUserTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
