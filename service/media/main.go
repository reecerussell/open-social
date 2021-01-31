package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/media/handler"
	"github.com/reecerussell/open-social/service/media/repository"
)

const (
	connectionStringVar = "CONNECTION_STRING"
)

func main() {
	ctn := buildServices()

	createMedia := ctn.GetService("CreateMediaHandler").(*handler.CreateMediaHandler)

	app := core.NewApp("0.0.0.0:80")
	app.Post("/media", createMedia)
	app.HealthCheck(core.HealthCheckHandler)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

func buildServices() *core.Container {
	ctn := core.NewContainer()

	ctn.AddService("MediaRepository", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		return repository.NewMediaRepository(url)
	})

	ctn.AddService("CreateMediaHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("MediaRepository").(repository.MediaRepository)

		return handler.NewCreateMediaHandler(repo)
	})

	return ctn
}
