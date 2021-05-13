package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/subosito/gotenv"

	"github.com/lorezi/boolang/controllers"
	"github.com/lorezi/boolang/middleware"
	"github.com/lorezi/boolang/pkg/metric"
	"github.com/lorezi/boolang/pkg/prometheus"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lorezi/boolang/docs"
)

// @title Boolang
// @version 1.0
// @description This is a CRUD application.
// @termsOfService http://swagger.io/terms/

// @contact.name Lawrence Onaulogho
// @contact.url https://github.com/lorezi/
// @contact.email lawrence[at][gmail][dot][com]

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /

func main() {

	gotenv.Load()

	r := mux.NewRouter()
	// Prometheus
	pro := prometheus.New(true)
	// metric
	mtr := metric.New(pro.Registry())
	// prometheus metric routes
	r.HandleFunc("/metrics", pro.Handler())

	bc := controllers.NewBookController()
	uc := controllers.NewUserController()
	pc := controllers.NewPermissionController()

	// metric controller
	mc := controllers.NewBalanceUpdate(mtr)

	subr := r.PathPrefix("/api").Subrouter()

	subr.HandleFunc("/balances", mc.Handle)

	// create a middleware for monitoring
	subr.Use(middleware.Monitoring)

	authr := subr.PathPrefix("/auth").Subrouter()
	authr.Use(middleware.Authentication)

	// book protected routes
	bkr := authr.PathPrefix("/api-books").Subrouter()
	bkr.Use(middleware.BookAuthorization)
	// bkr.HandleFunc("/home",  bc.HomePage).Methods("GET")
	bkr.HandleFunc("/books", bc.GetBooks).Methods("GET").Queries("limit", "{limit:[0-9]+}", "page", "{page:[0-9]+}")
	bkr.HandleFunc("/books/{id}", bc.GetBook).Methods("GET")
	bkr.HandleFunc("/books", bc.AddBook).Methods("POST")
	bkr.HandleFunc("/books/{id}", bc.UpdateBook).Methods("PATCH")
	bkr.HandleFunc("/books/{id}", bc.DeleteBook).Methods("DELETE")

	// user protected routes
	ur := authr.PathPrefix("/api-users").Subrouter()
	ur.Use(middleware.UserAuthorization)

	ur.HandleFunc("/users/{id}", uc.GetUser).Methods("GET")
	ur.HandleFunc("/users/{id}", uc.UpdateUser).Methods("PATCH")
	ur.HandleFunc("/users", uc.GetUsers).Methods("GET").Queries("limit", "{limit:[0-9]+}", "page", "{page:[0-9]+}")

	// permission protected routes
	pr := authr.PathPrefix("/api-permissions").Subrouter()
	pr.Use(middleware.PermissionAuthorization)
	pr.HandleFunc("/permissions/{id}", pc.GetPermission).Methods("GET")
	pr.HandleFunc("/permissions", pc.CreatePermission).Methods("POST")
	pr.HandleFunc("/permissions", pc.GetPermissions).Methods("GET")

	// unprotected routes
	subr.HandleFunc("/users/login", uc.Login).Methods("POST")
	subr.HandleFunc("/users/signup", uc.CreateUser).Methods("POST")
	// subr.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)

	subr.Handle("/home", handlers.LoggingHandler(os.Stdout, http.HandlerFunc(bc.HomePage))).Methods("GET")
	r.Handle("/favicon.ico", http.NotFoundHandler())
	r.PathPrefix("/").Handler(httpSwagger.WrapHandler)

	handler := cors.Default().Handler(r)
	srv := &http.Server{
		Handler: handler,
		// Addr:         ":3000",
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
