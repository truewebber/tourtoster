package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	UserBillingPage struct {
		Menu menu
		Me   *me
		Year int
	}
)

const (
	ConsoleUserBillingPath         = "/user/billing"
	ConsoleUserBillingTemplateName = "console-user-billing"
)

func (h *Handlers) ConsoleUserBillingPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)

	data := UserBillingPage{
		Menu: menu{},
		Me:   templateMe(u),
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleUserBillingTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleUserBillingTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
