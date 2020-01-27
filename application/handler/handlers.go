package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mgutz/logxi/v1"
	"golang.org/x/crypto/bcrypt"

	"tourtoster/hotel"
	"tourtoster/mail"
	"tourtoster/token"
	"tourtoster/user"
)

type (
	Handlers struct {
		user      user.Repository
		token     token.Repository
		hotel     hotel.Repository
		templates map[string]*template.Template
		mailer    mail.Mailer
	}

	Config struct {
		User      user.Repository
		Token     token.Repository
		Hotel     hotel.Repository
		Templates map[string]*template.Template
		Mailer    mail.Mailer
	}

	respError struct {
		Error string `json:"error"`
	}
)

var (
	internalError = respError{
		Error: "Server Error",
	}

	forbiddenError = respError{
		Error: "Access Denied",
	}

	inputInvalidError = respError{
		Error: "input data is invalid",
	}
)

const (
	ConsolePathPrefix = "/console"
	ApiPathPrefix     = ConsolePathPrefix + "/api"
)

func New(cfg *Config) *Handlers {
	return &Handlers{
		user:      cfg.User,
		token:     cfg.Token,
		hotel:     cfg.Hotel,
		templates: cfg.Templates,
		mailer:    cfg.Mailer,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func write(w http.ResponseWriter, obj interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(obj)
	if err != nil {
		log.Error("Error write response", "error", err.Error())
	}
}
