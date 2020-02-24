package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mgutz/logxi/v1"

	"tourtoster/handler"
	hotelRepo "tourtoster/hotel/repository"
	"tourtoster/mail"
	mailRepo "tourtoster/mail/repository"
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
	// -----------------------------------------------------------------------------------------------------------------
	tokenR := tokenRepo.NewMemory()
	_ = tokenR.Save(&token.Token{
		Token:  "blah",
		UserID: 1,
	})
	userR := userRepo.NewPostgres(db)
	hotelR := hotelRepo.NewPostgres(db)
	// -----------------------------------------------------------------------------------------------------------------
	mailer := newMailer()
	log.Debug("Init mailer", "_", mailer.Name())
	// -----------------------------------------------------------------------------------------------------------------
	handlers, handlersErr := handler.New(&handler.Config{
		User:          userR,
		Token:         tokenR,
		Hotel:         hotelR,
		Mailer:        mailer,
		TemplatesPath: templatePath,
	})
	if handlersErr != nil {
		println("error init handlers")
		panic(handlersErr)
	}
	log.Debug("templates init", "path", templatePath)
	// -----------------------------------------------------------------------------------------------------------------
	middlewares := middleware.New(tokenR, userR, hotelR)
	// -----------------------------------------------------------------------------------------------------------------

	// ---------------------------------------------------- ROUTER -----------------------------------------------------
	r := mux.NewRouter()
	r.HandleFunc("/{path:.*}/", func(w http.ResponseWriter, r *http.Request) {
		newPath := r.URL.Path[:len(r.URL.Path)-1]
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
	})
	r.HandleFunc(handler.LandingPageIndexPath, handlers.LandingIndexPage).Methods(http.MethodGet)
	// -----------------------------------------------------------------------------------------------------------------
	rc := r.PathPrefix(handler.ConsolePathPrefix).Subrouter()
	// ----------------------------------------------------- MAIN ------------------------------------------------------
	rc.HandleFunc(handler.ConsoleRegistrationPath, handlers.ConsoleRegistrationPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleAuthorizationPath, handlers.ConsoleAuthorizationPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleSignoutPath, handlers.ConsoleSignoutPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleIndexPath, handlers.ConsoleGTPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleGPPath, handlers.ConsoleGTPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsolePTPath, handlers.ConsolePTPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserPath, handlers.ConsoleUserPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserBillingPath, handlers.ConsoleUserBillingPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserProfilePath, handlers.ConsoleUserProfilePage).Methods(http.MethodGet)
	// -----------------------------------------------------------------------------------------------------------------
	rca := r.PathPrefix(handler.ApiPathPrefix).Subrouter()
	// --------------------------------------------------- MAIN API ----------------------------------------------------
	rca.HandleFunc(handler.AuthorizationApiPath, handlers.AuthorizationApi).Methods(http.MethodPost)
	rca.HandleFunc(handler.UserApiPath, handlers.ApiUserCreate).Methods(http.MethodPost)
	rca.HandleFunc(handler.UserApiPath, handlers.ApiUseDelete).Methods(http.MethodDelete)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelList).Methods(http.MethodGet)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelCreate).Methods(http.MethodPost)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelDelete).Methods(http.MethodDelete)
	// -------------------------------------------------- MIDDLEWARE ---------------------------------------------------
	rc.Use(middlewares.PageAuthMiddleware)
	rca.Use(middlewares.APIAuthMiddleware)
	// -----------------------------------------------------------------------------------------------------------------

	// ---------------------------------------------------- SERVER -----------------------------------------------------
	log.Debug("Starting server", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Error("Error start http server", "error", err.Error())
	}
	// -----------------------------------------------------------------------------------------------------------------
}

func newMailer() mail.Mailer {
	u := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_PASSWORD")

	switch os.Getenv("MAIL_SERVICE") {
	case mailRepo.GMailName:
		return mailRepo.NewGMail(u, pass)
	case mailRepo.YandexName:
		return mailRepo.NewYandex(u, pass)
	default:
		return mailRepo.NewNull()
	}
}
