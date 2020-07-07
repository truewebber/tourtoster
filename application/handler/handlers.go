package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"tourtoster/hotel"
	"tourtoster/log"
	"tourtoster/mail"
	"tourtoster/token"
	"tourtoster/tour"
	"tourtoster/user"
)

type (
	Handlers struct {
		user      user.Repository
		token     token.Repository
		hotel     hotel.Repository
		tour      tour.Repository
		templates map[string]*template.Template
		mailer    mail.Mailer
		logger    log.Logger
	}

	Config struct {
		User          user.Repository
		Token         token.Repository
		Tour          tour.Repository
		Hotel         hotel.Repository
		Mailer        mail.Mailer
		TemplatesPath string
		Logger        log.Logger
	}
)

type (
	respOK struct {
		Message string `json:"message"`
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

func New(cfg *Config) (*Handlers, error) {
	templates, err := templatesInit(cfg.TemplatesPath)
	if err != nil {
		return nil, err
	}

	return &Handlers{
		user:      cfg.User,
		token:     cfg.Token,
		tour:      cfg.Tour,
		hotel:     cfg.Hotel,
		templates: templates,
		mailer:    cfg.Mailer,
		logger:    cfg.Logger,
	}, nil
}

func templatesInit(templatePath string) (map[string]*template.Template, error) {
	filesName := []string{
		"parts/footer.gohtml",
		"parts/header/header.gohtml",
		"parts/header/header-mobile.gohtml",
		"parts/header/header-dropdown-user-menu.gohtml",
		// --
		"landing-index.gohtml",
		"console-authorization.gohtml",
		"console-registration.gohtml",
		"console-index.gohtml",
		"console-group_tours.gohtml",
		"console-group_tours-view_all.gohtml",
		"console-group_tours-edit.gohtml",
		"console-group_tours-edit-faq.gohtml",
		"console-private_tours.gohtml",
		"console-private_tours-view_all.gohtml",
		"console-private_tours-edit.gohtml",
		"console-private_tours-edit-faq.gohtml",
		"console-user-profile.gohtml",
		"console-user-billing.gohtml",
		"console-index.gohtml",
		"console-user.gohtml",
		"console-contact.gohtml",
		"console-about.gohtml",
	}

	pathes := make([]string, 0, len(filesName))
	for _, fileName := range filesName {
		pathes = append(pathes, templatePath+"/"+fileName)
	}

	tmpls, err := template.New("blah").Funcs(template.FuncMap{
		"UserShortName": user.ShortName,
	}).ParseFiles(pathes...)
	if err != nil {
		return nil, err
	}

	templateNames := []string{
		LandingIndexTemplateName,
		ConsoleAuthorizationTemplateName,
		ConsoleRegistrationTemplateName,
		ConsoleIndexTemplateName,
		ConsoleGTTemplateName,
		ConsoleGTViewAllSubTemplateName,
		ConsoleGTEditSubTemplateName,
		ConsoleGTEditFAQSubTemplateName,
		ConsolePTTemplateName,
		ConsolePTViewAllSubTemplateName,
		ConsolePTEditSubTemplateName,
		ConsolePTEditFAQSubTemplateName,
		ConsoleUserTemplateName,
		ConsoleUserBillingTemplateName,
		ConsoleUserProfileTemplateName,
		ConsoleContactTemplateName,
		ConsoleAboutTemplateName,
	}
	templates := make(map[string]*template.Template)
	for _, n := range templateNames {
		t := tmpls.Lookup(n)
		if t == nil {
			return nil, errors.Errorf("Template `%s` not found", n)
		}
		templates[n] = t
	}

	return templates, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func (h *Handlers) write(w http.ResponseWriter, obj interface{}) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(obj)
	if err != nil {
		h.logger.Error("Error write response", "error", err.Error())
	}
}
