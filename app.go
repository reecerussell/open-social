package core

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/reecerussell/open-social/util"
)

// PortEnvVar is the name of the environment variable user to read the port value.
const PortEnvVar = "PORT"

// DefaultPort is the default port number to use if the PortEnvVar is not set.
const DefaultPort = "8080"

// App pulls together components of a HTTP api, such as middleware, health checks
// and routing.
type App struct {
	addr         string
	middleware   []Middleware
	HealthChecks []HealthCheck
	router       *mux.Router

	HealthPath string
}

// NewApp returns a new instance of App.
func NewApp() *App {
	port := util.ReadEnv(PortEnvVar, DefaultPort)

	return &App{
		addr:         fmt.Sprintf("0.0.0.0:%s", port),
		router:       mux.NewRouter(),
		HealthChecks: []HealthCheck{},
		HealthPath:   "/health",
	}
}

// AddMiddleware adds middleware to the App.
func (app *App) AddMiddleware(middleware Middleware) {
	app.middleware = append(app.middleware, middleware)
}

// AddHealthCheck adds a health check to the App.
func (app *App) AddHealthCheck(healthCheck HealthCheck) {
	app.HealthChecks = append(app.HealthChecks, healthCheck)
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
	app.router.Handle(app.HealthPath, HealthCheckHandler(app.HealthChecks))

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
