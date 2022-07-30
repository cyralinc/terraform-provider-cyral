package cyral

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
<md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://sso-url-example.com/sso/saml"/>
</md:IDPSSODescriptor>
</md:EntityDescriptor>`, fakeCertificate)))
}

func TestAccIntegrationIdPSAMLResource(t *testing.T) {
	initialConfig, initialChecks := setupIntegrationIdPSAMLTest(
		"main_test", samlMetadataDocumentSample("fakeCertificateInitial"))
	updatedConfig, updatedChecks := setupIntegrationIdPSAMLTest(
		"main_test", samlMetadataDocumentSample("fakeCertificateUpdated"))

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
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"idp_metadata_document", "saml_draft_id"},
				ResourceName:            "cyral_integration_idp_saml.main_test",
			},
		},
	})
}

func setupIntegrationIdPSAMLTest(resName, metadataDoc string) (
	string,
	resource.TestCheckFunc,
) {
	var config string
	config += integrationIdPSAMLDraftResourceConfig(resName,
		"some-display-name", "some-idp-type")
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
			"idp_metadata_document", metadataDoc,
		),
	)

	return config, checkFunc
}

// integrationIdPSAMLDraftResourceConfig is a simplified version of
// formatGenericSAMLDraftIntoConfig. It only accepts custom display and idp
// type, and is used mostly to test the actual SAML integration (which needs a
// SAML draft to be created).
func integrationIdPSAMLDraftResourceConfig(resName, displayName, idpType string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml_draft" "%s" {
		display_name = "%s"
		idp_type = "%s"
	}`, resName, displayName, idpType)
}

func integrationIdPSAMLResourceConfig(resName, draftResName, metadataDoc string) string {
	return fmt.Sprintf(`
	resource "cyral_integration_idp_saml" "%s" {
		saml_draft_id = cyral_integration_idp_saml_draft.%s.id
		idp_metadata_document = "%s"
	}`, resName, draftResName, metadataDoc)
}
