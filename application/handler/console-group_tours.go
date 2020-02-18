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
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleGPPath         = "/group_tours"
	ConsoleGPTemplateName = "console-group_tours"
)

func (h *Handlers) ConsoleGTPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	data := GTPage{
		Menu: menu{GroupTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleGPTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleGPTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
