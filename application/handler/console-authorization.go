package handler

import (
	"net/http"

	"github.com/mgutz/logxi/v1"
)

const (
	ConsoleAuthorizationPath         = "/authorization"
	ConsoleAuthorizationTemplateName = "console-authorization"
)

func (h *Handlers) ConsoleAuthorizationPage(w http.ResponseWriter, _ *http.Request) {
	if err := h.templates[ConsoleAuthorizationTemplateName].Execute(w, nil); err != nil {
		log.Error("Error execute template", "template", ConsoleAuthorizationTemplateName,
			"error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
