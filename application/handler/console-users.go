package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/hotel"
	"tourtoster/user"
)

type (
	Filters map[string]interface{}

	UserPage struct {
		Menu menu
		Me   *me
		Year int
		//
		Users    []user.User
		Hotels   []hotel.Hotel
		EditUser *user.User
		Filters  Filters
	}
)

const (
	ConsoleUserPath         = "/users"
	ConsoleUserTemplateName = "console-user"
)

func (f Filters) ToURL() template.URL {
	uu := make([]string, 0, 2)

	for key, val := range f {
		if fmt.Sprintf("%d", val) == "-1" {
			continue
		}

		uu = append(uu, fmt.Sprintf("filter_%s=%d", key, val))
	}

	return template.URL(strings.Join(uu, "&"))
}

func (h *Handlers) ConsoleUserPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

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

	// default filter values
	queryFilters := map[string]interface{}{
		"status": user.Status(-1),
		"hotel":  int64(-1),
	}
	errs := make([]error, 0)
	for k := range queryFilters {
		if val := r.URL.Query().Get(fmt.Sprintf("filter_%s", k)); val != "" {
			switch k {
			case "status":
				i, err := strconv.Atoi(val)
				if err != nil {
					errs = append(errs, errors.Wrap(err, "error convert status to int"))

					continue
				}

				statusVal := user.Status(i)
				if err := user.ValidationFilterStatus(statusVal); err != nil {
					errs = append(errs, errors.Wrap(err, "error validate status"))

					continue
				}

				queryFilters[k] = statusVal
			case "hotel":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					errs = append(errs, errors.Wrap(err, "error convert hotel to int64"))

					continue
				}

				if i < 0 {
					// no filter
					continue
				}

				if err := hotel.ValidationFilterHotelID(i); err != nil {
					errs = append(errs, errors.Wrap(err, "error validate hotel"))

					continue
				}

				queryFilters[k] = i
			}
		}
	}

	if len(errs) != 0 {
		log.Error("Error validate input params", "query", r.URL.Query(),
			"errors", fmt.Sprintf("%#v", errs))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	users, errUsers := h.user.List(queryFilters)
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
		Filters:  queryFilters,
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
