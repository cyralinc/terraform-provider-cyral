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

func TestInvalidAuth0DomainFormat(t *testing.T) {
	test := func(keycloakProvider bool) {
		_, err := NewClient(
			"ExampleAuth0ClientIDvalue",
			"ExampleAuth0ClientSecretValue",
			"^^^exampleInvalidDomain",
			"ExampleAuth0Audience",
			"SomeControlPlane",
			keycloakProvider,
			true)

		if err == nil {
			t.Error(fmt.Errorf(
				"unexpected behavior in Client() when Auth0 domain has invalid format; err: %v",
				err.Error()))
		}
	}

	test(true)
	test(false)
}

func TestInvalidAuth0DomainValue(t *testing.T) {
	test := func(keycloakProvider bool) {
		_, err := NewClient(
			"ExampleAuth0ClientIDvalue",
			"ExampleAuth0ClientSecretValue",
			"invalidDomain",
			"ExampleAuth0Audience",
			"SomeControlPlane",
			keycloakProvider,
			true)

		if err == nil {
			t.Error(fmt.Errorf(
				"unexpected behavior in Client() when Auth0 domain has invalid value; "+
					"keycloakProvider: %t; err: %v",
				keycloakProvider, err.Error()))
		}
	}

	test(true)
	test(false)
}

func TestServerDown(t *testing.T) {
	test := func(keycloakProvider bool) {
		// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))

		_, err := NewClient(
			"ExampleAuth0ClientIDvalue",
			"ExampleAuth0ClientSecretValue",
			ts.URL[8:len(ts.URL)],
			"exampleAud",
			"SomeControlPlane",
			keycloakProvider,
			true)

		ts.URL = ts.URL + "/oauth/token"

		ts.Close()

		if err == nil {
			t.Error(fmt.Errorf(
				"unexpected behavior in Client() when server is down. "+
					"keycloakProvider: %t; err: %v",
				keycloakProvider, err.Error()))
		}
	}

	test(true)
	test(false)
}

func TestTimeoutResponse(t *testing.T) {
	test := func(keycloakProvider bool) {
		// Disables client's certificate authority validation, in order to
		// successfully mock https requests
		// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
		defer ts.Close()

		_, err := NewClient(
			"ExampleAuth0ClientIDvalue",
			"ExampleAuth0ClientSecretValue",
			ts.URL[8:len(ts.URL)],
			"exampleAud",
			"SomeControlPlane",
			keycloakProvider, true)

		ts.URL = ts.URL + "/oauth/token"

		if err == nil {
			t.Error(fmt.Errorf("error in timeoutResponse(); keycloakProvider: %t; err: %v",
				keycloakProvider, err.Error()))
		}
	}

	test(true)
	test(false)
}

func TestReqOK(t *testing.T) {
	test := func(keycloakProvider bool) {
		// Disables client's certificate authority validation, in order to
		// successfully mock https requests
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Add("content-type", "application/json")

			tokenRes := TokenResponse{
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
			ts.URL[8:len(ts.URL)],
			keycloakProvider,
			true)
		ts.URL = ts.URL + "/oauth/token"

		if err != nil {
			t.Error(fmt.Errorf("error in reqOK(); keycloakProvider: %t; err: %v",
				keycloakProvider, err.Error()))
		}
	}

	test(true)
	test(false)
}

func TestReqFail(t *testing.T) {
	test := func(keycloakProvider bool) {
		// Disables client's certificate authority validation, in order to
		// successfully mock https requests
		// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Any response different than 200 (http.StatusOK) is an error.
			w.WriteHeader(http.StatusBadRequest)
		}))

		_, err := NewClient(
			"ExampleAuth0ClientIDvalue",
			"ExampleAuth0ClientSecretValue",
			ts.URL[8:len(ts.URL)],
			"exampleAud",
			ts.URL[8:len(ts.URL)],
			keycloakProvider,
			true)
		ts.URL = ts.URL + "/oauth/token"

		if err != nil {
			if !strings.Contains(err.Error(), fmt.Sprintf("status code %d", http.StatusBadRequest)) {
				t.Error(fmt.Errorf("error in reqFail(); keycloakProvider: %t; err: %v",
					keycloakProvider, err.Error()))
			}
		}
		defer ts.Close()
	}

	test(true)
	test(false)
}
