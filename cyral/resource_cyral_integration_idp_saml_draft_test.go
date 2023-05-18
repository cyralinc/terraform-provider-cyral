package cyral

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	integrationIdPSAMLDraftResourceName = "integration-idp-saml-draft"
)

func genericSAMLDraftConfigInitial() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName: accTestName(
			integrationIdPSAMLDraftResourceName, "integration-1"),
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

func genericSAMLDraftConfigUpdated() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName: accTestName(
			integrationIdPSAMLDraftResourceName, "integration-2"),
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

func genericSAMLDraftConfigNoAttributes() *GenericSAMLDraft {
	return &GenericSAMLDraft{
		DisplayName: accTestName(
			integrationIdPSAMLDraftResourceName, "integration-2"),
		DisableIdPInitiatedLogin: true,
		IdpType:                  "some-idp-type-2",
	}
}

func TestAccIntegrationIdPSAMLDraftResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPSAMLDraftTest(t,
		genericSAMLDraftConfigInitial(), "main_test")
	updatedConfig, updatedChecks := setupIntegrationIdPSAMLDraftTest(t,
		genericSAMLDraftConfigUpdated(), "main_test")
	noAttributesConfig, noAttributesChecks := setupIntegrationIdPSAMLDraftTest(t,
		genericSAMLDraftConfigNoAttributes(), "no_attributes")

	resource.ParallelTest(t, resource.TestCase{
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
			{
				Config: noAttributesConfig,
				Check:  noAttributesChecks,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"attributes.", "toggle_recreation"},
				ResourceName:            "cyral_integration_idp_saml_draft.no_attributes",
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
	checkFuncs := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceFullName, "display_name",
			draft.DisplayName),
		resource.TestCheckResourceAttr(resourceFullName, "disable_idp_initiated_login",
			strconv.FormatBool(draft.DisableIdPInitiatedLogin)),
		resource.TestCheckResourceAttr(resourceFullName, "idp_type",
			draft.IdpType),
		resource.TestMatchResourceAttr(resourceFullName, "id",
			nonEmptyRegex),
		resource.TestMatchResourceAttr(resourceFullName, "sp_metadata",
			nonEmptyRegex),
	}

	if draft.Attributes != nil {
		checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
			resource.TestCheckResourceAttr(resourceFullName, "attributes.0.first_name",
				draft.Attributes.FirstName.Name),
			resource.TestCheckResourceAttr(resourceFullName, "attributes.0.last_name",
				draft.Attributes.LastName.Name),
			resource.TestCheckResourceAttr(resourceFullName, "attributes.0.email",
				draft.Attributes.Email.Name),
			resource.TestCheckResourceAttr(resourceFullName, "attributes.0.groups",
				draft.Attributes.Groups.Name),
		}...)
	}

	// checking SPMetadata content
	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.xml_document"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.entity_id"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.url"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.single_logout_url"),
	}...)
	checkFuncs = append(checkFuncs, []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.assertion_consumer_services.0.url"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.assertion_consumer_services.0.index"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.assertion_consumer_services.1.url"),
		resource.TestCheckResourceAttrSet(resourceFullName, "service_provider_metadata.0.assertion_consumer_services.1.index"),
	}...)

	return config, resource.ComposeTestCheckFunc(checkFuncs...)
}

func formatGenericSAMLDraftIntoConfig(draft *GenericSAMLDraft, resName string) string {
	var attributesStr string
	if draft.Attributes != nil {
		attributesStr += fmt.Sprintf(`
		attributes {
			first_name = "%s"
			last_name = "%s"
			email = "%s"
			groups = "%s"
		}`,
			draft.Attributes.FirstName.Name,
			draft.Attributes.LastName.Name,
			draft.Attributes.Email.Name,
			draft.Attributes.Groups.Name,
		)
	}
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		disable_idp_initiated_login = %t
		idp_type = "%s"
		%s
	}`, resName, draft.DisplayName, draft.DisableIdPInitiatedLogin,
		draft.IdpType, attributesStr)
}
