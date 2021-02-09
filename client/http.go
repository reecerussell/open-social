package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HTTP is a high level interface used to wrap the core http package,
// resulting in making HTTP requests simplier among services, using a
// standardised response pattern.
type HTTP interface {
	// Get makes a GET request to the given url, then reads a
	// JSON response to the given destination. dest can not be nil.
	Get(url string, dest interface{}) error

	// Post makes a POST request to the given url, with the given body,
	// which will be in JSON format. If dest is not nil, the response
	// body will be JSON decoded to the given destination.
	Post(url string, body, dest interface{}) error
}

// NewHTTP returns a new instance of HTTP with a base url.
func NewHTTP(baseURL string) HTTP {
	return &httpClient{
		base: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
	}
}

type httpClient struct {
	base    *http.Client
	baseURL string
}

func (hc *httpClient) Get(url string, dest interface{}) error {
	if dest == nil {
		return errors.New("http: dest can not contain a nil value for a GET request")
	}

	return hc.makeRequest(http.MethodGet, url, nil, dest)
}

func (hc *httpClient) Post(url string, body, dest interface{}) error {
	return hc.makeRequest(http.MethodPost, url, body, dest)
}

func (hc *httpClient) makeRequest(method, url string, body, respDest interface{}) error {
	reqBody := getRequestBody(method, body)
	req, _ := http.NewRequest(method, getRequestURL(hc.baseURL, url), reqBody)

	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := hc.base.Do(req)
	if err != nil {
		return fmt.Errorf("http: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.Header.Get("Content-Type") != "application/json" {
			return fmt.Errorf("http: server returned a %d status code", resp.StatusCode)
		}

		var data ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("http: failed to read json response, status code: %d", resp.StatusCode)
		}

		return NewError(resp.StatusCode, data.Message)
	}

	if respDest != nil {
		err := json.NewDecoder(resp.Body).Decode(respDest)
		if err != nil {
			return fmt.Errorf("http: failed to read successful response")
		}
	}

	return nil
}

// getRequestBody returns an io.Reader containing the JSON encoded value of body.
// If method is "GET" or if the body is nil, a nil-value will be returned.
func getRequestBody(method string, body interface{}) io.Reader {
	if method == http.MethodGet || body == nil {
		return nil
	}

	jsonBytes, _ := json.Marshal(body)
	return bytes.NewReader(jsonBytes)
}

// getRequestURL concatenates the base and url strings into a url.
func getRequestURL(base, url string) string {
	if base != "" && !strings.HasSuffix(base, "/") {
		base = base + "/"
	}

	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}

	return base + url
}
