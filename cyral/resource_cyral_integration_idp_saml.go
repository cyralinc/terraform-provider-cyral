package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

func GenericSAMLIdPInfoArguments() []string {
	return []string{"idp_metadata_url", "idp_metadata_xml"}
}

type CreateGenericSAMLRequest struct {
	SAMLDraftId string                  `json:"samlDraftId"`
	IdpMetadata *GenericSAMLIdpMetadata `json:"idpMetadata,omitempty"`

	// Currently unused fields. TODO: implement this -aholmquist 2022-08-03
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

func ValidateGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceValidation",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Id())
		},
		NewResourceData: func() ResourceData { return &IdentityProviderData{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &IdentityProviderData{} },
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
		Description: "Manages identity provider (IdP) integrations using SAML to allow " +
			"[Single Sing-On](https://cyral.com/docs/sso/overview) to Cyral.\n\nSee also " +
			"the remaining SAML-related resources and data sources.",
		CreateContext: CRUDResources(
			[]ResourceConfig{
				{
					Type:            create,
					OperationConfig: CreateGenericSAMLConfig(),
				},
				{
					Type:            read,
					OperationConfig: ReadGenericSAMLConfig(),
				},
				{
					Type:            update,
					OperationConfig: ValidateGenericSAMLConfig(),
				},
			},
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
