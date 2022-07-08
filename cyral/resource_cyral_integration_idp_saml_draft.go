package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type GenericSAMLDraft struct {
	DisplayName              string `json:"displayName"`
	DisableIdPInitiatedLogin string `json:"disableIdPInitiatedLogin,omitempty"`
}

func (data GenericSAMLDraft) WriteToSchema(d *schema.ResourceData) error {
	// TODO
	return nil
}

func (data *GenericSAMLDraft) ReadFromSchema(d *schema.ResourceData) error {
	// TODO
	return nil
}

func CreateSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/generic-saml/drafts", c.ControlPlane)
		},
		NewResponseData: func() ResponseData { return &GenericSAMLDraft{} },
	}
}

func ReadSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
		NewResponseData: func() ResponseData { return &GenericSAMLDraft{} },
	}
}

// func UpdateSAMLDraftConfig() ResourceOperationConfig {
// 	return ResourceOperationConfig{
// 		Name:       "SAMLDraftResourceUpdate",
// 		HttpMethod: http.MethodPost,
// 		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
// 			return fmt.Sprintf("https://%s/v1/generic-saml/drafts/%s", c.ControlPlane, d.Id())
// 		},
// 		NewResponseData: func() ResponseData { return &GenericSAMLDraft{} },
// 	}
// }

func DeleteSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
	}
}

func resourceIntegrationSAMLDraft() *schema.Resource {
	return &schema.Resource{
		Description: "Manages generic SAML integration drafts.",
		CreateContext: CreateResource(
			CreateSAMLDraftConfig(),
			ReadSAMLDraftConfig(),
		),
		ReadContext: ReadResource(ReadSAMLDraftConfig()),
		//		UpdateContext: UpdateResource(UpdateSAMLDraftConfig()),
		DeleteContext: DeleteResource(DeleteSAMLDraftConfig()),
		Schema: map[string]*schema.Schema{
			// All of the input arguments force a recreation of the
			// draft, because the API does not support updates.
			"display_name": {
				Description: "Display name used in the Cyral control plane.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"disable_idp_initiated_login": {
				Description: "Whether or not IdP-Initiated login should be disabled for this generic SAML integration draft. Defaults to false.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"saml_metadata_document": {
				Description: "The SP Metadata document describing the Cyral service provider for this integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
