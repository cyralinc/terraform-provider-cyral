package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func invalidAuth0DomainFormat() error {
	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		"^^^exampleInvalidDomain",
		"ExampleAuth0Audience",
		"SomeControlPlane")

	if err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when Auth0 domain has invalid format; err: %v",
			err.Error())
	}

	return nil
}

func invalidAuth0DomainValue() error {
	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		"invalidDomain",
		"ExampleAuth0Audience",
		"SomeControlPlane")

	if err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when Auth0 domain has invalid value; err: %v",
			err.Error())
	}

	return nil
}

func serverDown() error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		ts.URL[8:len(ts.URL)],
		"exampleAud",
		"SomeControlPlane")

	ts.URL = ts.URL + "/oauth/token"

	ts.Close()

	if err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when server is down; err: %v",
			err.Error())
	}

	return nil
}

func timeoutResponse() error {
	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for {
		}
	}))
	defer ts.Close()

	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		ts.URL[8:len(ts.URL)],
		"exampleAud",
		"SomeControlPlane")

	ts.URL = ts.URL + "/oauth/token"

	if err == nil {
		return fmt.Errorf("error in timeoutResponse(); err: %v", err.Error())
	}

	return nil
}

func reqOK() error {
	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json")

		tokenRes := Auth0TokenResponse{
			AccessToken: "ExampleAcessToken",
			TokenType:   "ExampleTokenType",
		}
		jsonTokenRes, _ := json.Marshal(tokenRes)
		w.Write(jsonTokenRes)
	}))
	defer ts.Close()

	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		ts.URL[8:len(ts.URL)],
		"exampleAud",
		"SomeControlPlane")
	ts.URL = ts.URL + "/oauth/token"

	if err != nil {
		return fmt.Errorf("error in reqOK(); err: %v", err.Error())
	}

	return nil
}

func reqFail() error {
	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Any response different than 200 (http.StatusOK) is an error.
		w.WriteHeader(http.StatusBadRequest)
	}))

	_, err := NewClient(
		"ExampleAuth0ClientIDvalue",
		"ExampleAuth0ClientSecretValue",
		ts.URL[8:len(ts.URL)],
		"exampleAud",
		"SomeControlPlane")
	ts.URL = ts.URL + "/oauth/token"

	if err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("status code %d", http.StatusBadRequest)) {
			return fmt.Errorf("error in reqFail(); err: %v", err.Error())
		}
	}
	defer ts.Close()

	return nil
}

func TestClient(t *testing.T) {
	if err := invalidAuth0DomainFormat(); err != nil {
		t.Error(err)
	}

	if err := invalidAuth0DomainValue(); err != nil {
		t.Error(err)
	}

	if err := serverDown(); err != nil {
		t.Error(err)
	}

	if err := reqOK(); err != nil {
		t.Error(err)
	}

	if err := reqFail(); err != nil {
		t.Error(err)
	}
}
