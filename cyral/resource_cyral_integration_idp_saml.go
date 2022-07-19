package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type CreateGenericSAMLRequest struct {
	SAMLDraftId string                  `json:"samlDraftId"`
	IdpMetadata *GenericSAMLIdpMetadata `json:"idpMetadata,omitempty"`

	// Currently unused fields
	IdpDescriptor *GenericSAMLIdpDescriptor `json:"idpDescriptor,omitempty"`
}

func (req *CreateGenericSAMLRequest) ReadFromSchema(d *schema.ResourceData) error {
	req.SAMLDraftId = d.Get("saml_draft_id").(string)
	if url := d.Get("idp_metadata_url").(string); url != "" {
		req.IdpMetadata = &GenericSAMLIdpMetadata{
			URL: url,
		}
	} else if xml := d.Get("idp_metadata_document").(string); xml != "" {
		req.IdpMetadata = &GenericSAMLIdpMetadata{
			XML: xml,
		}
	}
	return nil
}

type CreateGenericSAMLResponse struct {
	Integration GenericSAMLIntegration `json:"integration"`
}

func (resp *CreateGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return resp.Integration.WriteToSchema(d)
}

func (resp *CreateGenericSAMLResponse) ReadFromSchema(d *schema.ResourceData) error {
	return resp.Integration.ReadFromSchema(d)
}

type ReadGenericSAMLResponse struct {
	IdentityProvider GenericSAMLIntegration `json:"identityProvider"`
}

func (resp *ReadGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return resp.IdentityProvider.WriteToSchema(d)
}

func (resp *ReadGenericSAMLResponse) ReadFromSchema(d *schema.ResourceData) error {
	return resp.IdentityProvider.ReadFromSchema(d)
}

type UpdateGenericSAMLRequest struct {
	ID               string                 `json:"id"`
	IdentityProvider GenericSAMLIntegration `json:"identityProvider"`
}

func (resp *UpdateGenericSAMLRequest) ReadFromSchema(d *schema.ResourceData) error {
	resp.ID = d.Id()
	return resp.IdentityProvider.ReadFromSchema(d)
}

func CreateGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso", c.ControlPlane)
		},
		NewResourceData: func() ResourceData { return &CreateGenericSAMLRequest{} },
		NewResponseData: func() ResponseData { return &CreateGenericSAMLResponse{} },
	}
}

func ReadGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
		NewResponseData: func() ResponseData { return &ReadGenericSAMLResponse{} },
	}
}

func UpdateGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceUpdate",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/sso/%s", c.ControlPlane, d.Id())
		},
		NewResourceData: func() ResourceData { return &UpdateGenericSAMLRequest{} },
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
	idpMetadataTypes := []string{"idp_metadata_url", "idp_metadata_document"}

	return &schema.Resource{
		Description: "Manages SAML IdP integrations.",
		CreateContext: CreateResource(
			CreateGenericSAMLConfig(),
			ReadGenericSAMLConfig(),
		),
		ReadContext: ReadResource(ReadGenericSAMLConfig()),
		UpdateContext: UpdateResource(
			UpdateGenericSAMLConfig(),
			ReadGenericSAMLConfig(),
		),
		DeleteContext: DeleteResource(DeleteGenericSAMLConfig()),
		Schema: map[string]*schema.Schema{
			"saml_draft_id": {
				Description:  "A valid id for a SAML Draft. Must be at least 5 character long. See attribute `id` in resource `cyral_integration_idp_saml_draft`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validationStringLenAtLeast(5),
			},
			"idp_metadata_url": {
				Description:  "A SAML XML IdP Metadata document containing all configuration values required by the Cyral SP. Conflicts with `idp_metadata_document`.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: idpMetadataTypes,
			},
			"idp_metadata_document": {
				Description:  "Full SAML metadata XML document. Must be base64 encoded. Conflicts with `idp_metadata_url`.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: idpMetadataTypes,
				ValidateFunc: validation.StringIsBase64,
			},
			// TODO: computed values
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
