package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"

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
	u := context.Get(r, user.ContextKey).(*user.User)

	data := MainPage{
		Menu: menu{Dashboard: true},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleIndexTemplateName].Execute(w, data); err != nil {
		h.logger.Error("Error execute template", "template", ConsoleIndexTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
