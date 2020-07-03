package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	AboutPage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleAboutPath         = "/about"
	ConsoleAboutTemplateName = "console-about"
)

func (h *Handlers) ConsoleAboutPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	data := AboutPage{
		Menu: menu{},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleAboutTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleAboutTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
