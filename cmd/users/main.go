package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	hashpkg "github.com/reecerussell/adaptive-password-hasher"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/cmd/users/handler"
	"github.com/reecerussell/open-social/cmd/users/password"
	"github.com/reecerussell/open-social/cmd/users/provider"
	"github.com/reecerussell/open-social/cmd/users/repository"
	"github.com/reecerussell/open-social/database"
	"github.com/reecerussell/open-social/util"
)

const (
	connectionStringVar = "CONNECTION_STRING"
	configFileVar       = "CONFIG_FILE"
)

func main() {
	cnf := buildConfig()
	ctn := buildServices(cnf)
	db := ctn.GetService("Database").(database.Database)

	createUser := ctn.GetService("CreateUserHandler").(*handler.CreateUserHandler)
	getClaims := ctn.GetService("GetClaimsHandler").(*handler.GetClaimsHandler)
	getIDByReference := ctn.GetService("GetIDByReferenceHandler").(*handler.GetIDByReferenceHandler)
	getProfile := ctn.GetService("GetProfileHandler").(*handler.GetProfileHandler)
	getInfo := ctn.GetService("GetInfoHandler").(*handler.GetInfoHandler)

	app := core.NewApp()
	app.AddHealthCheck(database.NewHealthCheck(db))
	app.AddMiddleware(core.NewLoggingMiddleware())

	app.Post("/users", createUser)
	app.Get("/users/id/{referenceId}", getIDByReference)
	app.Post("/claims", getClaims)
	app.Get("/profile/{username}/{userReferenceID}", getProfile)
	app.Get("/info/{userReferenceID}", getInfo)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

// Config is a configuration model for the service.
type Config struct {
	PasswordValidatorOptions *password.Options     `json:"passwordValidator"`
	PasswordHasherOptions    *password.HashOptions `json:"passwordHasher"`
}

func buildConfig() *Config {
	filename := util.ReadEnv(configFileVar, "config.json")

	var cnf Config
	err := core.ReadConfig(filename, &cnf)
	if err != nil {
		panic(fmt.Errorf("failed to read config: %v", err))
	}

	return &cnf
}

func buildServices(cnf *Config) *core.Container {
	ctn := core.NewContainer()

	ctn.AddSingleton("Config", func(ctn *core.Container) interface{} {
		return cnf
	})

	ctn.AddSingleton("Database", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		db, err := database.New(url)
		if err != nil {
			panic(err)
		}

		return db
	})

	ctn.AddService("PasswordValidator", func(ctn *core.Container) interface{} {
		cnf := ctn.GetService("Config").(*Config)
		val := password.New(cnf.PasswordValidatorOptions)
		return val
	})

	ctn.AddService("PasswordHasher", func(ctn *core.Container) interface{} {
		cnf := ctn.GetService("Config").(*Config)
		hasher, err := hashpkg.New(
			cnf.PasswordHasherOptions.IterationCount,
			cnf.PasswordHasherOptions.SaltSize,
			cnf.PasswordHasherOptions.KeySize,
			cnf.PasswordHasherOptions.HashKey)
		if err != nil {
			panic(fmt.Errorf("failed to build PasswordHasher: %v", err))
		}

		return hasher
	})

	ctn.AddService("UserRepository", func(ctn *core.Container) interface{} {
		url := os.Getenv(connectionStringVar)
		return repository.NewUserRepository(url)
	})

	ctn.AddService("UserProvider", func(ctn *core.Container) interface{} {
		db := ctn.GetService("Database").(database.Database)
		return provider.NewUserProvider(db)
	})

	ctn.AddService("CreateUserHandler", func(ctn *core.Container) interface{} {
		val := ctn.GetService("PasswordValidator").(password.Validator)
		hasher := ctn.GetService("PasswordHasher").(hashpkg.Hasher)
		repo := ctn.GetService("UserRepository").(repository.UserRepository)

		return handler.NewCreateUserHandler(val, hasher, repo)
	})

	ctn.AddService("GetClaimsHandler", func(ctn *core.Container) interface{} {
		hasher := ctn.GetService("PasswordHasher").(hashpkg.Hasher)
		repo := ctn.GetService("UserRepository").(repository.UserRepository)

		return handler.NewGetClaimsHandler(hasher, repo)
	})

	ctn.AddService("GetIDByReferenceHandler", func(ctn *core.Container) interface{} {
		repo := ctn.GetService("UserRepository").(repository.UserRepository)

		return handler.NewGetIDByReferenceHandler(repo)
	})

	ctn.AddService("GetProfileHandler", func(ctn *core.Container) interface{} {
		provider := ctn.GetService("UserProvider").(provider.UserProvider)

		return handler.NewGetProfileHandler(provider)
	})

	ctn.AddService("GetInfoHandler", func(ctn *core.Container) interface{} {
		provider := ctn.GetService("UserProvider").(provider.UserProvider)

		return handler.NewGetInfoHandler(provider)
	})

	return ctn
}
