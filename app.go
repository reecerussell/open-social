package core

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	addr       string
	middleware []Middleware
	r          *mux.Router
}

func NewApp(addr string) *App {
	return &App{
		addr: addr,
		r:    mux.NewRouter(),
	}
}

func (app *App) AddMiddleware(middleware Middleware) {
	app.middleware = append(app.middleware, middleware)
}

func (app *App) HealthCheck(h http.HandlerFunc) {
	app.r.HandleFunc("/health", h).Methods("GET")
}

func (app *App) Get(path string, h http.Handler) {
	app.r.Handle(path, h).Methods(http.MethodGet)
}

func (app *App) GetFunc(path string, h http.HandlerFunc) {
	app.r.HandleFunc(path, h).Methods(http.MethodGet)
}

func (app *App) Post(path string, h http.Handler) {
	app.r.Handle(path, h).Methods(http.MethodPost)
}

func (app *App) PostFunc(path string, h http.HandlerFunc) {
	app.r.HandleFunc(path, h).Methods(http.MethodPost)
}

func (app *App) Serve() {
	s := http.Server{
		Addr:    app.addr,
		Handler: ChainMiddleware(app.r, app.middleware...),
	}

	log.Printf("Listening on %s\n", app.addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}
