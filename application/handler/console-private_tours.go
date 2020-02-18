package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	PTPage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsolePTPath         = "/private_tours"
	ConsolePTTemplateName = "console-index"
)

func (h *Handlers) ConsolePTPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	data := PTPage{
		Menu: menu{PrivateTours: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsolePTTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsolePTTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
