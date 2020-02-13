package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/hotel"
	"tourtoster/user"
)

type (
	UserPage struct {
		Menu menu
		Me   *me
		Year int
		//
		Users    []user.User
		Hotels   []hotel.Hotel
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

	editUserStr := r.URL.Query().Get("edit_id")
	editUser, editUserErr := h.editUser(editUserStr)
	if editUserErr != nil {
		log.Warn("Error get user to edit", "value", editUserStr, "error", editUserErr.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	users, errUsers := h.user.List()
	if errUsers != nil {
		log.Error("Error get user list", "error", errUsers.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	hotels, errHotels := h.hotel.List()
	if errHotels != nil {
		log.Error("Error get hotel list", "error", errHotels.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	data := UserPage{
		Menu:     menu{Users: true},
		Me:       templateMe(u),
		Year:     time.Now().Year(),
		Users:    users,
		Hotels:   hotels,
		EditUser: editUser,
	}

	if err := h.templates[ConsoleUserTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleUserTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) editUser(idStr string) (*user.User, error) {
	if idStr == "" {
		return &user.User{
			Hotel: &hotel.Hotel{},
		}, nil
	}

	editUserID, parseErr := toInt64(idStr)
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "error parse toInt64")
	}

	u, err := h.user.User(editUserID)
	if err != nil {
		return nil, errors.Wrap(err, "error find user")
	}

	if u == nil {
		return nil, errors.New("user not found")
	}

	return u, nil
}
