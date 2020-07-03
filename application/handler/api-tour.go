package handler

import (
	"net/http"

	"github.com/gorilla/context"

	"tourtoster/user"
)

const (
	TourApiPath = "/tour"
)

func (h *Handlers) ApiTourCreate(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.EditToursPermission) {
		w.WriteHeader(http.StatusForbidden)
		write(w, forbiddenError)

		return
	}
}

func (h *Handlers) ApiTourDelete(w http.ResponseWriter, r *http.Request) {

}
