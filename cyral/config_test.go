// To generate detailed test cover:
// go test -v -coverprofile cover.out && go tool cover -html=cover.out

package cyral

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func invalidHTTPrequest() error {
	// Disables certificate authority validation in order to
	// successfully mock https requests
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	os.Setenv("AUTH0_CLIENT_ID", "Auth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "Auth0ClientIDsecret")

	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))

	c := Config{Auth0Domain: ts.URL[8:len(ts.URL)], Auth0Audience: "cyral"}
	ts.URL = ts.URL + "/oauth/token"
	defer ts.Close()

	_, err := c.Client()

	return err
}

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

func TestClient(t *testing.T) {

	if err := absentEnvVarAuth0ClientID(); err != nil {
		t.Error(err)
	}
	if err := absentEnvVarAuth0ClientSecret(); err != nil {
		t.Error(err)
	}

	if err := invalidHTTPrequest(); err != nil {
		t.Error(err)
	}

}
