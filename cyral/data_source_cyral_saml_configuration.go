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
		ReadContext: dataSourceSAMLConfigurationRead,
		Schema: map[string]*schema.Schema{
			"saml_metadata_url": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: metadataTypes,
			},
			"base_64_saml_metadata_document": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: metadataTypes,
			},
			"disable_using_jwks_url": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sync_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_id_policy_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"principal_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"saml_xml_key_name_tranformer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hide_on_login_page": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"back_channel_supported": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_response": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_authn_request": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_logout": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_want_authn_requests_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_want_assertions_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"want_assertions_encrypted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_force_authentication": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_validate_signature": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gui_order": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"single_sign_on_service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"single_logout_service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"xml_sig_key_info_key_name_transformer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allowed_clock_skew": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ldap_group_attribute": {
				Type:     schema.TypeString,
				Computed: true,
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
