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
	hotelRepo "tourtoster/hotel/repository"
	"tourtoster/middleware"
	"tourtoster/token"
	tokenRepo "tourtoster/token/repository"
	userRepo "tourtoster/user/repository"
)

var (
	port         string
	host         string
	templatePath string
	dbFilePath   string

	templates = make(map[string]*template.Template)
)

func init() {
	flag.StringVar(&port, "port", "9000", "")
	flag.StringVar(&host, "host", "localhost", "")
	flag.StringVar(&templatePath, "template-path", "/Users/truewebber/tourtoster/templates", "")
	flag.StringVar(&dbFilePath, "db", "/Users/truewebber/tourtoster/ttdb.sqlite", "")

	flag.Parse()
}

func main() {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		println("error connect to db")
		panic(err)
	}
	if err := db.Ping(); err != nil {
		println("ping error")
		panic(err)
	}
	log.Debug("connection to db established", "db", dbFilePath)
	// ----------------------------------------------------------------
	if err := templatesInit(templatePath); err != nil {
		println("error init template")
		panic(err)
	}
	log.Debug("templates init", "path", templatePath)
	// ----------------------------------------------------------------
	tokenR := tokenRepo.NewMemory()
	_ = tokenR.Save(&token.Token{
		Token:  "blah",
		UserID: 1,
	})
	userR := userRepo.NewPostgres(db)
	hotelR := hotelRepo.NewPostgres(db)
	//
	handlers := handler.New(&handler.Config{
		User:      userR,
		Token:     tokenR,
		Templates: templates,
	})
	middlewares := middleware.New(tokenR, userR, hotelR)
	// ----------------------------------------------------------------

	// ---------------------------- ROUTER ----------------------------
	r := mux.NewRouter()
	r.HandleFunc("/{path:.*}/", func(w http.ResponseWriter, r *http.Request) {
		newPath := r.URL.Path[:len(r.URL.Path)-1]
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
	})
	r.HandleFunc(handler.LandingPageIndexPath, handlers.LandingIndexPage).Methods(http.MethodGet)

	// ----------------------------------------------------------------
	rc := r.PathPrefix(handler.ConsolePathPrefix).Subrouter()
	// ----------------------------- MAIN -----------------------------
	rc.HandleFunc(handler.ConsoleAuthorizationPath, handlers.ConsoleAuthorizationPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleIndexPath, handlers.ConsoleIndexPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleLogoutPath, handlers.ConsoleLogoutPage).Methods(http.MethodGet)
	// --------------------------- MAIN API ---------------------------

	//// ----------------------------------------------------------------
	//ra := r.PathPrefix(handler.AdminPathPrefix).Subrouter()
	//// ---------------------------- ADMIN -----------------------------
	////ra.HandleFunc(handler.AuthorizationAdminPagePath, handlers.AuthorizationAdminPage).Methods(http.MethodGet)
	//ra.HandleFunc(handler.AdminPageIndexPath, handlers.AdminIndexPage).Methods(http.MethodGet)
	////ra.HandleFunc(handler.ServicesAdminPagePath, handlers.ServicesAdminPage).Methods(http.MethodGet)
	////ra.HandleFunc(handler.ProjectsAdminPagePath, handlers.ProjectsAdminPage).Methods(http.MethodGet)
	//// -------------------------- ADMIN API ---------------------------
	//// -------------------------- MIDDLEWARE --------------------------
	rc.Use(middlewares.AuthMiddleware)
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
		// --
		"landing-index.gohtml",
		"console-authorization.gohtml",
		"console-index.gohtml",
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
		handler.LandingIndexTemplateName,
		handler.ConsoleAuthorizationTemplateName,
		handler.ConsoleIndexTemplateName,
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
