package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"
	"github.com/pkg/errors"

	"tourtoster/handler"
	"tourtoster/middleware"
	"tourtoster/token"
	tokenRepo "tourtoster/token/repository"
	userRepo "tourtoster/user/repository"
)

var (
	port       string
	host       string
	configPath string

	templates = make(map[string]*template.Template)
)

const (
	defaultConfigPath = ""
)

func init() {
	flag.StringVar(&port, "port", "9000", "")
	flag.StringVar(&host, "host", "localhost", "")
	flag.StringVar(&configPath, "config-path", defaultConfigPath, "")

	flag.Parse()
}

func main() {
	//cfg, err := config.New("prod", configPath)
	//if err != nil {
	//	log.Error("Error get config", "error", err.Error())
	//
	//	return
	//}
	// ----------------------------------------------------------------
	db, err := sql.Open("sqlite3", "/Users/truewebber/tourtoster/ttdb.sqlite")
	if err != nil {
		log.Fatal("error connect to db", "error", err.Error())
	}
	// ----------------------------------------------------------------
	err = templatesInit("/Users/truewebber/tourtoster/templates")
	if err != nil {
		log.Fatal("Error start http server", "error", err.Error())
	}
	log.Debug("templates init")
	// ----------------------------------------------------------------
	tokenR := tokenRepo.NewMemory()
	tokenR.Save(&token.Token{
		Token:  "blah",
		UserID: 1,
	})
	userR := userRepo.NewPostgres(db)
	//serviceR := serviceRepo.NewPostgres(db)
	//projectR := projectRepo.NewPostgres(db)
	//
	handlers := handler.New(&handler.Config{
		User:      userR,
		Token:     tokenR,
		Templates: templates,
	})
	middlewares := middleware.New(tokenR, userR)
	// ----------------------------------------------------------------

	// ---------------------------- ROUTER ----------------------------
	r := mux.NewRouter()
	r.HandleFunc("/{path:.*}/", func(w http.ResponseWriter, r *http.Request) {
		newPath := r.URL.Path[:len(r.URL.Path)-1]
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
	})
	// ----------------------------- MAIN -----------------------------
	r.HandleFunc(handler.MainAuthorizationPagePath, handlers.MainAuthorizationPage).Methods(http.MethodGet)
	r.HandleFunc(handler.MainPagePath, handlers.MainPage).Methods(http.MethodGet)
	// ---------------------------- ADMIN -----------------------------
	//ra := r.PathPrefix(handler.AdminPathPrefix).Subrouter()
	// ----------------------------- PAGE -----------------------------
	//tw.HandleFunc(handler.AuthorizationAdminPagePath, handlers.AuthorizationAdminPage).Methods(http.MethodGet)
	//tw.HandleFunc(handler.LogoutAdminPagePath, handlers.LogoutAdminPage).Methods(http.MethodGet)
	//tw.HandleFunc(handler.MainAdminPagePath, handlers.MainAdminPage).Methods(http.MethodGet)
	//tw.HandleFunc(handler.ServicesAdminPagePath, handlers.ServicesAdminPage).Methods(http.MethodGet)
	//tw.HandleFunc(handler.ProjectsAdminPagePath, handlers.ProjectsAdminPage).Methods(http.MethodGet)
	//// ----------------------------- API ------------------------------
	//tw.HandleFunc(handler.AuthorizationAdminApiPath, handlers.AuthorizationAdminApi).Methods(http.MethodPost)
	////
	//tw.HandleFunc(handler.LanguageAdminApiPath, handlers.LanguageAdminApiPost).Methods(http.MethodPost)
	//tw.HandleFunc(handler.LanguageAdminApiPath, handlers.LanguageAdminApiDelete).Methods(http.MethodDelete)
	////
	//tw.HandleFunc(handler.ServiceAdminApiPath, handlers.ServiceAdminApiPost).Methods(http.MethodPost)
	//tw.HandleFunc(handler.ServiceAdminApiPath, handlers.ServiceAdminApiDelete).Methods(http.MethodDelete)
	////
	//tw.HandleFunc(handler.ProjectAdminApiPath, handlers.ProjectAdminApiPost).Methods(http.MethodPost)
	//tw.HandleFunc(handler.ProjectAdminApiPath, handlers.ProjectAdminApiDelete).Methods(http.MethodDelete)
	////
	//tw.HandleFunc(handler.LocalizationAdminApiPath, handlers.LocalizationAdminApiGet).Methods(http.MethodGet)
	//tw.HandleFunc(handler.LocalizationAdminApiPath, handlers.LocalizationAdminApiPost).Methods(http.MethodPost)
	////
	//tw.HandleFunc(handler.FilesAdminApiPath, handlers.FilesAdminApiPost).Methods(http.MethodPost)
	//tw.HandleFunc(handler.FilesAdminApiPath, handlers.FilesAdminApiDelete).Methods(http.MethodDelete)
	//// -------------------------- MIDDLEWARE --------------------------
	r.Use(middlewares.AuthMiddleware)
	// ----------------------------------------------------------------

	log.Debug("Starting server", host, port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Error("Error start http server", "error", err.Error())
	}
}

func templatesInit(templatePath string) error {
	filesName := []string{
		//"parts/footer.gohtml",
		//"parts/bottom.gohtml",
		//"parts/header.gohtml",
		//"parts/sidebar.gohtml",
		//"parts/top.gohtml",
		// -- pages
		// --
		"index.gohtml",
		"authorization.gohtml",
	}

	pathes := make([]string, 0, len(filesName))
	for _, fileName := range filesName {
		pathes = append(pathes, templatePath+"/"+fileName)
	}

	tmpls, err := template.New("blah").ParseFiles(pathes...)
	if err != nil {
		return err
	}

	templateNames := []string{
		handler.MainTemplateName,
		handler.MainAuthorizationTemplateName,
	}
	for _, n := range templateNames {
		t := tmpls.Lookup(n)
		if t == nil {
			return errors.Errorf("Template `%s` not found", n)
		}
		templates[n] = t
	}

	return nil
}
