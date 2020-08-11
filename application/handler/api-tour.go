package handler

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/teambition/rrule-go"

	"github.com/truewebber/tourtoster/tour"
	"github.com/truewebber/tourtoster/user"
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

	tourTypeID, tourTypeErr := intValue("tour_type", r.Form)
	if tourTypeErr != nil {
		h.logger.Error("Error parse tourType", "error", tourTypeErr.Error(),
			"tour_type", r.Form.Get("tour_type"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}
	if err := tour.ValidateType(tour.Type(tourTypeID)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	title := stringEscapedValue("title", r.Form)
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	statusID, statusErr := intValue("tour_status", r.Form)
	if statusErr != nil {
		h.logger.Error("Error parse tour status", "error", statusErr.Error(),
			"tour_status", r.Form.Get("tour_status"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}
	if err := tour.ValidateStatus(tour.Status(statusID)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	rruleSetStr := stringValue("rrule_set", r.Form)
	s, err := rrule.StrToRRuleSet(rruleSetStr)
	if err != nil {
		h.logger.Error("Error parse rrule_set", "error", err.Error(),
			"rrule_set", r.Form.Get("rrule_set"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	image := stringEscapedValue("image", r.Form)
	if image == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	description := stringValue("description", r.Form)
	if description == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	pricePerAdults, adultsErr := intValue("price_per_adults", r.Form)
	if adultsErr != nil {
		h.logger.Error("Error parse pricePerAdults", "error", adultsErr.Error(),
			"price_per_adults", r.Form.Get("price_per_adults"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	pricePerChildrenSevenSeventeen, children717Err := intValue("price_per_children_seven_seventeen", r.Form)
	if children717Err != nil {
		h.logger.Error("Error parse pricePerChildrenSevenSeventeen", "error", children717Err.Error(),
			"price_per_children_seven_seventeen", r.Form.Get("price_per_children_seven_seventeen"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	pricePerChildrenZeroSix, children06Err := intValue("price_per_children_zero_six", r.Form)
	if children06Err != nil {
		h.logger.Error("Error parse pricePerChildrenZeroSix", "error", children06Err.Error(),
			"price_per_children_zero_six", r.Form.Get("price_per_children_zero_six"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	if pricePerAdults < 0 || pricePerChildrenSevenSeventeen < 0 || pricePerChildrenZeroSix < 0 {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	t := &tour.Tour{
		Type:                           tour.Type(tourTypeID),
		Creator:                        u,
		Status:                         tour.Status(statusID),
		Recurrence:                     s,
		Title:                          title,
		Image:                          image,
		Description:                    description,
		PricePerAdults:                 tour.NewRUB(pricePerAdults),
		PricePerChildrenSevenSeventeen: tour.NewRUB(pricePerChildrenSevenSeventeen),
		PricePerChildrenZeroSix:        tour.NewRUB(pricePerChildrenZeroSix),
		//disabled
		Map:                      "",
		MaxPersons:               0,
		PricePerChildrenThreeSix: 0,
	}

	if err := h.tour.Save(t); err != nil {
		h.logger.Error("Error create/save tour", "error", err.Error(), "tour", t)
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	h.write(w, t)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) ApiTourDelete(w http.ResponseWriter, r *http.Request) {

}
