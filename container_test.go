package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainer_GetNonExistantService_Panics(t *testing.T) {
	ctn := &Container{
		srvs:    make(map[string]interface{}),
		srvConf: make(map[string]*ServiceConfig),
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	_ = ctn.GetService("myService")
}

func TestContainer_AddService(t *testing.T) {
	builder := func(ctn *Container) interface{} {
		return nil
	}
	name := "myService"

	ctn := NewContainer()
	ctn.AddService(name, builder)

	srv, ok := ctn.srvConf[name]
	assert.True(t, ok)
	assert.False(t, srv.Singleton)
}

func TestContainer_AddSingleton(t *testing.T) {
	builder := func(ctn *Container) interface{} {
		return nil
	}
	name := "myService"

	ctn := NewContainer()
	ctn.AddSingleton(name, builder)

	srv, ok := ctn.srvConf[name]
	assert.True(t, ok)
	assert.True(t, srv.Singleton)
}

func TestContainer_SingletonNotRecreated(t *testing.T) {
	ctn := &Container{
		srvs:    make(map[string]interface{}),
		srvConf: make(map[string]*ServiceConfig),
	}

	ctn.srvConf["test"] = &ServiceConfig{
		Singleton: true,
		Build: func(ctn *Container) interface{} {
			srvValue := "My super cool service"
			return &srvValue
		},
	}

	srv := ctn.GetService("test")
	srv2 := ctn.GetService("test")

	if srv != srv2 {
		t.Error("Expected the services to be equal")
	}
}

func TestContainer_TransientIsRecreated(t *testing.T) {
	ctn := &Container{
		srvs:    make(map[string]interface{}),
		srvConf: make(map[string]*ServiceConfig),
	}

	ctn.srvConf["test"] = &ServiceConfig{
		Singleton: false,
		Build: func(ctn *Container) interface{} {
			srvValue := "My super cool service"
			return &srvValue
		},
	}

	srv := ctn.GetService("test")
	srv2 := ctn.GetService("test")

	if srv == srv2 {
		t.Error("Expected the services to not be equal")
	}
}
