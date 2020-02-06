package handler

import (
	"net/http"
	"time"

	"github.com/mgutz/logxi/v1"
)

type (
	RegistrationPage struct {
		Year int
	}
)

const (
	ConsoleRegistrationPath         = "/registration"
	ConsoleRegistrationTemplateName = "console-registration"
)

func (h *Handlers) ConsoleRegistrationPage(w http.ResponseWriter, _ *http.Request) {
	data := &RegistrationPage{
		Year: time.Now().Year(),
	}

	if err := h.templates[ConsoleRegistrationTemplateName].Execute(w, data); err != nil {
		log.Error("Error execute template", "template", ConsoleRegistrationTemplateName,
			"error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
