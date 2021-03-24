package cyral

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Example value: "stable.dev.cyral.com:8000"
const EnvVarControlPlaneBaseURL = "TEST_TPC_CP_URL"

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cyral": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(EnvVarControlPlaneBaseURL); v == "" {
		t.Fatalf("%q must be set for acceptance tests", EnvVarControlPlaneBaseURL)
	}
}
