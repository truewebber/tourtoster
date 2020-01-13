package handler

import (
	"net/http"

	"github.com/mgutz/logxi/v1"
)

const (
	MainAuthorizationPagePath     = "/authorization"
	MainAuthorizationTemplateName = "main-authorization"
)

func (h *Handlers) MainAuthorizationPage(w http.ResponseWriter, _ *http.Request) {
	err := h.templates[MainAuthorizationTemplateName].Execute(w, nil)
	if err != nil {
		log.Error("Error execute template", "template", MainAuthorizationTemplateName, "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
