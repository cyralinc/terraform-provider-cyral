package cyral

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func absentEnvVarAuth0ClientID() error {
	c := Config{}

	// Remove env var AUTH0_CLIENT_ID, if present
	os.Unsetenv("AUTH0_CLIENT_ID")

	_, err := c.Client()

	if _, err2 := c.getEnv("AUTH0_CLIENT_ID"); err2.Error() != err.Error() {
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
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
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
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
		msg := "Unexpected Client() behavior when receives invalid Auth0 domain format"
		return fmt.Errorf(msg)
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
		msg := "Unexpected Client() behavior when receives invalid Auth0 domain value"
		return fmt.Errorf(msg)
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
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
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
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
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
		w.WriteHeader(201)
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
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
	}

	return nil
}

func requisitionContentProblem() error {

	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// General requisition error test: 4xx family
		w.WriteHeader(401)
	}))
	defer ts.Close()

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	if _, err := c.Client(); err != nil {
		if err.Error()[0:36] == "unable to get access token from json" {
			return fmt.Errorf("Absent error warning for requisition values errors condition.")
		}

		return nil
	}

	return fmt.Errorf("Unexpected behavior in Client() function.\n")
}

func serverInternalError() error {

	os.Setenv("AUTH0_CLIENT_ID", "ExampleAuth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "ExampleAuth0ClientSecretValue")

	// Disables client's certificate authority validation, in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Any server error: 5xx family
		w.WriteHeader(500)
	}))
	defer ts.Close()

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "exampleAud"}
	ts.URL = ts.URL + "/oauth/token"

	if _, err := c.Client(); err != nil {

		if err.Error()[0:36] == fmt.Errorf("unable to get access token from json").Error() {
			return fmt.Errorf("Absent error warning for internal server errors condition.\n")
		}
		return nil
	}

	return fmt.Errorf("Unexpected behavior in Client() function")
}

// Validation of server certficate also works but its not used in here
// due to always logging behavior, making unclean results for go test commmand

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

	// Needs thread / goroutine to stop after certain time if it doesnt
	// timeout (as is behaving)
	/*if err := timeoutResponse(); err != nil {
		t.Error(err)
	}*/

	if err := reqOK(); err != nil {
		t.Error(err)
	}

	if err := serverInternalError(); err != nil {
		t.Error(err)
	}

	if err := requisitionContentProblem(); err != nil {
		t.Error(err)
	}

}
