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

	"tourtoster/hotel"
	"tourtoster/user"
)

const (
	UserApiPath = "/user"
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

	var ID int64
	if stringID := strings.TrimSpace(values.Get("id")); len(stringID) > 0 {
		var err error
		ID, err = strconv.ParseInt(stringID, 10, 64)
		if err != nil {
			log.Error("Error parse project id", "error", err.Error(), "id", stringID)
			w.WriteHeader(http.StatusBadRequest)
			write(w, inputInvalidError)

			return
		}
	}

	email := html.EscapeString(strings.TrimSpace(values.Get("email")))

	firstName := html.EscapeString(strings.TrimSpace(values.Get("first_name")))
	secondName := html.EscapeString(strings.TrimSpace(values.Get("second_name")))
	lastName := html.EscapeString(strings.TrimSpace(values.Get("last_name")))

	phone := html.EscapeString(strings.TrimSpace(values.Get("phone")))
	note := html.EscapeString(strings.TrimSpace(values.Get("note")))

	// hotel:start
	hotelName := html.EscapeString(strings.TrimSpace(values.Get("hotel_name")))
	hotelIDStr := html.EscapeString(strings.TrimSpace(values.Get("hotel_id")))
	hotelID, parseHotelIDErr := strconv.ParseInt(hotelIDStr, 10, 64)
	if parseHotelIDErr != nil {
		log.Error("Error parse hotel id", "error", parseHotelIDErr.Error(), "hotel_id", hotelIDStr)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}

	var (
		hotelObj *hotel.Hotel
		hotelErr error
	)
	if hotelID == 0 {
		hotelObj, hotelErr = h.hotel.Hotel(hotelID)
	} else if hotelName != "" {
		hotelObj = &hotel.Hotel{
			Name: hotelName,
		}
		hotelErr = h.hotel.Save(hotelObj)
	}
	if hotelErr != nil {
		log.Error("Error get hotel object", "error", hotelErr.Error(),
			"hotel_id", hotelIDStr, "hotel_name", hotelName)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}
	// hotel:end

	// permissions:start
	var permissions user.Permission
	for _, p := range values["permission[]"] {
		p := html.EscapeString(strings.TrimSpace(p))
		if p == "" {
			continue
		}

		pInt, err := strconv.Atoi(p)
		if err != nil {
			log.Error("Error parse permission id", "error", err.Error(), "permission", p)
			w.WriteHeader(http.StatusBadRequest)
			write(w, inputInvalidError)

			return
		}
		permission := user.Permission(pInt)

		if err := user.ValidationPermission(permission); err != nil {
			log.Error("Permission invalid", "error", err.Error(), "permission", p)
			w.WriteHeader(http.StatusBadRequest)
			write(w, inputInvalidError)

			return
		}

		permissions = permissions | permission
	}
	// permissions:end

	// status:start
	statusStr := html.EscapeString(strings.TrimSpace(values.Get("status")))
	statusInt, parseStatusErr := strconv.Atoi(statusStr)
	if parseStatusErr != nil {
		log.Error("Error parse status id", "error", parseStatusErr.Error(), "status", statusStr)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}
	status := user.Status(statusInt)

	if err := user.ValidationStatus(status); err != nil {
		log.Error("Status invalid", "error", err.Error(), "status", statusStr)
		w.WriteHeader(http.StatusBadRequest)
		write(w, inputInvalidError)

		return
	}
	// status:end

	password := html.EscapeString(strings.TrimSpace(values.Get("password")))
	if password == "" {
		password = RandString(10)
	}
	passwordHash, hashErr := HashPassword(password)
	if hashErr != nil {
		log.Error("Error hash password", "error", hashErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}

	if values.Get("send_mail") == "1" {
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
		Hotel:        hotelObj,
		Note:         note,
		Email:        email,
		Phone:        phone,
		Status:       status,
		Permissions:  permissions,
		PasswordHash: passwordHash,
		Token:        nil,
	}

	if err := h.user.Save(newUser); err != nil {
		log.Error("Error create/save user", "error", err.Error(), "id", ID)
		w.WriteHeader(http.StatusInternalServerError)
		write(w, internalError)

		return
	}
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
