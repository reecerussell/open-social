package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/posts/handler"
	"github.com/reecerussell/open-social/service/posts/repository"
)

const (
	connectionStringVar = "CONNECTION_STRING"
	usersAPIVar         = "USERS_API_URL"
)

func main() {
	ctn := buildServices()

	createPost := ctn.GetService("CreatePostHandler").(*handler.CreatePostHandler)

	app := core.NewApp("0.0.0.0:80")
	app.Post("/posts", createPost)
	app.HealthCheck(core.HealthCheckHandler)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

func buildServices() *core.Container {
	ctn := core.NewContainer()

	ctn.AddService("PostRepository", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		return repository.NewPostRepository(url)
	})

	ctn.AddSingleton("UserClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(usersAPIVar)
		return users.New(url)
	})

	ctn.AddService("CreatePostHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("PostRepository").(repository.PostRepository)
		client := ctn.GetService("UserClient").(users.Client)

		return handler.NewCreatePostHandler(repo, client)
	})

	return ctn
}
