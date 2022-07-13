package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func initialGenericSAMLConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-1",
		DisableIdPInitiatedLogin: false,
		IdpType:                  "some-idp-type-1",
		RequiredUserAttributes: NewRequiredUserAttributes(
			"first-name-1",
			"last-name-1",
			"email-1",
			"groups-1",
		),
	}
}

func updatedGenericSAMLConfig() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName:              "tf-test-saml-draft-2",
		DisableIdPInitiatedLogin: true,
		IdpType:                  "some-idp-type-2",
		RequiredUserAttributes: NewRequiredUserAttributes(
			"first-name-2",
			"last-name-2",
			"email-2",
			"groups-2",
		),
	}
}

func TestAccIntegrationIdPGenericSAMLResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPGenericSAMLTest(t,
		initialGenericSAMLConfig(), "initial_test")
	updatedConfig, updatedChecks := setupIntegrationIdPGenericSAMLTest(t,
		updatedGenericSAMLConfig(), "updated_test")

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

func setupIntegrationIdPGenericSAMLTest(t *testing.T, draft *GenericSAMLDraft, resName string) (
	string,
	resource.TestCheckFunc,
) {
	config := formatGenericSAMLDraftIntoConfig(draft, resName)

	nonEmptyRegex, err := regexp.Compile(".+")
	require.NoError(t, err)

	resourceFullName := fmt.Sprintf("cyral_integration_idp_generic_saml.%s", resName)
	checkFunc := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceFullName, "display_name",
			draft.DisplayName),
		resource.TestCheckResourceAttr(resourceFullName, "disable_idp_initiated_login",
			strconv.FormatBool(draft.DisableIdPInitiatedLogin)),
		resource.TestCheckResourceAttr(resourceFullName, "idp_type",
			draft.IdpType),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.first_name",
			draft.RequiredUserAttributes.FirstName.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.last_name",
			draft.RequiredUserAttributes.LastName.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.email",
			draft.RequiredUserAttributes.Email.Name),
		resource.TestCheckResourceAttr(resourceFullName, "attributes.0.groups",
			draft.RequiredUserAttributes.Groups.Name),
		resource.TestMatchResourceAttr(resourceFullName, "id",
			nonEmptyRegex),
		resource.TestMatchResourceAttr(resourceFullName, "sp_metadata",
			nonEmptyRegex),
	)

	return config, checkFunc
}

func formatGenericSAMLDraftIntoConfig(draft *GenericSAMLDraft, resName string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_generic_saml" "%s" {
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
		draft.RequiredUserAttributes.FirstName.Name,
		draft.RequiredUserAttributes.LastName.Name,
		draft.RequiredUserAttributes.Email.Name,
		draft.RequiredUserAttributes.Groups.Name,
	)
}
