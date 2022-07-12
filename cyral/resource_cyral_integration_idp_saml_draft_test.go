package cyral

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func initialSAMLDraftConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-1",
		DisableIdPInitiatedLogin: true,
	}
}

func updatedSAMLDraftConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-2",
		DisableIdPInitiatedLogin: false,
	}
}

func TestAccIntegrationIdPSAMLDraftResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPSAMLDraftTest(
		initialSAMLDraftConfig(), "initial_test")
	updatedConfig, updatedChecks := setupIntegrationIdPSAMLDraftTest(
		updatedSAMLDraftConfig(), "updated_test")

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check:  initialChecks,
			},
			{
				Config: updatedConfig,
				Check:  updatedChecks,
			},
		},
	})
}

func setupIntegrationIdPSAMLDraftTest(draft *GenericSAMLDraft, resName string) (
	string,
	resource.TestCheckFunc,
) {
	config := formatGenericSAMLDraftIntoConfig(draft, resName)

	resourceFullName := fmt.Sprintf("cyral_integration_idp_saml_draft.%s", resName)
	checkFunc := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "display_name",
			draft.DisplayName),
		resource.TestCheckResourceAttr(resourceFullName, "disable_idp_initiated_login",
			strconv.FormatBool(draft.DisableIdPInitiatedLogin)),
	)

	return config, checkFunc
}

func formatGenericSAMLDraftIntoConfig(draft *GenericSAMLDraft, resName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		disable_idp_initiated_login = %t
	}`, resName, draft.DisplayName, draft.DisableIdPInitiatedLogin)
}
