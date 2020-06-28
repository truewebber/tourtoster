package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/tour"
	"tourtoster/user"
)

type (
	GTPage struct {
		Menu menu
		Me   *me
		Year int
		//
		Page     string
		Tours    []tour.Tour
		EditTour *tour.Tour
	}
)

const (
	ConsoleGTPath        = "/group_tours"
	ConsoleGTEditPath    = "/group_tours/edit"
	ConsoleGTEditFAQPath = "/group_tours/edit/faq"

	ConsoleGTTemplateName           = "console-group_tours"
	ConsoleGTViewAllSubTemplateName = "console-group_tours-view_all"
	ConsoleGTEditSubTemplateName    = "console-group_tours-edit"
	ConsoleGTEditFAQSubTemplateName = "console-group_tours-edit-faq"
)

func (h *Handlers) ConsoleGTPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	h.renderGT(w, &GTPage{
		Page: ConsoleGTViewAllSubTemplateName,
		Menu: menu{GroupTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	})
}

func (h *Handlers) ConsoleGTEditPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	editTourStr := r.URL.Query().Get("id")
	editTour, editUserErr := h.editGroupTour(editTourStr)
	if editUserErr != nil {
		log.Warn("Error get group tour to edit", "value", editTourStr, "error", editUserErr.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	tours, err := h.tour.List()
	if err != nil {
		log.Error("Error get group tours list", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.renderGT(w, &GTPage{
		Tours:    tours,
		EditTour: editTour,
		Page:     ConsoleGTEditSubTemplateName,
		Menu:     menu{GroupTours: true},
		Me:       templateMe(u),
		Year:     time.Now().Year(),
	})
}

func (h *Handlers) ConsoleGTEditFAQPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	h.renderGT(w, &GTPage{
		Page: ConsoleGTEditFAQSubTemplateName,
		Menu: menu{GroupTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	})
}

func (h *Handlers) renderGT(w http.ResponseWriter, gtPage *GTPage) {
	if err := h.templates[ConsoleGTTemplateName].Execute(w, *gtPage); err != nil {
		log.Error("Error execute template", "template", ConsoleGTTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) editGroupTour(idStr string) (*tour.Tour, error) {
	if idStr == "" {
		return new(tour.Tour), nil
	}

	editTourID, parseErr := toInt64(idStr)
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "error parse toInt64")
	}

	t, err := h.tour.Tour(editTourID)
	if err != nil {
		return nil, errors.Wrap(err, "error find group tour")
	}

	if t == nil {
		return nil, errors.New("group tour not found")
	}

	return t, nil
}
