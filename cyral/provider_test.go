package cyral

/*
import (
	"fmt"
	"os"
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
*/
