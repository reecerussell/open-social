package main

import (
	"log"
	"os"
	"os/signal"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/media"
	"github.com/reecerussell/open-social/media/gcp"
	"github.com/reecerussell/open-social/service/media/handler"
	"github.com/reecerussell/open-social/service/media/repository"
)

const (
	connectionStringVar = "CONNECTION_STRING"
	mediaBucketVar      = "MEDIA_BUCKET"
)

func main() {
	ctn := buildServices()

	createMedia := ctn.GetService("CreateMediaHandler").(*handler.CreateMediaHandler)
	getMediaContent := ctn.GetService("GetMediaContentHandler").(*handler.GetMediaContentHandler)

	app := core.NewApp("0.0.0.0:80")
	app.Post("/media", createMedia)
	app.Get("/media/content/{referenceID}", getMediaContent)
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

	ctn.AddService("MediaService", func(ctn *core.Container) interface{} {
		uploader := gcp.New(os.Getenv(mediaBucketVar))
		return uploader
	})

	ctn.AddService("CreateMediaHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("MediaRepository").(repository.MediaRepository)
		uploader := ctn.GetService("MediaService").(media.Service)

		return handler.NewCreateMediaHandler(repo, uploader)
	})

	ctn.AddService("GetMediaContentHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("MediaRepository").(repository.MediaRepository)
		downloader := ctn.GetService("MediaService").(media.Service)

		return handler.NewGetMediaContentHandler(repo, downloader)
	})

	return ctn
}
