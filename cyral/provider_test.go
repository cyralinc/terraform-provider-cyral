package cyral

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	EnvVarTFAcc = "TF_ACC"

	// Ex: stable.dev.cyral.com:8000
	EnvVarCPURL = "CYRAL_TF_CP_URL"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var providerFactories = map[string]func() (*schema.Provider, error){
	"cyral": func() (*schema.Provider, error) {
		p := Provider()
		p.Schema["control_plane"] = &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			DefaultFunc: schema.EnvDefaultFunc(EnvVarCPURL, nil),
		}
		return p, nil
	},
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(EnvVarTFAcc); v == "" {
		t.Fatalf("%q must be set for acceptance tests", EnvVarTFAcc)
	}
}
