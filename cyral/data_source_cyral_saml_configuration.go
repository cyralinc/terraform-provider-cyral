package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	metadataTypes = []string{"base_64_saml_metadata_document", "saml_metadata_url"}
)

func dataSourceSAMLConfiguration() *schema.Resource {
	return &schema.Resource{
		Description: "Parses a SAML metadata URL or a Base64 document into a SAML configuration.",
		ReadContext: dataSourceSAMLConfigurationRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Computed ID for this resource (locally computed to be used in Terraform state).",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"saml_metadata_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: metadataTypes,
				Description: "(Required unless using `base_64_saml_metadata_document`) This is the full SAML metadata URL we " +
					"should use to parse to a SAML configuration.",
			},
			"base_64_saml_metadata_document": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: metadataTypes,
				Description: "(Required unless using `saml_metadata_url`) This is the full SAML metadata document that should " +
					"be used to parse a SAML configuration, Base64 encoded.",
			},
			"disable_using_jwks_url": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "By default, the jwks URL is used for all SAML connections.",
			},
			"sync_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `FORCE` if unset.",
			},
			"name_id_policy_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified` if unset.",
			},
			"principal_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `SUBJECT` if unset.",
			},
			"signature_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `RSA_SHA256` if unset.",
			},
			"saml_xml_key_name_tranformer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `KEY_ID` if unset.",
			},
			"hide_on_login_page": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Defaults to `false` if unset.",
			},
			"back_channel_supported": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Defaults to `false` if unset.",
			},
			"disable_post_binding_response": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` " +
					"binding will be used.",
			},
			"disable_post_binding_authn_request": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Indicates whether the AuthnRequest must be sent using `HTTP-POST` binding. If `true`, " +
					"`HTTP-REDIRECT` binding will be used.",
			},
			"disable_post_binding_logout": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Indicates whether to respond to requests using `HTTP-POST` binding. If `true`, `HTTP-REDIRECT` " +
					"binding will be used.",
			},
			"disable_want_authn_requests_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the identity provider expects a signed AuthnRequest.",
			},
			"disable_want_assertions_signed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the service provider expects a signed Assertion.",
			},
			"want_assertions_encrypted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether the service provider expects an encrypted Assertion.",
			},
			"disable_force_authentication": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Indicates whether the identity provider must authenticate the presenter directly rather than rely " +
					"on a previous security context.",
			},
			"disable_validate_signature": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable/Disable signature validation of SAML responses. Highly recommended for minimum security.",
			},
			"gui_order": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUI order.",
			},
			"single_sign_on_service_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL that must be used to send authentication requests (SAML AuthnRequest).",
			},
			"single_logout_service_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL that must be used to send logout requests.",
			},
			"xml_sig_key_info_key_name_transformer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defaults to `KEY_ID` if unset.",
			},
			"signing_certificate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Signing certificate used to validate signatures. Required if signature validation is enabled.",
			},
			"allowed_clock_skew": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Clock skew in seconds that is tolerated when validating identity provider tokens. Default value is `0`.",
			},
			"ldap_group_attribute": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Type of `LDAP Group RDN` that identifies the name of a group within a DN. For example, if an " +
					"LDAP DN sent in a SAML assertion is `cn=Everyone`, `ou=groups`, `dc=openam`, `dc=forgerock`, `dc=org` and " +
					"the `LDAP Group RDN` Type is `cn` Cyral will interpret `Everyone` as the group name.",
			},
		},
	}
}

func dataSourceSAMLConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceSAMLConfigurationRead")
	c := m.(*client.Client)

	metadataRequest := getSAMLMetadataRequestFromSchema(d)

	url := fmt.Sprintf("https://%v/v1/integrations/saml/parse", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, metadataRequest)
	if err != nil {
		return createError("Unable to retrieve saml configuration", fmt.Sprintf("%v", err))
	}

	response := SAMLConfiguration{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(uuid.New().String())
	setSAMLConfigurationToSchema(d, response)

	log.Printf("[DEBUG] End dataSourceSAMLConfigurationRead")

	return diag.Diagnostics{}
}

func getSAMLMetadataRequestFromSchema(d *schema.ResourceData) ParseSAMLMetadataRequest {
	return ParseSAMLMetadataRequest{
		SamlMetadataURL:            d.Get("saml_metadata_url").(string),
		Base64SamlMetadataDocument: d.Get("base_64_saml_metadata_document").(string),
	}
}

func setSAMLConfigurationToSchema(d *schema.ResourceData, data SAMLConfiguration) {
	if data.Config != nil {
		d.Set("disable_using_jwks_url", data.Config.DisableUsingJWKSUrl)
		d.Set("sync_mode", data.Config.SyncMode)
		d.Set("name_id_policy_format", data.Config.NameIDPolicyFormat)
		d.Set("principal_type", data.Config.PrincipalType)
		d.Set("signature_type", data.Config.SignatureType)
		d.Set("saml_xml_key_name_tranformer", data.Config.SamlXmlKeyNameTranformer)
		d.Set("hide_on_login_page", data.Config.HideOnLoginPage)
		d.Set("back_channel_supported", data.Config.BackChannelSupported)
		d.Set("disable_post_binding_response", data.Config.DisablePostBindingResponse)
		d.Set("disable_post_binding_authn_request", data.Config.DisablePostBindingAuthnRequest)
		d.Set("disable_post_binding_logout", data.Config.DisablePostBindingLogout)
		d.Set("disable_want_authn_requests_signed", data.Config.DisableWantAuthnRequestsSigned)
		d.Set("disable_want_assertions_signed", data.Config.DisableWantAssertionsSigned)
		d.Set("want_assertions_encrypted", data.Config.WantAssertionsEncrypted)
		d.Set("disable_force_authentication", data.Config.DisableForceAuthentication)
		d.Set("disable_validate_signature", data.Config.DisableValidateSignature)
		d.Set("gui_order", data.Config.GuiOrder)
		d.Set("single_sign_on_service_url", data.Config.SingleSignOnServiceURL)
		d.Set("single_logout_service_url", data.Config.SingleLogoutServiceURL)
		d.Set("xml_sig_key_info_key_name_transformer", data.Config.XmlSigKeyInfoKeyNameTransformer)
		d.Set("signing_certificate", data.Config.SigningCertificate)
		d.Set("allowed_clock_skew", data.Config.AllowedClockSkew)
		d.Set("saml_metadata_url", data.Config.SamlMetadataURL)
		d.Set("base_64_saml_metadata_document", data.Config.Base64SamlMetadataDocument)
		d.Set("ldap_group_attribute", data.Config.LdapGroupAttribute)
	}
}
