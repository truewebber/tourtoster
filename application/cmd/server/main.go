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
	//db, err := sql.Open("sqlite3", "/home/truewebber/web/tourtoster.truewebber.com/ttdb.sqlite")
	if err != nil {
		log.Fatal("error connect to db", "error", err.Error())
	}
	// ----------------------------------------------------------------
	err = templatesInit("/Users/truewebber/tourtoster/templates")
	//err = templatesInit("/home/truewebber/web/tourtoster.truewebber.com/app/templates")
	if err != nil {
		log.Fatal("Error start http server", "error", err.Error())
	}
	log.Debug("templates init")
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
	// ----------------------------- MAIN -----------------------------
	r.HandleFunc(handler.MainPageAuthorizationPath, handlers.MainAuthorizationPage).Methods(http.MethodGet)
	r.HandleFunc(handler.MainPageIndexPath, handlers.MainIndexPage).Methods(http.MethodGet)
	r.HandleFunc(handler.MainPageLogoutPath, handlers.MainLogoutPage).Methods(http.MethodGet)
	// --------------------------- MAIN API ---------------------------

	// ----------------------------------------------------------------
	ra := r.PathPrefix(handler.AdminPathPrefix).Subrouter()
	// ---------------------------- ADMIN -----------------------------
	//ra.HandleFunc(handler.AuthorizationAdminPagePath, handlers.AuthorizationAdminPage).Methods(http.MethodGet)
	ra.HandleFunc(handler.AdminPageIndexPath, handlers.AdminIndexPage).Methods(http.MethodGet)
	//ra.HandleFunc(handler.ServicesAdminPagePath, handlers.ServicesAdminPage).Methods(http.MethodGet)
	//ra.HandleFunc(handler.ProjectsAdminPagePath, handlers.ProjectsAdminPage).Methods(http.MethodGet)
	// -------------------------- ADMIN API ---------------------------
	// -------------------------- MIDDLEWARE --------------------------
	r.Use(middlewares.AuthMiddleware)
	// ----------------------------------------------------------------

	log.Debug("Starting server", host, port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Error("Error start http server", "error", err.Error())
	}
}

func templatesInit(templatePath string) error {
	filesName := []string{
		"parts/footer.gohtml",
		"parts/bottom.gohtml",
		"parts/header.gohtml",
		"parts/sidebar.gohtml",
		"parts/top.gohtml",
		// -- pages
		// --
		"admin-index.gohtml",
		"main-index.gohtml",
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
		handler.MainIndexTemplateName,
		handler.AdminIndexTemplateName,
		handler.MainPageAuthorizationTemplateName,
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
