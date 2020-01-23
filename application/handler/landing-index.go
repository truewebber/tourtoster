package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"

	"tourtoster/user"
)

type (
	LandingPage struct {
		User *user.User
		Year int
	}
)

const (
	LandingPageIndexPath     = "/"
	LandingIndexTemplateName = "landing-index"
)

func (h *Handlers) LandingIndexPage(w http.ResponseWriter, r *http.Request) {
	u, _ := context.Get(r, "user").(*user.User)

	data := LandingPage{
		User: u,
		Year: time.Now().Year(),
	}

	if err := h.templates[LandingIndexTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", LandingIndexTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
