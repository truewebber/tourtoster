package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	AdminPage struct {
		User         *user.User
		AllowUserAdd bool
		Hostname     string
		Year         int
	}
)

const (
	AdminPageIndexPath     = ""
	AdminIndexTemplateName = "admin-index"
)

func (h *Handlers) AdminIndexPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	data := AdminPage{
		User:         u,
		AllowUserAdd: true,
		Hostname:     r.Host,
		Year:         time.Now().Year(),
	}

	err := h.templates[AdminIndexTemplateName].Execute(w, data)
	if err != nil {
		log.Error("Error execute template", "template", AdminIndexTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
