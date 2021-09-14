package cyral

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSAMLConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSAMLConfigRead,
		Schema: map[string]*schema.Schema{
			"saml_metadata_url": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"base_64_saml_metadata_document",
				},
			},
			"base_64_saml_metadata_document": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"saml_metadata_url",
				},
			},
			"allowed_clock_skew": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"back_channel_supported": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_force_authentication": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_authn_request": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_response": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_post_binding_logout": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_using_jwks_url": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_validate_signature": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_want_assertions_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disable_want_authn_requests_signed": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gui_order": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hide_on_login_page": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ldap_group_attribute": {
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
			"saml_xml_key_name_transformer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signature_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signing_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"single_logout_service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"single_sign_on_service_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sync_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"want_assertions_encrypted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"xml_sig_key_info_key_name_transformer": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSAMLConfigRead(c context.Context, rd *schema.ResourceData, i interface{}) diag.Diagnostics {
	var doc string
	if url := rd.Get("saml_metadata_url").(string); url != "" {
		resp, err := http.Get(url)
		if err != nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode != 200 {
			return diag.Errorf("response from metadata url not 200. Status: %v", resp.Status)
		}
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return diag.FromErr(err)
		}
		doc = string(buf)
	} else if docb := rd.Get("base_64_saml_metadata_document").(string); docb != "" {
		doc = docb
	}
	req := ParseRequest{doc}

	return ReadResource(
		ResourceOperationConfig{
			Name:       "SAMLConfigurationReadResource",
			HttpMethod: http.MethodGet, CreateURL: func(d *schema.ResourceData, c *client.Client) string {
				return fmt.Sprintf("https://%v/v1/integrations/saml/parse", c.ControlPlane)
			}, ResourceData: &req, ResponseData: &SAMLConfig{}},
	)(c, rd, i)
}
