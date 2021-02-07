package core

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// App pulls together components of a HTTP api, such as middleware, health checks
// and routing.
type App struct {
	addr         string
	middleware   []Middleware
	healthChecks []HealthCheck
	router       *mux.Router

	HealthPath string
}

// NewApp returns a new instance of App.
func NewApp(addr string) *App {
	return &App{
		addr:         addr,
		router:       mux.NewRouter(),
		healthChecks: []HealthCheck{},
		HealthPath:   "/health",
	}
}

// AddMiddleware adds middleware to the App.
func (app *App) AddMiddleware(middleware Middleware) {
	app.middleware = append(app.middleware, middleware)
}

// AddHealthCheck adds a health check to the App.
func (app *App) AddHealthCheck(healthCheck HealthCheck) {
	app.healthChecks = append(app.healthChecks, healthCheck)
}

func (app *App) Get(path string, h http.Handler) {
	app.router.Handle(path, h).Methods(http.MethodGet)
}

func (app *App) GetFunc(path string, h http.HandlerFunc) {
	app.router.HandleFunc(path, h).Methods(http.MethodGet)
}

func (app *App) Post(path string, h http.Handler) {
	app.router.Handle(path, h).Methods(http.MethodPost)
}

func (app *App) PostFunc(path string, h http.HandlerFunc) {
	app.router.HandleFunc(path, h).Methods(http.MethodPost)
}

func (app *App) Serve() {
	app.router.Handle(app.HealthPath, HealthCheckHandler(app.healthChecks))

	s := http.Server{
		Addr:    app.addr,
		Handler: ChainMiddleware(app.router, app.middleware...),
	}

	log.Printf("Listening on %s\n", app.addr)
	err := s.ListenAndServe()
	if err != nil {
		log.Panicln(err)
	}
}
