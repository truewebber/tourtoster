package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/truewebber/tourtoster/conn"
	"github.com/truewebber/tourtoster/currency"
	currencyRepo "github.com/truewebber/tourtoster/currency/repository"
	"github.com/truewebber/tourtoster/handler"
	hotelRepo "github.com/truewebber/tourtoster/hotel/repository"
	"github.com/truewebber/tourtoster/log"
	"github.com/truewebber/tourtoster/mail"
	mailRepo "github.com/truewebber/tourtoster/mail/repository"
	"github.com/truewebber/tourtoster/middleware"
	"github.com/truewebber/tourtoster/token"
	tokenRepo "github.com/truewebber/tourtoster/token/repository"
	"github.com/truewebber/tourtoster/tour"
	tourRepo "github.com/truewebber/tourtoster/tour/repository"
	userRepo "github.com/truewebber/tourtoster/user/repository"
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
	logger := log.NewZap()
	defer func() {
		if err := logger.Close(); err != nil {
			println("error close logger", err.Error())
		}
	}()

	db, err := conn.NewConn(dbFilePath)
	if err != nil {
		println("error create db connection")
		panic(err)
	}

	logger.Info("connection to db established", "db", dbFilePath)
	// -----------------------------------------------------------------------------------------------------------------
	tokenR := tokenRepo.NewMemory()
	_ = tokenR.Save(&token.Token{
		Token:  "blah",
		UserID: 1,
	})
	userR := userRepo.NewSQLite(db, logger)
	hotelR := hotelRepo.NewSQLite(db, logger)
	tourR := tourRepo.NewSQLite(db, userR, logger)

	if err := initCurrencies(db, logger); err != nil {
		println("error init currencies")
		panic(err)
	}
	if err := initFeatures(tourR); err != nil {
		println("error init features")
		panic(err)
	}
	// -----------------------------------------------------------------------------------------------------------------
	mailer := newMailer(logger)
	logger.Info("Init mailer", "_", mailer.Name())
	// -----------------------------------------------------------------------------------------------------------------
	handlers, handlersErr := handler.New(&handler.Config{
		User:          userR,
		Token:         tokenR,
		Tour:          tourR,
		Hotel:         hotelR,
		Mailer:        mailer,
		TemplatesPath: templatePath,
		Logger:        logger.With("module", "handlers"),
	})
	if handlersErr != nil {
		println("error init handlers")
		panic(handlersErr)
	}
	logger.Info("templates init", "path", templatePath)
	// -----------------------------------------------------------------------------------------------------------------
	middlewares := middleware.New(tokenR, userR, hotelR, logger)
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
	rc.HandleFunc(handler.ConsoleGTPath, handlers.ConsoleGTPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleGTEditPath, handlers.ConsoleGTEditPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleGTEditFAQPath, handlers.ConsoleGTEditFAQPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsolePTPath, handlers.ConsolePTPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsolePTEditPath, handlers.ConsolePTEditPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsolePTEditFAQPath, handlers.ConsolePTEditFAQPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserPath, handlers.ConsoleUserPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserBillingPath, handlers.ConsoleUserBillingPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleUserProfilePath, handlers.ConsoleUserProfilePage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleContactPath, handlers.ConsoleContactPage).Methods(http.MethodGet)
	rc.HandleFunc(handler.ConsoleAboutPath, handlers.ConsoleAboutPage).Methods(http.MethodGet)
	// -----------------------------------------------------------------------------------------------------------------
	rca := r.PathPrefix(handler.ApiPathPrefix).Subrouter()
	// --------------------------------------------------- MAIN API ----------------------------------------------------
	rca.HandleFunc(handler.ForgetApiPath, handlers.ApiForget).Methods(http.MethodPost)
	rca.HandleFunc(handler.RegistrationApiPath, handlers.ApiRegistration).Methods(http.MethodPost)
	rca.HandleFunc(handler.AuthorizationApiPath, handlers.AuthorizationApi).Methods(http.MethodPost)
	//
	rca.HandleFunc(handler.UserApiPath, handlers.ApiUserCreate).Methods(http.MethodPost)
	rca.HandleFunc(handler.UserApiPath, handlers.ApiUserDelete).Methods(http.MethodDelete)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelList).Methods(http.MethodGet)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelCreate).Methods(http.MethodPost)
	rca.HandleFunc(handler.HotelApiPath, handlers.ApiHotelDelete).Methods(http.MethodDelete)
	rca.HandleFunc(handler.TourApiPath, handlers.ApiTourCreate).Methods(http.MethodPost)
	rca.HandleFunc(handler.TourApiPath, handlers.ApiTourDelete).Methods(http.MethodDelete)
	// -------------------------------------------------- MIDDLEWARE ---------------------------------------------------
	r.Use(middlewares.LogMiddleware)
	rc.Use(middlewares.PageAuthMiddleware)
	rca.Use(middlewares.APIAuthMiddleware)
	// -----------------------------------------------------------------------------------------------------------------
	// ---------------------------------------------------- SERVER -----------------------------------------------------
	logger.Info("Starting server", host, port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		logger.Error("Error start http server", "error", err.Error())
	}
	// -----------------------------------------------------------------------------------------------------------------
}

func initFeatures(repo tour.Repository) error {
	ff, err := repo.Features()
	if err != nil {
		return err
	}

	tour.FeaturesByType[tour.PrivateType] = make([]tour.Feature, 0)
	tour.FeaturesByType[tour.GroupType] = make([]tour.Feature, 0)

	for i := 0; i < len(ff); i++ {
		tour.FeaturesByID[ff[i].ID] = ff[i]
		tour.FeaturesByType[ff[i].TourType] = append(tour.FeaturesByType[ff[i].TourType], ff[i])
	}

	return nil
}

func initCurrencies(db *sql.DB, logger log.Logger) error {
	repo := currencyRepo.NewSQLite(db, logger)
	values, err := repo.List(currency.USDName, currency.EURName)
	if err != nil {
		return err
	}

	if value, ok := values[currency.USDName]; ok {
		currency.USD = value
	}

	if value, ok := values[currency.EURName]; ok {
		currency.EUR = value
	}

	return nil
}

func newMailer(logger log.Logger) mail.Mailer {
	u := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_PASSWORD")

	switch os.Getenv("MAIL_SERVICE") {
	case mailRepo.GMailName:
		return mailRepo.NewGMail(u, pass)
	case mailRepo.YandexName:
		return mailRepo.NewYandex(u, pass)
	default:
		return mailRepo.NewNull(logger)
	}
}
