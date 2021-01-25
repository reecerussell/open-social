package main

import (
	"crypto"
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/gojwt/rsa"

	"github.com/reecerussell/open-social/client/auth"
	"github.com/reecerussell/open-social/client/posts"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/core"
	"github.com/reecerussell/open-social/service/backend/handler"
	"github.com/reecerussell/open-social/service/backend/middleware"
)

const (
	usersAPIVar       = "USERS_API_URL"
	authAPIVar        = "AUTH_API_URL"
	postsAPIVar       = "POSTS_API_URL"
	tokenPublicKeyVar = "TOKEN_PUBLIC_KEY"
)

func main() {
	ctn := buildServices()

	userHandler := ctn.GetService("UserHandler").(*handler.UserHandler)
	postHandler := ctn.GetService("PostHandler").(*handler.PostHandler)
	authMiddleware := ctn.GetService("AuthMiddleware").(*middleware.Authentication)

	app := core.NewApp("0.0.0.0:80")
	app.HealthCheck(core.HealthCheckHandler)
	app.AddMiddleware(core.NewLoggingMiddleware())
	app.AddMiddleware(authMiddleware)

	app.PostFunc("/users/register", userHandler.Register)
	app.PostFunc("/posts", postHandler.Create)
	app.GetFunc("/feed", postHandler.GetFeed)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

func buildServices() *core.Container {
	ctn := core.NewContainer()

	ctn.AddService("UserClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(usersAPIVar)
		client := users.New(url)
		return client
	})

	ctn.AddService("AuthClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(authAPIVar)
		client := auth.New(url)
		return client
	})

	ctn.AddService("PostClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(postsAPIVar)
		client := posts.New(url)
		return client
	})

	ctn.AddService("AuthMiddleware", func(ctn *core.Container) interface{} {
		path := os.Getenv(tokenPublicKeyVar)
		alg, err := rsa.NewFromFile(path, crypto.SHA256)
		if err != nil {
			panic(err)
		}

		h := middleware.NewAuthentication(alg)
		return h
	})

	ctn.AddService("UserHandler", func(ctn *core.Container) interface{} {
		client := ctn.GetService("UserClient").(users.Client)
		authClient := ctn.GetService("AuthClient").(auth.Client)
		h := handler.NewUserHandler(client, authClient)
		return h
	})

	ctn.AddService("PostHandler", func(ctn *core.Container) interface{} {
		client := ctn.GetService("PostClient").(posts.Client)
		h := handler.NewPostHandler(client)
		return h
	})

	return ctn
}
