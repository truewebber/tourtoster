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

	if !u.AdminPage() {
		http.Redirect(w, r, MainPageIndexPath, http.StatusFound)

		return
	}

	data := AdminPage{
		User:         u,
		AllowUserAdd: true,
		Hostname:     r.Host,
		Year:         time.Now().Year(),
	}

	if err := h.templates[AdminIndexTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", AdminIndexTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
