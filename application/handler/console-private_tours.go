package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/pkg/errors"

	"tourtoster/tour"
	"tourtoster/user"
)

type (
	PTPage struct {
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
	ConsolePTPath        = "/private_tours"
	ConsolePTEditPath    = "/private_tours/edit"
	ConsolePTEditFAQPath = "/private_tours/edit/faq"

	ConsolePTTemplateName           = "console-private_tours"
	ConsolePTViewAllSubTemplateName = "console-private_tours-view_all"
	ConsolePTEditSubTemplateName    = "console-private_tours-edit"
	ConsolePTEditFAQSubTemplateName = "console-private_tours-edit-faq"
)

func (h *Handlers) ConsolePTPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	h.renderPT(w, &PTPage{
		Page: ConsolePTViewAllSubTemplateName,
		Menu: menu{PrivateTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	})
}

func (h *Handlers) ConsolePTEditPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	editTourStr := r.URL.Query().Get("id")
	editTour, editUserErr := h.editPrivateTour(editTourStr)
	if editUserErr != nil {
		h.logger.Warn("Error get private tour to edit", "value", editTourStr, "error", editUserErr.Error())
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	tours, err := h.tour.List(tour.NewOrder("title", "asc"), tour.FilterTourType(tour.PrivateType))
	if err != nil {
		h.logger.Error("Error get private tours list", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	h.renderPT(w, &PTPage{
		Tours:    tours,
		EditTour: editTour,
		Page:     ConsolePTEditSubTemplateName,
		Menu:     menu{PrivateTours: true},
		Me:       templateMe(u),
		Year:     time.Now().Year(),
	})
}

func (h *Handlers) ConsolePTEditFAQPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	h.renderPT(w, &PTPage{
		Page: ConsolePTEditFAQSubTemplateName,
		Menu: menu{PrivateTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	})
}

func (h *Handlers) renderPT(w http.ResponseWriter, ptPage *PTPage) {
	if err := h.templates[ConsolePTTemplateName].Execute(w, *ptPage); err != nil {
		h.logger.Error("Error execute template", "template", ConsolePTTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (h *Handlers) editPrivateTour(idStr string) (*tour.Tour, error) {
	if idStr == "" {
		return new(tour.Tour), nil
	}

	editTourID, parseErr := toInt64(idStr)
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "error parse toInt64")
	}

	t, err := h.tour.Tour(editTourID)
	if err != nil {
		return nil, errors.Wrap(err, "error find private tour")
	}

	if t == nil {
		return nil, errors.New("private tour not found")
	}

	if t.Type != tour.PrivateType {
		return nil, errors.Errorf("tour type is not Private, %v", t.Type)
	}

	return t, nil
}
