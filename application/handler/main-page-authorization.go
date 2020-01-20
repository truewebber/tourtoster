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
	if err := h.templates[MainPageAuthorizationTemplateName].Execute(w, nil); err != nil {
		log.Error("Error execute template", "template", MainPageAuthorizationTemplateName,
			"error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
