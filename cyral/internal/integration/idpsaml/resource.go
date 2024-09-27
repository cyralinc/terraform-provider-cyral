package idpsaml

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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

func CreateGenericSAMLConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "GenericSAMLResourceCreate",
		Type:         operationtype.Create,
		HttpMethod:   http.MethodPost,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso", c.ControlPlane)
		},
		SchemaReaderFactory: func() core.SchemaReader { return &CreateGenericSAMLRequest{} },
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &CreateGenericSAMLResponse{} },
	}
}

func CreateIdPConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "GenericSAMLResourceValidation",
		Type:         operationtype.Create,
		HttpMethod:   http.MethodPost,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/conf/identityProviders/%s", c.ControlPlane, d.Id())
		},
		SchemaReaderFactory: func() core.SchemaReader { return &IdentityProviderData{} },
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &IdentityProviderData{} },
	}
}

func ReadGenericSAMLConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "GenericSAMLResourceRead",
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &ReadGenericSAMLResponse{} },
		RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Generic SAML"},
	}
}

func DeleteGenericSAMLConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "GenericSAMLResourceDelete",
		Type:         operationtype.Delete,
		HttpMethod:   http.MethodDelete,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
	}
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages identity provider (IdP) integrations using SAML to allow " +
			"[Single Sing-On](https://cyral.com/docs/integrations/authentication/idp/) to Cyral.\n\nSee also " +
			"the remaining SAML-related resources and data sources.",
		CreateContext: core.CRUDResources(
			[]core.ResourceOperationConfig{
				CreateGenericSAMLConfig(),
				ReadGenericSAMLConfig(),
				CreateIdPConfig(),
			},
		),
		ReadContext:   core.ReadResource(ReadGenericSAMLConfig()),
		DeleteContext: core.DeleteResource(DeleteGenericSAMLConfig()),
		Schema: map[string]*schema.Schema{
			// Input arguments
			//
			"saml_draft_id": {
				Description:  "A valid id for a SAML Draft. Must be at least 5 character long. See attribute `id` in resource `cyral_integration_idp_saml_draft`.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: utils.ValidationStringLenAtLeast(5),
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
