package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

func GenericSAMLIdPInfoArguments() []string {
	return []string{"idp_metadata_url", "idp_metadata_xml", "idp_descriptor"}
}

type CreateGenericSAMLRequest struct {
	SAMLDraftId   string                    `json:"samlDraftId"`
	IdpMetadata   *GenericSAMLIdpMetadata   `json:"idpMetadata,omitempty"`
	IdpDescriptor *GenericSAMLIdpDescriptor `json:"idpDescriptor,omitempty"`
}

func (req *CreateGenericSAMLRequest) ReadFromSchema(d *schema.ResourceData) error {
	req.SAMLDraftId = d.Get("saml_draft_id").(string)
	if url := d.Get("idp_metadata_url").(string); url != "" {
		req.IdpMetadata = &GenericSAMLIdpMetadata{
			URL: url,
		}
	} else if xml := d.Get("idp_metadata_xml").(string); xml != "" {
		req.IdpMetadata = &GenericSAMLIdpMetadata{
			XML: xml,
		}
	} else if idpDescriptorList := d.Get("idp_descriptor").(*schema.Set).List(); len(idpDescriptorList) > 0 {
		idpDescriptorMap := idpDescriptorList[0].(map[string]interface{})
		req.IdpDescriptor = &GenericSAMLIdpDescriptor{
			SingleSignOnServiceURL:     idpDescriptorMap["single_sign_on_service_url"].(string),
			SigningCertificate:         idpDescriptorMap["signing_certificate"].(string),
			DisableForceAuthentication: idpDescriptorMap["disable_force_authentication"].(bool),
			SingleLogoutServiceURL:     idpDescriptorMap["single_logout_service_url"].(string),
		}
	} else {
		panic(fmt.Sprintf("Expected one of the arguments to be set: %v.",
			GenericSAMLIdPInfoArguments()))
	}
	return nil
}

type CreateGenericSAMLResponse struct {
	Integration GenericSAMLIntegration `json:"integration"`
}

func (resp *CreateGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return resp.Integration.WriteToSchema(d)
}

type ReadGenericSAMLResponse struct {
	IdentityProvider GenericSAMLIntegration `json:"identityProvider"`
}

func (resp *ReadGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return resp.IdentityProvider.WriteToSchema(d)
}

func CreateGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso", c.ControlPlane)
		},
		NewResourceData: func() ResourceData { return &CreateGenericSAMLRequest{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &CreateGenericSAMLResponse{} },
	}
}

func ReadGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ReadGenericSAMLResponse{} },
	}
}

func DeleteGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
	}
}

func resourceIntegrationIdPSAML() *schema.Resource {
	return &schema.Resource{
		Description: "Manages SAML IdP integrations.",
		CreateContext: CreateResource(
			CreateGenericSAMLConfig(),
			ReadGenericSAMLConfig(),
		),
		ReadContext:   ReadResource(ReadGenericSAMLConfig()),
		DeleteContext: DeleteResource(DeleteGenericSAMLConfig()),
		Schema: map[string]*schema.Schema{
			// Input arguments
			//
			"saml_draft_id": {
				Description:  "A valid id for a SAML Draft. Must be at least 5 character long. See attribute `id` in resource `cyral_integration_idp_saml_draft`.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validationStringLenAtLeast(5),
			},
			// Do not allow in-place updates of idp_metadata_url or
			// idp_metadata_xml. This would require us to parse the
			// URL or XML locally to obtain the SSO URL,
			// certificate, etc. Thus it would be a huge pain in the
			// back to keep in sync with the API backend.
			//
			// TODO: in a future implementation we should allow
			// in-place updates so as to avoid problems when
			// updating these parameters. -aholmquist 2022-08-04
			"idp_metadata_url": {
				Description:  "The web address of an IdP SAML Metadata XML document. Conflicts with `idp_metadata_xml`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: GenericSAMLIdPInfoArguments(),
			},
			"idp_metadata_xml": {
				Description:  "Full SAML metadata XML document. Must be base64 encoded. Conflicts with `idp_metadata_url`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: GenericSAMLIdPInfoArguments(),
				ValidateFunc: validation.StringIsBase64,
			},
			"idp_descriptor": {
				Description: "The configuration information required by the Cyral SP, provided by the IdP.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"single_sign_on_service_url": {
							Description:  "The IdP’s Single Sign-on Service (SS0) URL, where Cyral SP will send SAML AuthnRequests via SAML-POST binding.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validationStringPrefix("https://"),
						},
						"signing_certificate": {
							Description: "The signing certificate used by the Cyral SP to validate signed SAML assertions sent by the IdP.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"disable_force_authentication": {
							Description: "Indicates whether the identity provider must authenticate the presenter directly rather than rely on a previous security context.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     false,
						},
						"single_logout_service_url": {
							Description: "The IdP’s Single Log-out Service (SL0) URL, where Cyral will send SAML AuthnRequests via SAML-POST binding. If supplied, SLO will be enabled.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},

			// Computed arguments
			//
			"id": {
				Description: "ID of this resource in the Cyral environment.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Computed:    true,
			},
			"single_sign_on_service_url": {
				Description: "The IdP’s Single Sign-on Service (SSO) URL, where Cyral SP will send SAML AuthnRequests via SAML-POST binding.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
