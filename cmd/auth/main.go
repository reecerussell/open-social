package main

import (
	"crypto"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/reecerussell/gojwt"
	"github.com/reecerussell/gojwt/rsa"

	core "github.com/reecerussell/open-social"
	"github.com/reecerussell/open-social/client/users"
	"github.com/reecerussell/open-social/cmd/auth/handler"
	"github.com/reecerussell/open-social/util"
)

const (
	usersAPIVar        = "USERS_API_URL"
	configFileVar      = "CONFIG_FILE"
	tokenPrivateKeyVar = "TOKEN_PRIVATE_KEY"
)

func main() {
	cnf := buildConfig()
	ctn := buildServices(cnf)

	userHandler := ctn.GetService("TokenHandler").(*handler.TokenHandler)

	app := core.NewApp("0.0.0.0:80")
	app.AddMiddleware(core.NewLoggingMiddleware())

	app.Post("/token", userHandler)

	go app.Serve()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	<-stop

	log.Println("App stopped.")
}

// Config is a configuration model for the service.
type Config struct {
	Token *TokenConfig `json:"token"`
}

// TokenConfig contains config for generating access tokens.
type TokenConfig struct {
	ExpiryMinutes int `json:"expiryMinutes"`
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

	ctn.AddService("UserClient", func(ctn *core.Container) interface{} {
		url := os.Getenv(usersAPIVar)
		client := users.New(url)
		return client
	})

	ctn.AddService("TokenAlg", func(ctn *core.Container) interface{} {
		alg, err := rsa.NewFromFile(os.Getenv(tokenPrivateKeyVar), crypto.SHA256)
		if err != nil {
			panic(err)
		}

		return alg
	})

	ctn.AddService("TokenHandler", func(ctn *core.Container) interface{} {
		client := ctn.GetService("UserClient").(users.Client)
		alg := ctn.GetService("TokenAlg").(gojwt.Algorithm)
		cnf := ctn.GetService("Config").(*Config)
		h := handler.NewTokenHandler(client, alg, cnf.Token.ExpiryMinutes)
		return h
	})

	return ctn
}
