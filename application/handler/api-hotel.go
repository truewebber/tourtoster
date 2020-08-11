package handler

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/context"

	"github.com/truewebber/tourtoster/hotel"
	"github.com/truewebber/tourtoster/user"
)

const (
	HotelApiPath = "/hotel"
)

func (h *Handlers) ApiHotelList(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		h.write(w, forbiddenError)

		return
	}

	hh, err := h.hotel.List()
	if err != nil {
		h.logger.Error("Error list hotel", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	h.write(w, hh)
}

func (h *Handlers) ApiHotelCreate(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		h.write(w, forbiddenError)

		return
	}

	if err := r.ParseForm(); err != nil {
		h.logger.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}
	values := r.Form

	name := html.EscapeString(strings.TrimSpace(values.Get("name")))
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)
		return
	}

	ID, IDErr := toInt64(strings.TrimSpace(values.Get("id")))
	if IDErr != nil {
		h.logger.Error("Error parse hotel id", "error", IDErr.Error(), "id", values.Get("id"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	newHotel := &hotel.Hotel{
		ID:   ID,
		Name: name,
	}

	if err := h.hotel.Save(newHotel); err != nil {
		h.logger.Error("Error create/save hotel", "error", err.Error(), "id", ID, "name", name)
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	h.write(w, newHotel)
}

func (h *Handlers) ApiHotelDelete(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		h.write(w, forbiddenError)

		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		h.logger.Error("Error read body", "error", readErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	values, parseErr := url.ParseQuery(string(body))
	if parseErr != nil {
		h.logger.Warn("Error parse body", "error", parseErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	stringID := strings.TrimSpace(values.Get("id"))
	ID, convertErr := strconv.ParseInt(stringID, 10, 64)
	if convertErr != nil {
		h.logger.Warn("Error parse hotel id", "error", convertErr.Error(), "id", stringID)
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	if err := h.hotel.Delete(ID); err != nil {
		h.logger.Error("Error delete hotel", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
