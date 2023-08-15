package cyral

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/require"
)

const (
	integrationIdPSAMLDataSourceName = "data-integration-idp-saml"
)

func integrationIdPSAMLDataSourceTestIdps() []GenericSAMLIntegration {
	return []GenericSAMLIntegration{
		{
			ID: "id-1",
			DisplayName: accTestName(
				integrationIdPSAMLDataSourceName, "integration-1"),
			IdpType:  "idp-type-1",
			Disabled: false,
			IdpDescriptor: &GenericSAMLIdpDescriptor{
				SingleSignOnServiceURL:     "sso-url-1",
				SigningCertificate:         "signing-certificate-1",
				DisableForceAuthentication: false,
				SingleLogoutServiceURL:     "slo-url-1",
			},
			SPMetadata: &GenericSAMLSPMetadata{
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
			ID: "id-2",
			DisplayName: accTestName(
				integrationIdPSAMLDataSourceName, "integration-2"),
			IdpType:  "idp-type-2",
			Disabled: true,
			IdpDescriptor: &GenericSAMLIdpDescriptor{
				SingleSignOnServiceURL:     "sso-url-2",
				SigningCertificate:         "signing-certificate-2",
				DisableForceAuthentication: true,
				SingleLogoutServiceURL:     "slo-url-2",
			},
			SPMetadata: &GenericSAMLSPMetadata{
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

func testIntegrationIdPSAMLDataSourceName1() string {
	return accTestName(integrationIdPSAMLDataSourceName, "1")
}

func testIntegrationIdPSAMLDataSourceName2() string {
	return accTestName(integrationIdPSAMLDataSourceName, "2")
}

func TestAccIntegrationIdPSAMLDataSource(t *testing.T) {
	testConfig1, testFunc1 := testIntegrationIdPSAMLDataSource(t,
		"test1", testIntegrationIdPSAMLDataSourceName1(), "type1")
	testConfig2, testFunc2 := testIntegrationIdPSAMLDataSource(t,
		"test2", testIntegrationIdPSAMLDataSourceName2(), "type2")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig1,
				Check:  testFunc1,
			},
			{
				Config: testConfig2,
				Check:  testFunc2,
			},
		},
	})
}

func testIntegrationIdPSAMLDataSource(t *testing.T, resName string, nameFilter, typeFilter string) (
	string, resource.TestCheckFunc,
) {
	return testIntegrationIdPSAMLDataSourceConfig(resName, nameFilter, typeFilter),
		testIntegrationIdPSAMLDataSourceChecks(t, resName, nameFilter, typeFilter)
}

// Setup two integrations that are retrieved and checked by the data source.
func testIntegrationIdPSAMLDataSourceConfigDependencies(resName string) string {
	resName1 := resName + "_1"
	resName2 := resName + "_2"

	var config string
	config += formatBasicIntegrationIdPSAMLDraftIntoConfig(
		resName1,
		testIntegrationIdPSAMLDataSourceName1(),
		"type1")
	config += integrationIdPSAMLResourceConfig(
		resName1,
		resName1,
		samlMetadataDocumentSample("fake-certificate"))
	config += formatBasicIntegrationIdPSAMLDraftIntoConfig(
		resName2,
		testIntegrationIdPSAMLDataSourceName2(),
		"type2")
	config += integrationIdPSAMLResourceConfig(
		resName2,
		resName2,
		samlMetadataDocumentSample("fake-certificate"))
	return config
}

func testIntegrationIdPSAMLDataSourceConfig(resName, nameFilter, typeFilter string) string {
	var config string
	config += testIntegrationIdPSAMLDataSourceConfigDependencies(resName)
	config += integrationIdPSAMLDataSourceConfig(
		resName,
		[]string{
			fmt.Sprintf("cyral_integration_idp_saml.%s_1", resName),
			fmt.Sprintf("cyral_integration_idp_saml.%s_2", resName),
		},
		nameFilter,
		typeFilter)
	return config
}

// The checks assume that there exists two SAML integrations in the Terraform
// state, but only one passed the filter.
func testIntegrationIdPSAMLDataSourceChecks(t *testing.T, resName, nameFilter, typeFilter string) resource.TestCheckFunc {
	dataSourceFullName := fmt.Sprintf("data.cyral_integration_idp_saml.%s", resName)

	nonEmptyRegex, err := regexp.Compile(".+")
	require.NoError(t, err)

	testFunctions := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.#", "1",
		),
		resource.TestMatchResourceAttr(dataSourceFullName,
			"idp_list.0.id", nonEmptyRegex,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.disabled", "false",
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.idp_descriptor.#", "1",
		),
		resource.TestMatchResourceAttr(dataSourceFullName,
			"idp_list.0.idp_descriptor.0.single_sign_on_service_url", nonEmptyRegex,
		),
		resource.TestMatchResourceAttr(dataSourceFullName,
			"idp_list.0.idp_descriptor.0.signing_certificate", nonEmptyRegex,
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.sp_metadata.#", "1",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.xml_document",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.url",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.entity_id",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.single_logout_url",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.assertion_consumer_services.#",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.assertion_consumer_services.0.url",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.assertion_consumer_services.0.index",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.assertion_consumer_services.1.url",
		),
		resource.TestCheckResourceAttrSet(dataSourceFullName,
			"idp_list.0.sp_metadata.0.assertion_consumer_services.1.index",
		),
		resource.TestCheckResourceAttr(dataSourceFullName,
			"idp_list.0.attributes.#", "1",
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

	if nameFilter != "" {
		testFunctions = append(testFunctions,
			resource.TestCheckResourceAttr(dataSourceFullName,
				"idp_list.0.display_name", nameFilter,
			),
		)
	}

	return resource.ComposeTestCheckFunc(testFunctions...)
}

func filterSAMLIdps(idps []GenericSAMLIntegration, nameFilter, typeFilter string) []GenericSAMLIntegration {
	var filteredIdps []GenericSAMLIntegration
	for _, idp := range idps {
		if (nameFilter == "" || idp.DisplayName == nameFilter) &&
			(typeFilter == "" || idp.IdpType == typeFilter) {
			filteredIdps = append(filteredIdps, idp)
		}
	}
	return filteredIdps
}

func integrationIdPSAMLDataSourceConfig(resName string, dependsOn []string,
	nameFilter, typeFilter string,
) string {
	return fmt.Sprintf(`
	data "cyral_integration_idp_saml" "%s" {
		depends_on = %s
		display_name = "%s"
		idp_type = "%s"
	}`, resName, utils.ListToStr(dependsOn), nameFilter, typeFilter)
}
