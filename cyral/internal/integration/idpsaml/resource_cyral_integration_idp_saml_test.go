package idpsaml_test

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"

	"github.com/cyralinc/terraform-provider-cyral/cyral/provider"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	integrationIdPSAMLResourceName = "integration-idp-saml"

	testSSOURL = "https://sso-url-example.com/sso/saml"
)

func samlMetadataDocumentSample(fakeCertificate string) string {
	// Do not add sensitive information here!
	//
	// The XML document is sanity-checked by the API, so we need a sample
	// for ACC tests. Every XML element below is necessary.
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(`
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="http://www.entity-id-example.com/1234567890">
<md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
<md:KeyDescriptor use="signing">
<ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
<ds:X509Data>
<ds:X509Certificate>
%s
</ds:X509Certificate>
</ds:X509Data>
</ds:KeyInfo>
</md:KeyDescriptor>
<md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="%s"/>
</md:IDPSSODescriptor>
</md:EntityDescriptor>`, fakeCertificate, testSSOURL)))
}

func TestAccIntegrationIdPSAMLResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPSAMLTest(
		"upgrade_test", samlMetadataDocumentSample("fakeCertificateInitial"))
	updatedConfig, _ := setupIntegrationIdPSAMLTest(
		"upgrade_test", samlMetadataDocumentSample("fakeCertificateUpdated"))
	updatedAgainConfig, updatedAgainChecks := setupIntegrationIdPSAMLTest(
		"upgrade_test", samlMetadataDocumentSample("fakeCertificateUpdated"))
	newConfig, newChecks := setupIntegrationIdPSAMLTest(
		"new_test", samlMetadataDocumentSample("fakeCertificateNew"))

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: provider.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check:  initialChecks,
			},
			{
				Config: updatedConfig,
				// This is a known case where the resource
				// fails: when creating a new SAML integration
				// pointing to an old, completed, SAML draft.
				ExpectError: regexp.MustCompile("Error running apply: exit status 1"),
			},
			{
				Config: updatedAgainConfig,
				// If user runs apply again, it should work.
				Check: updatedAgainChecks,
			},
			{
				Config: newConfig,
				// When a new SAML draft and a new integration
				// are created, there should be no no problem.
				Check: newChecks,
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"idp_metadata_xml", "saml_draft_id"},
				ResourceName:            "cyral_integration_idp_saml.new_test",
			},
		},
	})
}

func setupIntegrationIdPSAMLTest(resName, metadataDoc string) (
	string,
	resource.TestCheckFunc,
) {
	var config string
	config += utils.FormatBasicIntegrationIdPSAMLDraftIntoConfig(resName,
		utils.AccTestName(integrationIdPSAMLResourceName, "saml-draft"),
		"some-idp-type")
	config += integrationIdPSAMLResourceConfig(resName, resName, metadataDoc)

	resourceFullName := fmt.Sprintf("cyral_integration_idp_saml.%s", resName)
	resourceFullNameDraft := fmt.Sprintf("cyral_integration_idp_saml_draft.%s",
		resName)
	checkFunc := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			resourceFullName, "saml_draft_id",
			resourceFullNameDraft, "id",
		),
		// Unfortunately, we can only test this resource using a
		// metadata document. Using a URL would require an active
		// external SAML endpoint during the ACC tests.
		resource.TestCheckResourceAttr(resourceFullName,
			"idp_metadata_xml", metadataDoc,
		),
		resource.TestCheckResourceAttr(resourceFullName,
			"single_sign_on_service_url", testSSOURL,
		),
	)

	return config, checkFunc
}

func integrationIdPSAMLResourceConfig(resName, draftResName, metadataDoc string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml" "%s" {
		saml_draft_id = cyral_integration_idp_saml_draft.%s.id
		idp_metadata_xml = "%s"
	}`, resName, draftResName, metadataDoc)
}
