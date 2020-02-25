package handler

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/user"
)

const (
	ForgetApiPath       = "/forget"
	RegistrationApiPath = "/registration"
)

var (
	unknownUserError = respError{
		Error: "E-mail not found",
	}

	passwordEmptyError = respError{
		Error: "Password field is empty",
	}

	passwordInvalidError = respError{
		Error: "Password fields not equal",
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
	values := r.Form

	firstName := html.EscapeString(strings.TrimSpace(values.Get("first_name")))
	if firstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userFNameInvalidError)
		return
	}
	secondName := html.EscapeString(strings.TrimSpace(values.Get("second_name")))
	lastName := html.EscapeString(strings.TrimSpace(values.Get("last_name")))
	if lastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userLNameInvalidError)
		return
	}

	email := html.EscapeString(strings.TrimSpace(values.Get("email")))
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userEmailInvalidError)
		return
	}
	phone := html.EscapeString(strings.TrimSpace(values.Get("phone")))
	if phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userPhoneInvalidError)
		return
	}

	hotelName := html.EscapeString(strings.TrimSpace(values.Get("hotel_name")))
	htl, htlErr := h.newHotel(0, hotelName, false)
	if htlErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)
		return
	}

	// password:start
	password := html.EscapeString(strings.TrimSpace(values.Get("password")))
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, passwordEmptyError)
		return
	}

	passwordRepeat := html.EscapeString(strings.TrimSpace(values.Get("password_repeat")))
	if passwordRepeat == "" {
		w.WriteHeader(http.StatusBadRequest)
		write(w, passwordEmptyError)
		return
	}

	if password != passwordRepeat {
		w.WriteHeader(http.StatusBadRequest)
		write(w, passwordInvalidError)
		return
	}

	passwordHash, hashErr := HashPassword(password)
	if hashErr != nil {
		log.Error("Error hash password", "error", hashErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	// password:end

	newUser := &user.User{
		FirstName:    firstName,
		SecondName:   secondName,
		LastName:     lastName,
		Hotel:        htl,
		Note:         "",
		Email:        email,
		Phone:        phone,
		Status:       user.StatusNew,
		Permissions:  user.Permission(0),
		PasswordHash: passwordHash,
	}

	if err := h.user.Save(newUser); err != nil {
		code := http.StatusBadRequest
		e := userEmailPhoneUniqueError

		if errors.Cause(err) != user.PhoneEmailUniqueError {
			log.Error("Error create/save user", "error", err.Error())
			code = http.StatusInternalServerError
			e = internalError
		}

		w.WriteHeader(code)
		write(w, e)

		return
	}

	write(w, respOK{"success"})
}
