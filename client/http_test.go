package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello World"}`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	var data map[string]string
	err := hc.Get("/test", &data)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", data["message"])
}

func TestHTTPGet_ReturnsErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"message":"an error occured"}`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	var data map[string]string
	err := hc.Get("/test", &data)
	assert.Equal(t, "an error occured", err.Error())
}

func TestHTTPGet_GivenInvalidURL_ReturnsError(t *testing.T) {
	hc := NewHTTP("")
	var data map[string]string
	err := hc.Get("324@2-asd", &data)
	assert.NotNil(t, err)
}

func TestHTTPGet_GivenNilDest_ReturnsError(t *testing.T) {
	hc := NewHTTP("")
	err := hc.Get("", nil)
	assert.Equal(t, "http: dest can not contain a nil value for a GET request", err.Error())
}

func TestHTTPGet_ReturnsNonJSONError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	var data map[string]string
	err := hc.Get("/test", &data)
	assert.Equal(t, "http: server returned a 500 status code", err.Error())
}

func TestHTTPGet_ReturnsErrorWithInvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"Hello World"`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	var data map[string]string
	err := hc.Get("/test", &data)
	assert.Equal(t, "http: failed to read json response, status code: 500", err.Error())
}

func TestHTTPGet_ReturnsOKWithInvalidJSON_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/test", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello World"`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	var data map[string]string
	err := hc.Get("/test", &data)
	assert.Equal(t, "http: failed to read successful response", err.Error())
}

func TestHTTPPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/test", r.URL.Path)

		var data map[string]string
		_ = json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		assert.Equal(t, "Hello World", data["message"])

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"Hello World"}`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	data := map[string]string{
		"message": "Hello World",
	}
	var res map[string]string

	err := hc.Post("/test", data, &res)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", res["message"])
}

func TestHTTPPost_ReturnsErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/test", r.URL.Path)

		var data map[string]string
		_ = json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		assert.Equal(t, "Hello World", data["message"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"message":"an error occured"}`))
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	data := map[string]string{
		"message": "Hello World",
	}
	err := hc.Post("/test", data, nil)
	assert.Equal(t, "an error occured", err.Error())
}

func TestHTTPPost_GivenInvalidURL_ReturnsError(t *testing.T) {
	hc := NewHTTP("")
	err := hc.Post("324@2-asd", nil, nil)
	assert.NotNil(t, err)
}

func TestHTTPPost_ReturnsNonJSONError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/test", r.URL.Path)

		var data map[string]string
		_ = json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	hc := NewHTTP(server.URL)

	data := map[string]string{
		"message": "Hello World",
	}
	err := hc.Post("/test", data, nil)
	assert.Equal(t, "http: server returned a 500 status code", err.Error())
}
