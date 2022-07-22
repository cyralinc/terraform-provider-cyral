package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

func integrationIdPSAMLDataSourceTestIdps() []GenericSAMLIntegration {
	return []GenericSAMLIntegration{
		{
			ID:          "id-1",
			DisplayName: "display-name-1",
			IdpType:     "idp-type-1",
			Disabled:    false,
			IdpDescriptor: &GenericSAMLIdpDescriptor{
				SingleSignOnServiceURL:     "sso-url-1",
				SigningCertificate:         "signing-certificate-1",
				DisableForceAuthentication: false,
				SingleLogoutServiceURL:     "slo-url-1",
			},
			SPMetadata: &SPMetadata{
				XMLDocument: "xml-document-1",
			},
			Attributes: NewRequiredUserAttributes(
				"first-name-1",
				"last-name-1",
				"email-1",
				"groups-1",
			),
		},
		{
			ID:          "id-2",
			DisplayName: "display-name-2",
			IdpType:     "idp-type-2",
			Disabled:    true,
			IdpDescriptor: &GenericSAMLIdpDescriptor{
				SingleSignOnServiceURL:     "sso-url-2",
				SigningCertificate:         "signing-certificate-2",
				DisableForceAuthentication: true,
				SingleLogoutServiceURL:     "slo-url-2",
			},
			SPMetadata: &SPMetadata{
				XMLDocument: "xml-document-2",
			},
			Attributes: NewRequiredUserAttributes(
				"first-name-2",
				"last-name-2",
				"email-2",
				"groups-2",
			),
		},
	}
}

func TestAccIntegrationIdPSAMLDataSource(t *testing.T) {
	testConfigNameFilter, testFuncNameFilter := testIntegrationIdPSAMLDataSource(t,
		"name_filter", "name_filter_1", "")
	testConfigTypeFilter, testFuncTypeFilter := testIntegrationIdPSAMLDataSource(t,
		"type_filter", "", "type_filter_2")

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigNameFilter,
				Check:  testFuncNameFilter,
			},
			{
				Config: testConfigTypeFilter,
				Check:  testFuncTypeFilter,
			},
		},
	})
}

func testIntegrationIdPSAMLDataSource(t *testing.T, resName string, displayName, idpType string) (
	string, resource.TestCheckFunc,
) {
	return testIntegrationIdPSAMLDataSourceConfig(resName, displayName, idpType),
		testIntegrationIdPSAMLDataSourceChecks(t, resName, displayName, idpType)
}

func testIntegrationIdPSAMLDataSourceConfig(resName, displayName, idpType string) string {
	var config string
	// Setup two integrations
	config += integrationIdPSAMLDraftResourceConfig(
		fmt.Sprintf("%s_1", resName),
		fmt.Sprintf("%s_1", resName),
		fmt.Sprintf("%s_1", resName))
	config += integrationIdPSAMLResourceConfig(
		fmt.Sprintf("%s_1", resName),
		fmt.Sprintf("%s_1", resName),
		samlMetadataDocumentSample("fake-certificate"))
	config += integrationIdPSAMLDraftResourceConfig(
		fmt.Sprintf("%s_2", resName),
		fmt.Sprintf("%s_2", resName),
		fmt.Sprintf("%s_2", resName))
	config += integrationIdPSAMLResourceConfig(
		fmt.Sprintf("%s_2", resName),
		fmt.Sprintf("%s_2", resName),
		samlMetadataDocumentSample("fake-certificate"))
	config += integrationIdPSAMLDataSourceConfig(
		resName,
		[]string{
			fmt.Sprintf("cyral_integration_idp_saml.%s_1", resName),
			fmt.Sprintf("cyral_integration_idp_saml.%s_2", resName),
		},
		displayName,
		idpType)
	return config
}

// The checks assume that there exists two SAML integrations in the Terraform
// state, but only one passed the filter.
func testIntegrationIdPSAMLDataSourceChecks(t *testing.T, resName, displayName, idpType string) resource.TestCheckFunc {
	dataSourceFullName := fmt.Sprintf("data.cyral_integration_idp_saml.%s", resName)

	nonEmptyRegex, err := regexp.Compile(".+")
	require.NoError(t, err)

	testFunctions := []resource.TestCheckFunc{
		resource.TestMatchResourceAttr(dataSourceFullName,
			"idp_list.0.id", nonEmptyRegex,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.disabled", "false",
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.attributes.0.first_name", defaultUserAttributeFirstName,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.attributes.0.last_name", defaultUserAttributeLastName,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.attributes.0.email", defaultUserAttributeEmail,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.attributes.0.groups", defaultUserAttributeGroups,
		),
	}

	if displayName != "" {
		testFunctions = append(testFunctions,
			resource.TestCheckResourceAttr(dataSourceFullName,
				"idp_list.0.display_name", displayName,
			),
		)
	}
	if idpType != "" {
		testFunctions = append(testFunctions,
			resource.TestCheckResourceAttr(dataSourceFullName,
				"idp_list.0.idp_type", idpType,
			),
		)
	}

	return resource.ComposeTestCheckFunc(testFunctions...)
}

func filterSAMLIdps(idps []GenericSAMLIntegration, displayName, idpType string) []GenericSAMLIntegration {
	var filteredIdps []GenericSAMLIntegration
	for _, idp := range idps {
		if (displayName == "" || idp.DisplayName == displayName) &&
			(idpType == "" || idp.IdpType == idpType) {
			filteredIdps = append(filteredIdps, idp)
		}
	}
	return filteredIdps
}

func integrationIdPSAMLDataSourceConfig(resName string, dependsOn []string,
	displayName, idpType string,
) string {
	return fmt.Sprintf(`
	data "cyral_integration_idp_saml" "%s" {
		depends_on = [%s]
		display_name = "%s"
		idp_type = "%s"
	}`, resName, formatAttributes(dependsOn), displayName, idpType)
}
