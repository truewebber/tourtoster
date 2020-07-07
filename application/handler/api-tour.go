package handler

import (
	"html"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/teambition/rrule-go"

	"tourtoster/tour"
	"tourtoster/user"
)

const (
	TourApiPath = "/tour"
)

func (h *Handlers) ApiTourCreate(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.EditToursPermission) {
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

	title := html.EscapeString(strings.TrimSpace(values.Get("title")))
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)
		return
	}

	image := html.EscapeString(strings.TrimSpace(values.Get("image")))
	if image == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)
		return
	}

	description := strings.TrimSpace(values.Get("description"))
	if description == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)
		return
	}

	//rruleFreqInt, err := toInt(values.Get("rrule_freq"))
	//if err != nil {
	//	log.Error("Error parse rrule freq", "error", err.Error(), "rrule_freq", values.Get("rrule_freq"))
	//	w.WriteHeader(http.StatusBadRequest)
	//	h.write(w, inputInvalidError)
	//
	//	return
	//}
	//rruleFreq := rrule.Frequency(rruleFreqInt)

	rs := &rrule.Set{}
	//rs.RRule(r)

	t := &tour.Tour{
		Recurrence: rs,
	}

	if err := h.tour.Save(t); err != nil {
		h.logger.Error("Error create/save tour", "error", err.Error(), "tour", t)
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}
}

func (h *Handlers) ApiTourDelete(w http.ResponseWriter, r *http.Request) {

}
