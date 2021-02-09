package main

import (
	"log"
	"os"
	"os/signal"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/cmd/posts/handler"
	"github.com/reecerussell/open-social/cmd/posts/provider"
	"github.com/reecerussell/open-social/cmd/posts/repository"
	"github.com/reecerussell/open-social/database"
)

const (
	connectionStringVar = "CONNECTION_STRING"
	usersAPIVar         = "USERS_API_URL"
	kafkaHostVar        = "KAFKA_HOST"
)

func main() {
	ctn := buildServices()
	db := ctn.GetService("Database").(database.Database)

	createPost := ctn.GetService("CreatePostHandler").(*handler.CreatePostHandler)
	feedhandler := ctn.GetService("FeedHandler").(*handler.FeedHandler)
	profileFeedHandler := ctn.GetService("ProfileFeedHandler").(*handler.ProfileFeedHandler)
	likePost := ctn.GetService("LikePostHandler").(*handler.LikePostHandler)
	unlikePost := ctn.GetService("UnlikePostHandler").(*handler.UnlikePostHandler)
	getPost := ctn.GetService("GetPostHandler").(*handler.GetPostHandler)

	app := core.NewApp()
	app.AddHealthCheck(database.NewHealthCheck(db))
	app.AddMiddleware(core.NewLoggingMiddleware())

	app.Post("/posts", createPost)
	app.Get("/posts/{postReferenceID}/{userReferenceID}", getPost)
	app.Post("/posts/like", likePost)
	app.Post("/posts/unlike", unlikePost)
	app.Get("/feed/{userReferenceId}", feedhandler)
	app.Get("/profile/feed/{username}/{userReferenceID}", profileFeedHandler)

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

	ctn.AddService("PostProvider", func(ctn *core.Container) interface{} {
		db := ctn.GetService("Database").(database.Database)
		return provider.NewPostProvider(db)
	})

	ctn.AddService("PostRepository", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		return repository.NewPostRepository(url)
	})

	ctn.AddService("LikeRepository", func(ctn *core.Container) interface{} {
		db := ctn.GetService("Database").(database.Database)
		return repository.NewLikeRepository(db)
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

	ctn.AddService("FeedHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("PostRepository").(repository.PostRepository)

		return handler.NewFeedHandler(repo)
	})

	ctn.AddService("ProfileFeedHandler", func(ctn *core.Container) interface{} {
		provider := ctn.GetService("PostProvider").(provider.PostProvider)

		return handler.NewProfileFeedHandler(provider)
	})

	ctn.AddService("LikePostHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("PostRepository").(repository.PostRepository)
		likes := ctn.GetService("LikeRepository").(repository.LikeRepository)

		return handler.NewLikePostHandler(repo, likes)
	})

	ctn.AddService("UnlikePostHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("PostRepository").(repository.PostRepository)
		likes := ctn.GetService("LikeRepository").(repository.LikeRepository)

		return handler.NewUnlikePostHandler(repo, likes)
	})

	ctn.AddService("GetPostHandler", func(ctn *core.Container) interface{} {
		provider := ctn.GetService("PostProvider").(provider.PostProvider)

		return handler.NewGetPostHandler(provider)
	})

	return ctn
}
