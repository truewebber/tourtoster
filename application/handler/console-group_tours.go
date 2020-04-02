package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	GTPage struct {
		Page string
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleGTPath        = "/group_tours"
	ConsoleGTEditPath    = "/group_tours/edit"
	ConsoleGTEditFAQPath = "/group_tours/edit/faq"

	ConsoleGTTemplateName = "console-group_tours"

	ConsoleGTViewAllSubTemplateName = "console-group_tours-view_all"
	ConsoleGTEditSubTemplateName    = "console-group_tours-edit"
	ConsoleGTEditFAQSubTemplateName = "console-group_tours-edit-faq"
)

func (h *Handlers) ConsoleGTPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	h.renderGT(w, u, ConsoleGTViewAllSubTemplateName)
}

func (h *Handlers) ConsoleGTEditPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	h.renderGT(w, u, ConsoleGTEditSubTemplateName)
}

func (h *Handlers) ConsoleGTEditFAQPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	h.renderGT(w, u, ConsoleGTEditFAQSubTemplateName)
}

func (h *Handlers) renderGT(w http.ResponseWriter, u *user.User, subTemplate string) {
	data := GTPage{
		Page: subTemplate,
		Menu: menu{GroupTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleGTTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleGTTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
