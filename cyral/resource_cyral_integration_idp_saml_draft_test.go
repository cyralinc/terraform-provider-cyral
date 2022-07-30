package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func initialGenericSAMLDraftConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-1",
		DisableIdPInitiatedLogin: false,
		IdpType:                  "some-idp-type-1",
		Attributes: NewRequiredUserAttributes(
			"first-name-1",
			"last-name-1",
			"email-1",
			"groups-1",
		),
	}
}

func updatedGenericSAMLDraftConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-2",
		DisableIdPInitiatedLogin: true,
		IdpType:                  "some-idp-type-2",
		Attributes: NewRequiredUserAttributes(
			"first-name-2",
			"last-name-2",
			"email-2",
			"groups-2",
		),
	}
}

func TestAccIntegrationIdPSAMLDraftResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPSAMLDraftTest(t,
		initialGenericSAMLDraftConfig(), "main_test")
	updatedConfig, updatedChecks := setupIntegrationIdPSAMLDraftTest(t,
		updatedGenericSAMLDraftConfig(), "main_test")

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

func setupIntegrationIdPSAMLDraftTest(t *testing.T, draft *GenericSAMLDraft, resName string) (
	string,
	resource.TestCheckFunc,
) {
	config := formatGenericSAMLDraftIntoConfig(draft, resName)

	nonEmptyRegex, err := regexp.Compile(".+")
	require.NoError(t, err)

	resourceFullName := fmt.Sprintf("cyral_integration_idp_saml_draft.%s", resName)
	checkFunc := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "display_name",
			draft.DisplayName),
		resource.TestCheckResourceAttr(resourceFullName, "disable_idp_initiated_login",
			strconv.FormatBool(draft.DisableIdPInitiatedLogin)),
		resource.TestCheckResourceAttr(resourceFullName, "idp_type",
			draft.IdpType),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.first_name",
			draft.Attributes.FirstName.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.last_name",
			draft.Attributes.LastName.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.email",
			draft.Attributes.Email.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.groups",
			draft.Attributes.Groups.Name),
		resource.TestMatchResourceAttr(resourceFullName, "id",
			nonEmptyRegex),
		resource.TestMatchResourceAttr(resourceFullName, "sp_metadata",
			nonEmptyRegex),
	)

	return config, checkFunc
}

func formatGenericSAMLDraftIntoConfig(draft *GenericSAMLDraft, resName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		disable_idp_initiated_login = %t
		idp_type = "%s"
		attributes {
			first_name = "%s"
			last_name = "%s"
			email = "%s"
			groups = "%s"
		}
	}`, resName, draft.DisplayName, draft.DisableIdPInitiatedLogin, draft.IdpType,
		draft.Attributes.FirstName.Name,
		draft.Attributes.LastName.Name,
		draft.Attributes.Email.Name,
		draft.Attributes.Groups.Name,
	)
}
