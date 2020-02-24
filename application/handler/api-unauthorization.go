package handler

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/mgutz/logxi/v1"
)

const (
	ForgetApiPath       = "/forget"
	RegistrationApiPath = "/registration"
)

var (
	unknownUserError = respError{
		Error: "E-mail not found",
	}
)

func (h *Handlers) ApiForget(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	values := r.Form

	email := html.EscapeString(strings.TrimSpace(values.Get("email")))
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userEmailInvalidError)
		return
	}

	u, uErr := h.user.UserWithEmail(email)
	if uErr != nil {
		log.Error("Error get user", "email", email, "error", uErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, inputInvalidError)

		return
	}
	if u == nil {
		log.Warn("User not found", "email", email)
		w.WriteHeader(http.StatusNotFound)
		write(w, unknownUserError)

		return
	}

	// password:start
	password := RandString(10)
	passwordHash, hashErr := HashPassword(password)
	if hashErr != nil {
		log.Error("Error hash password", "error", hashErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	// password:end

	if err := h.user.Password(u.ID, passwordHash); err != nil {
		log.Error("Error set new password", "error", err.Error(), "id", u.ID)
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	title := "New password on Tourtoster!"
	body := "Hello! \n" + fmt.Sprintf("Your new password: `%s`", password)

	go func() {
		if err := h.mailer.Send(email, title, body); err != nil {
			log.Error("Error send email", "error", err.Error())
		}
	}()

	w.WriteHeader(http.StatusOK)
	write(w, respOK{Message: "success"})
}

func (h *Handlers) ApiRegistration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	//values := r.Form

	w.WriteHeader(http.StatusNoContent)
}
