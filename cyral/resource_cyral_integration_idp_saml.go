package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type CreateGenericSAMLRequest struct {
	SAMLDraftId   string                    `json:"samlDraftId"`
	IdpDescriptor *GenericSAMLIdpDescriptor `json:"idpDescriptor,omitempty"`
	IdpMetadata   *GenericSAMLIdpMetadata   `json:"idpMetadata,omitempty"`
}

func (req *CreateGenericSAMLRequest) ReadFromSchema(d *schema.ResourceData) error {
	return nil
}

type CreateGenericSAMLResponse struct {
	Integration GenericSAMLIntegration `json:"integration"`
}

func (resp *CreateGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return nil
}

func (resp *CreateGenericSAMLResponse) ReadFromSchema(d *schema.ResourceData) error {
	return nil
}

type ReadGenericSAMLResponse struct {
	// TODO
}

func (resp *ReadGenericSAMLResponse) WriteToSchema(d *schema.ResourceData) error {
	return nil
}

func (resp *ReadGenericSAMLResponse) ReadFromSchema(d *schema.ResourceData) error {
	return nil
}

type UpdateGenericSAMLRequest struct {
	// TODO
}

func (resp *UpdateGenericSAMLRequest) ReadFromSchema(d *schema.ResourceData) error {
	return nil
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
		ReadContext:   ReadResource(ReadGenericSAMLConfig()),
		DeleteContext: DeleteResource(DeleteGenericSAMLConfig()),
		Schema: map[string]*schema.Schema{
			"saml_draft_id": {
				Description:  "A valid id for an SAML Draft. Must be at least 5 character long. See resource `cyral_integration_idp_saml_draft`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validationStringLenAtLeast(5),
			},
			"idp_metadata_url": {
				Description:  "A SAML XML IdP Metadata document containing all configuration values required by the Cyral SP.",
				Type:         schema.TypeString,
				Required:     true,
				ExactlyOneOf: idpMetadataTypes,
			},
			"idp_metadata_document": {
				Description:  "Full SAML metadata XML document. Must be base64 encoded.",
				Type:         schema.TypeString,
				Required:     true,
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
