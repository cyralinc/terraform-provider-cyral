package provider

import (
	"testing"
)

const (
	EnvVarTFAcc = "TF_ACC"
)

func TestAccProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
