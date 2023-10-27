package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	EnvVarTFAcc = "TF_ACC"
)

func TestAccProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

var providerFactories = map[string]func() (*schema.Provider, error){
	"cyral": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}
