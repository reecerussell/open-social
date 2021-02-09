package main

import (
	"log"
	"os"
	"os/signal"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/media/handler"
	"github.com/reecerussell/open-social/cmd/media/repository"
	"github.com/reecerussell/open-social/database"
	"github.com/reecerussell/open-social/media"
	"github.com/reecerussell/open-social/media/gcp"
)

const (
	connectionStringVar = "CONNECTION_STRING"
	mediaBucketVar      = "MEDIA_BUCKET"
)

func main() {
	ctn := buildServices()
	db := ctn.GetService("Database").(database.Database)

	createMedia := ctn.GetService("CreateMediaHandler").(*handler.CreateMediaHandler)
	getMediaContent := ctn.GetService("GetMediaContentHandler").(*handler.GetMediaContentHandler)

	app := core.NewApp()
	app.AddHealthCheck(database.NewHealthCheck(db))
	app.AddMiddleware(core.NewLoggingMiddleware())

	app.Post("/media", createMedia)
	app.Get("/media/content/{referenceID}", getMediaContent)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

func buildServices() *core.Container {
	ctn := core.NewContainer()

	ctn.AddSingleton("Database", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		db, err := database.New(url)
		if err != nil {
			panic(err)
		}

		return db
	})

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
