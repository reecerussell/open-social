package main

import (
	"log"
	"os"
	"os/signal"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/media"
	"github.com/reecerussell/open-social/cmd/media-download/handler"
)

const (
	mediaAPIVar = "MEDIA_API_URL"
)

func main() {
	ctn := buildServices()

	downloadHandler := ctn.GetService("DownloadHandler").(*handler.DownloadHandler)

	app := core.NewApp()
	app.AddMiddleware(core.NewLoggingMiddleware())

	app.Get(app.HealthPath, core.HealthCheckHandler(app.HealthChecks))
	app.Get("/{referenceID}", downloadHandler)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

func buildServices() *core.Container {
	ctn := core.NewContainer()

	ctn.AddService("MediaClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(mediaAPIVar)
		client := media.New(url)
		return client
	})

	ctn.AddService("DownloadHandler", func(ctn *core.Container) interface{} {
		client := ctn.GetService("MediaClient").(media.Client)
		h := handler.NewDownloadHandler(client)
		return h
	})

	return ctn
}
