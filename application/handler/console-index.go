package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	MainPage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleIndexPath         = ""
	ConsoleIndexTemplateName = "console-index"
)

func (h *Handlers) ConsoleIndexPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	data := MainPage{
		Menu: menu{Dashboard: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleIndexTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleIndexTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
