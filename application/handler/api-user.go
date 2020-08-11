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
	"github.com/pkg/errors"

	"github.com/truewebber/tourtoster/hotel"
	"github.com/truewebber/tourtoster/user"
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
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		h.write(w, forbiddenError)

		return
	}

	if err := r.ParseForm(); err != nil {
		h.logger.Error("Error read form", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}
	values := r.Form

	ID, IDErr := toInt64(strings.TrimSpace(values.Get("id")))
	if IDErr != nil {
		h.logger.Error("Error parse user id", "error", IDErr.Error(), "id", values.Get("id"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	firstName := html.EscapeString(strings.TrimSpace(values.Get("first_name")))
	if firstName == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, userFNameInvalidError)
		return
	}
	secondName := html.EscapeString(strings.TrimSpace(values.Get("second_name")))
	lastName := html.EscapeString(strings.TrimSpace(values.Get("last_name")))
	if lastName == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, userLNameInvalidError)
		return
	}

	email := html.EscapeString(strings.TrimSpace(values.Get("email")))
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, userEmailInvalidError)
		return
	}
	phone := html.EscapeString(strings.TrimSpace(values.Get("phone")))
	if phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, userPhoneInvalidError)
		return
	}

	note := html.EscapeString(strings.TrimSpace(values.Get("note")))

	// hotel:start
	saveAsNewHotel := checkbox(html.EscapeString(strings.TrimSpace(values.Get("save_new_hotel"))))
	hotelName := html.EscapeString(strings.TrimSpace(values.Get("hotel_name")))
	hotelID, hotelIDErr := toInt64(html.EscapeString(strings.TrimSpace(values.Get("hotel_id"))))
	if hotelIDErr != nil {
		h.logger.Error("hotel_id is not int value", "error", hotelIDErr.Error(),
			"hotel_id", values.Get("hotel_id"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	htl, hotelErr := h.newHotel(hotelID, hotelName, saveAsNewHotel)
	if hotelErr != nil {
		h.logger.Error("user save hotel error", "error", hotelErr.Error(),
			"hotel_id", hotelID, "hotel_name", hotelName)
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}
	if htl == nil {
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, userHotelInvalidError)

		return
	}
	// hotel:end

	// status:start
	statusInt, parseStatusErr := toInt(strings.TrimSpace(values.Get("status")))
	if parseStatusErr != nil {
		h.logger.Error("Error parse status id", "error", parseStatusErr.Error(), "status", values.Get("status"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}
	status := user.Status(statusInt)

	if err := user.ValidationStatus(status); err != nil {
		h.logger.Error("Status invalid", "error", err.Error(), "status", values.Get("status"))
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}
	// status:end

	// permissions:start
	var pp user.Permission
	for _, pStr := range values["permission[]"] {
		pInt, err := toInt(strings.TrimSpace(pStr))
		if err != nil {
			h.logger.Error("Error parse permission id", "error", err.Error(), "permission", pStr)
			w.WriteHeader(http.StatusBadRequest)
			h.write(w, inputInvalidError)

			return
		}
		p := user.Permission(pInt)

		if err := user.ValidationPermission(p); err != nil {
			h.logger.Error("Permission invalid", "error", err.Error(), "permission", p)
			w.WriteHeader(http.StatusBadRequest)
			h.write(w, inputInvalidError)

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
			h.logger.Error("Error hash password", "error", hashErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			h.write(w, internalError)

			return
		}
	}
	// password:end

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
			h.logger.Error("Error create/save user", "error", err.Error(), "id", ID)
			code = http.StatusInternalServerError
			e = internalError
		}

		w.WriteHeader(code)
		h.write(w, e)

		return
	}

	if (passwordHash != "") && checkbox(values.Get("send_mail")) {
		title := "New password on Tourtoster!"
		body := "Welcome to Tourtoster! \n" + fmt.Sprintf("Your new password: `%s`", password)

		if ID != 0 {
			body = "Hello! \n" + fmt.Sprintf("Your new password: `%s`", password)
		}

		go func() {
			if err := h.mailer.Send(email, title, body); err != nil {
				h.logger.Error("Error send email", "error", err.Error())
			}
		}()
	}

	w.WriteHeader(http.StatusOK)
	h.write(w, u)
}

func (h *Handlers) ApiUserDelete(w http.ResponseWriter, r *http.Request) {
	u := context.Get(r, user.ContextKey).(*user.User)
	if !u.HasPermission(user.CreateNewUserPermission) {
		w.WriteHeader(http.StatusForbidden)
		h.write(w, forbiddenError)

		return
	}

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		h.logger.Error("Error read body", "error", readErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

		return
	}

	values, parseErr := url.ParseQuery(string(body))
	if parseErr != nil {
		h.logger.Warn("Error parse body", "error", parseErr.Error())
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	stringID := strings.TrimSpace(values.Get("id"))
	ID, convertErr := strconv.ParseInt(stringID, 10, 64)
	if convertErr != nil {
		h.logger.Warn("Error parse user id", "error", convertErr.Error(), "id", stringID)
		w.WriteHeader(http.StatusBadRequest)
		h.write(w, inputInvalidError)

		return
	}

	if err := h.user.Delete(ID); err != nil {
		h.logger.Error("Error delete project", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		h.write(w, internalError)

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

	htl, checkErr := h.hotel.HotelByName(hotelName)
	if checkErr != nil {
		return nil, errors.Wrap(checkErr, "error check hotel exists")
	}
	if htl != nil {
		return htl, nil
	}

	htl = &hotel.Hotel{Name: hotelName}
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
