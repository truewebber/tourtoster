package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	ContactPage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleContactPath         = "/contact"
	ConsoleContactTemplateName = "console-contact"
)

func (h *Handlers) ConsoleContactPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	data := ContactPage{
		Menu: menu{},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleContactTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleContactTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
