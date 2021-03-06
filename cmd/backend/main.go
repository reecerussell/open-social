package main

import (
	"crypto"
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/gojwt/rsa"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/auth"
	"github.com/reecerussell/open-social/client/media"
	"github.com/reecerussell/open-social/client/posts"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/cmd/backend/handler"
	"github.com/reecerussell/open-social/cmd/backend/middleware"
)

const (
	usersAPIVar       = "USERS_API_URL"
	authAPIVar        = "AUTH_API_URL"
	postsAPIVar       = "POSTS_API_URL"
	mediaAPIVar       = "MEDIA_API_URL"
	tokenPublicKeyVar = "TOKEN_PUBLIC_KEY"
	bucketName        = "MEDIA_BUCKET"
)

func main() {
	ctn := buildServices()

	userHandler := ctn.GetService("UserHandler").(*handler.UserHandler)
	postHandler := ctn.GetService("PostHandler").(*handler.PostHandler)
	authHandler := ctn.GetService("AuthHandler").(*handler.AuthHandler)
	authMiddleware := ctn.GetService("AuthMiddleware").(*middleware.Authentication)

	app := core.NewApp()
	app.AddMiddleware(core.NewLoggingMiddleware())
	app.AddMiddleware(middleware.NewCors())
	app.AddMiddleware(authMiddleware)

	// User endpoints
	app.PostFunc("/users/follow/{userReferenceID}", userHandler.Follow)
	app.PostFunc("/users/unfollow/{userReferenceID}", userHandler.Unfollow)

	// Post endpoints
	app.PostFunc("/posts/like/{id}", postHandler.Like)
	app.PostFunc("/posts/unlike/{id}", postHandler.Unlike)
	app.PostFunc("/posts", postHandler.Create)
	app.GetFunc("/posts/{id}", postHandler.GetPost)

	// Auth endpoints
	app.PostFunc("/auth/register", authHandler.Register)
	app.PostFunc("/auth/token", authHandler.Token)

	// Frontend endpoints
	app.GetFunc("/feed", postHandler.GetFeed)
	app.GetFunc("/profile/{username}", userHandler.GetProfile)
	app.GetFunc("/me", userHandler.GetInfo)

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

	ctn.AddService("MediaClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(mediaAPIVar)
		client := media.New(url)
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
		postsClient := ctn.GetService("PostClient").(posts.Client)
		h := handler.NewUserHandler(client, authClient, postsClient)
		return h
	})

	ctn.AddService("PostHandler", func(ctn *core.Container) interface{} {
		client := ctn.GetService("PostClient").(posts.Client)
		mediaClient := ctn.GetService("MediaClient").(media.Client)
		h := handler.NewPostHandler(client, mediaClient)
		return h
	})

	ctn.AddService("AuthHandler", func(ctn *core.Container) interface{} {
		usersClient := ctn.GetService("UserClient").(users.Client)
		authClient := ctn.GetService("AuthClient").(auth.Client)
		h := handler.NewAuthHandler(usersClient, authClient)
		return h
	})

	return ctn
}
