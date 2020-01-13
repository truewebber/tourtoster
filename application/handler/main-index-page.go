package handler

import (
	"net/http"
	"time"

	"github.com/mgutz/logxi/v1"
)

type (
	MainPage struct {
		Year int
	}
)

const (
	MainPagePath     = "/"
	MainTemplateName = "main-index"
)

func (h *Handlers) MainPage(w http.ResponseWriter, r *http.Request) {
	data := MainPage{
		Year: time.Now().Year(),
	}

	err := h.templates[MainTemplateName].Execute(w, data)
	if err != nil {
		log.Error("Error execute template", "template", MainTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
