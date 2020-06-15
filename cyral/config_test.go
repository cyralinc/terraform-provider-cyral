package cyral

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func absentEnvVarAuth0ClientID() error {
	c := Config{}

	// Remove env var AUTH0_CLIENT_ID, if present
	os.Unsetenv("AUTH0_CLIENT_ID")

	_, err := c.Client()

	if _, err2 := c.getEnv("AUTH0_CLIENT_ID"); err2.Error() != err.Error() {
		return fmt.Errorf(
			"unexpected behavior in Client() when AUTH0_CLIENT_ID is empty; err: %v",
			err.Error())
	}
	return nil
}

func absentEnvVarAuth0ClientSecret() error {
	c := Config{}
	os.Setenv("AUTH0_CLIENT_ID", "Auth0ClientIDvalue")

	// Remove env var AUTH0_CLIENT_SECRET, if present
	os.Unsetenv("AUTH0_CLIENT_SECRET")

	_, err := c.Client()

	if _, err2 := c.getEnv("AUTH0_CLIENT_SECRET"); err2.Error() != err.Error() {
		return fmt.Errorf(
			"unexpected behavior in Client() when AUTH0_CLIENT_SECRET is empty; err: %v",
			err.Error())
	}
	return nil
}

func invalidAuth0DomainFormat() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")
	c := Config{
		Auth0Domain: "^^^exampleInvalidDomain",
	}

	if _, err := c.Client(); err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when Auth0 domain has invalid format; err: %v",
			err.Error())
	}

	return nil
}

func invalidAuth0DomainValue() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")
	c := Config{
		Auth0Domain: "invalidDomain",
	}

	if _, err := c.Client(); err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when Auth0 domain has invalid value; err: %v",
			err.Error())
	}

	return nil
}

func serverDown() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	ts.Close()

	if _, err := c.Client(); err == nil {
		return fmt.Errorf(
			"unexpected behavior in Client() when server is down; err: %v",
			err.Error())
	}

	return nil
}

func timeoutResponse() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for {
		}
	}))
	defer ts.Close()

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	if _, err := c.Client(); err == nil {
		return fmt.Errorf("error in timeoutResponse(); err: %v", err.Error())
	}

	return nil
}

func reqOK() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("content-type", "application/json")

		tokenRes := auth0TokenResponse{
			AccessToken: "ExampleAcessToken",
			TokenType:   "ExampleTokenType",
		}
		jsonTokenRes, _ := json.Marshal(tokenRes)
		w.Write(jsonTokenRes)
	}))
	defer ts.Close()

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	if _, err := c.Client(); err != nil {
		return fmt.Errorf("error in reqOK(); err: %v", err.Error())
	}

	return nil
}

func reqFail() error {
	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Any response different than 200 (http.StatusOK) is an error.
		w.WriteHeader(http.StatusBadRequest)
	}))

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	if _, err := c.Client(); err != nil {
		if !strings.Contains(err.Error(), fmt.Sprintf("status code %d", http.StatusBadRequest)) {
			return fmt.Errorf("error in reqFail(); err: %v", err.Error())
		}
	}
	defer ts.Close()

	return nil
}

func TestClient(t *testing.T) {
	if err := absentEnvVarAuth0ClientID(); err != nil {
		t.Error(err)
	}

	if err := absentEnvVarAuth0ClientSecret(); err != nil {
		t.Error(err)
	}

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
