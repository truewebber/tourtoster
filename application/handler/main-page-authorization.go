package handler

import (
	"net/http"

	"github.com/mgutz/logxi/v1"
)

const (
	MainPageAuthorizationPath         = "/authorization"
	MainPageAuthorizationTemplateName = "main-authorization"
)

func (h *Handlers) MainAuthorizationPage(w http.ResponseWriter, _ *http.Request) {
	err := h.templates[MainPageAuthorizationTemplateName].Execute(w, nil)
	if err != nil {
		log.Error("Error execute template", "template", MainPageAuthorizationTemplateName,
			"error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
