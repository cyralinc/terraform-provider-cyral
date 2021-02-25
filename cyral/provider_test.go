package cyral

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cyral": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

/*
func absentEnvVarAuth0ClientID() error {
	c := NewClient("",
		"Auth0ClientSecretValue",
		"Auth0DomainValue",
		"Auth0AudienceValue",
		"ControlPlaneValue")

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
