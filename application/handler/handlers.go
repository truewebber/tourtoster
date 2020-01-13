package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/mgutz/logxi/v1"

	"tourtoster/token"
	"tourtoster/user"
)

type (
	Handlers struct {
		user      user.Repository
		token     token.Repository
		templates map[string]*template.Template
	}

	Config struct {
		User      user.Repository
		Token     token.Repository
		Templates map[string]*template.Template
	}

	respError struct {
		Error string `json:"error"`
	}
)

var (
	internalError = respError{
		Error: "server error",
	}
)

func New(cfg *Config) *Handlers {
	return &Handlers{
		user:      cfg.User,
		token:     cfg.Token,
		templates: cfg.Templates,
	}
}

func write(w http.ResponseWriter, obj interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(obj)
	if err != nil {
		log.Error("Error write response", "error", err.Error())
	}
}
