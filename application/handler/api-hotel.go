package handler

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/hotel"
	"tourtoster/user"
)

const (
	HotelApiPath = "/hotel"
)

func (h *Handlers) ApiHotelList(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		write(w, forbiddenError)

		return
	}

	hh, err := h.hotel.List()
	if err != nil {
		log.Error("Error list hotel", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	write(w, hh)
}

func (h *Handlers) ApiHotelCreate(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		write(w, forbiddenError)

		return
	}

	if err := r.ParseForm(); err != nil {
		log.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	values := r.Form

	name := html.EscapeString(strings.TrimSpace(values.Get("name")))
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)
		return
	}

	ID, IDErr := toInt64(strings.TrimSpace(values.Get("id")))
	if IDErr != nil {
		log.Error("Error parse hotel id", "error", IDErr.Error(), "id", values.Get("id"))
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	newHotel := &hotel.Hotel{
		ID:   ID,
		Name: name,
	}

	if err := h.hotel.Save(newHotel); err != nil {
		log.Error("Error create/save hotel", "error", err.Error(), "id", ID, "name", name)
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	write(w, newHotel)
}

func (h *Handlers) ApiHotelDelete(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		write(w, forbiddenError)

		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Error("Error read body", "error", readErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	values, parseErr := url.ParseQuery(string(body))
	if parseErr != nil {
		log.Warn("Error parse body", "error", parseErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	stringID := strings.TrimSpace(values.Get("id"))
	ID, convertErr := strconv.ParseInt(stringID, 10, 64)
	if convertErr != nil {
		log.Warn("Error parse hotel id", "error", convertErr.Error(), "id", stringID)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	if err := h.hotel.Delete(ID); err != nil {
		log.Error("Error delete hotel", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
