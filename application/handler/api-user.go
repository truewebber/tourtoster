package handler

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/hotel"
	"tourtoster/user"
)

const (
	UserApiPath = "/user"
)

var (
	userEmailInvalidError = respError{
		Error: "E-mail empty or invalid",
	}

	userPhoneInvalidError = respError{
		Error: "Phone number empty or invalid",
	}

	userEmailPhoneUniqueError = respError{
		Error: "Phone number or E-mail is not unique",
	}

	userFNameInvalidError = respError{
		Error: "First name empty or invalid",
	}

	userLNameInvalidError = respError{
		Error: "Last name empty or invalid",
	}

	userHotelInvalidError = respError{
		Error: "Select user hotel",
	}
)

func (h *Handlers) ApiUserCreate(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)

		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	values := r.Form

	ID, IDErr := toInt64(strings.TrimSpace(values.Get("id")))
	if IDErr != nil {
		log.Error("Error parse project id", "error", IDErr.Error(), "id", values.Get("id"))
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

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

	note := html.EscapeString(strings.TrimSpace(values.Get("note")))

	// hotel:start
	saveAsNewHotel := checkbox(html.EscapeString(strings.TrimSpace(values.Get("save_new_hotel"))))
	hotelName := html.EscapeString(strings.TrimSpace(values.Get("hotel_name")))
	hotelID, hotelIDErr := toInt64(html.EscapeString(strings.TrimSpace(values.Get("hotel_id"))))
	if hotelIDErr != nil {
		log.Error("hotel_id is not int value", "error", hotelIDErr.Error(),
			"hotel_id", values.Get("hotel_id"))
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	htl, hotelErr := h.newHotel(hotelID, hotelName, saveAsNewHotel)
	if hotelErr != nil {
		log.Error("user save hotel error", "error", hotelErr.Error(),
			"hotel_id", hotelID, "hotel_name", hotelName)
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
	if htl == nil {
		w.WriteHeader(http.StatusBadRequest)
		write(w, userHotelInvalidError)

		return
	}
	// hotel:end

	// status:start
	statusInt, parseStatusErr := toInt(strings.TrimSpace(values.Get("status")))
	if parseStatusErr != nil {
		log.Error("Error parse status id", "error", parseStatusErr.Error(), "status", values.Get("status"))
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}
	status := user.Status(statusInt)

	if err := user.ValidationStatus(status); err != nil {
		log.Error("Status invalid", "error", err.Error(), "status", values.Get("status"))
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}
	// status:end

	// permissions:start
	var pp user.Permission
	for _, pStr := range values["permission[]"] {
		pInt, err := toInt(strings.TrimSpace(pStr))
		if err != nil {
			log.Error("Error parse permission id", "error", err.Error(), "permission", pStr)
			w.WriteHeader(http.StatusBadRequest)
			write(w, inputInvalidError)

			return
		}
		p := user.Permission(pInt)

		if err := user.ValidationPermission(p); err != nil {
			log.Error("Permission invalid", "error", err.Error(), "permission", p)
			w.WriteHeader(http.StatusBadRequest)
			write(w, inputInvalidError)

			return
		}

		pp = pp | p
	}
	// permissions:end

	// password:start
	password := html.EscapeString(strings.TrimSpace(values.Get("password")))
	if password == "" && ID == 0 {
		password = RandString(10)
	}

	passwordHash := ""
	if password != "" {
		var hashErr error
		passwordHash, hashErr = HashPassword(password)
		if hashErr != nil {
			log.Error("Error hash password", "error", hashErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			write(w, internalError)

			return
		}
	}
	// password:end

	if checkbox(values.Get("send_mail")) {
		body := []byte(
			"Welcome to Tourtoster! \n" +
				fmt.Sprintf("Your new password: `%s`", password),
		)

		if err := h.mailer.Send(email, body); err != nil {
			log.Error("Error send email", "error", err.Error())
		}
	}

	newUser := &user.User{
		ID:           ID,
		FirstName:    firstName,
		SecondName:   secondName,
		LastName:     lastName,
		Hotel:        htl,
		Note:         note,
		Email:        email,
		Phone:        phone,
		Status:       status,
		Permissions:  pp,
		PasswordHash: passwordHash,
	}

	if err := h.user.Save(newUser); err != nil {
		code := http.StatusBadRequest
		e := userEmailPhoneUniqueError

		if errors.Cause(err) != user.PhoneEmailUniqueError {
			log.Error("Error create/save user", "error", err.Error(), "id", ID)
			code = http.StatusInternalServerError
			e = internalError
		}

		w.WriteHeader(code)
		write(w, e)

		return
	}

	w.WriteHeader(http.StatusOK)
	write(w, u)
}

func (h *Handlers) ApiUseDelete(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, "user").(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		write(w, forbiddenError)

		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Error("Error read body", "error", readErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	values, parseErr := url.ParseQuery(string(body))
	if parseErr != nil {
		log.Warn("Error parse body", "error", parseErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	stringID := strings.TrimSpace(values.Get("id"))
	ID, convertErr := strconv.ParseInt(stringID, 10, 64)
	if convertErr != nil {
		log.Warn("Error parse user id", "error", convertErr.Error(), "id", stringID)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	if err := h.user.Delete(ID); err != nil {
		log.Error("Error delete project", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) newHotel(hotelID int64, hotelName string, new bool) (*hotel.Hotel, error) {
	if hotelID != 0 {
		return h.hotel.Hotel(hotelID)
	}

	if hotelName == "" {
		return nil, nil
	}

	htl := &hotel.Hotel{Name: hotelName}
	if new {
		if err := h.hotel.Save(htl); err != nil {
			return nil, err
		}
	}

	return htl, nil
}

func checkbox(s string) bool {
	if s != "" {
		return true
	}

	return false
}

func toInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	return strconv.ParseInt(s, 10, 64)
}

func toInt(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	return strconv.Atoi(s)
}
