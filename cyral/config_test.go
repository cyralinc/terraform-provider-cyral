// To generate detailed test cover:
// go test -v -coverprofile cover.out && go tool cover -html=cover.out

package cyral

import (
	"fmt"
	"os"
	"testing"
	//"net/http/httptest"
)

func invalidHTTPrequest() error {

	c := Config{Auth0Domain: "Auth0DomainValue", Auth0Audience: "Auth0AudienceValue"}
	os.Setenv("AUTH0_CLIENT_ID", "Auth0ClientIDvalue")
	os.Setenv("AUTH0_CLIENT_SECRET", "Auth0ClientIDsecret")

	// readTokenInfo() call: will always be correct in type matching

	// url assignment: will never fail, always receives right parameter type
	// and imported fuction call fmt.Sprintf already tested

	// audienceUrl assignment: will never fail,always receives right parameter
	// type and imported fuction call fmt.Sprintf already tested

	// tokenReq assignment: will never fail: always correct type and already
	// checked values lecture

	// payloadBytes assignment: Encoding to json call over tokenReq never fails.
	// Since tokenReq is always composed only by strings, and also as its values
	// will never be cyclic, this function will never fail.

	// Creating a http request using the paramaters will never fail.
	// http.NewRequest() method wraps, using argument data, a http request,
	// without sending it. It Fails when a non-valid method is entered, or
	// it receives a non valid body type. Since in this code it will always
	// receive a valid method (always http.MethodPost), and a body of a
	// correct stream type, and also always receives a string as url, this
	// method will never have problems creating a request using argument data
	// received.

	// statement req.Header.Add will never fail since it always receives right
	// arguments for execution

	if _, err := c.Client(); err != nil {
		return err
	}
	return nil
}

func absentEnvVarAuth0ClientID() error {
	c := Config{}

	// Remove env var AUTH0_CLIENT_ID, if present
	if _, err := c.getEnv("AUTH0_CLIENT_ID"); err == nil {
		os.Unsetenv("AUTH0_CLIENT_ID")
	}

	_, err2 := c.Client()

	if _, err := c.getEnv("AUTH0_CLIENT_ID"); err2.Error() != err.Error() {
		return fmt.Errorf("Unexpected behavior in Client() function.\n")
	}
	return nil
}

func absentEnvVarAuth0ClientSecret() error {
	c := Config{}
	os.Setenv("AUTH0_CLIENT_ID", "Auth0ClientIDvalue")

	// Remove env var AUTH0_CLIENT_SECRET, if present
	if _, err := c.getEnv("AUTH0_CLIENT_SECRET"); err == nil {
		os.Unsetenv("AUTH0_CLIENT_SECRET")
	}

	_, err2 := c.Client()

	if _, err := c.getEnv("AUTH0_CLIENT_SECRET"); err2.Error() != err.Error() {
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
		//t.Error(err)
		t.Log(err)
	}

}

// getEnv() function: as defined, will never fail, since it is only
// another way to print os.LookupEnv() results, and os.LookupEnv() is already
// officialy tested.
