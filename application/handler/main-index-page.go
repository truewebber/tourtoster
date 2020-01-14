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
		User *user.User
		Year int
	}
)

const (
	MainPagePath     = "/"
	MainTemplateName = "main-index"
)

func (h *Handlers) MainPage(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)

	data := MainPage{
		User: u,
		Year: time.Now().Year(),
	}

	err := h.templates[MainTemplateName].Execute(w, data)
	if err != nil {
		log.Error("Error execute template", "template", MainTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
