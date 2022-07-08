package cyral

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type CreateGenericSAMLDraftRequest struct {
	DisplayName              string `json:"displayName"`
	DisableIdPInitiatedLogin bool   `json:"disableIdpInitiatedLogin"`
}

func (req *CreateGenericSAMLDraftRequest) ReadFromSchema(d *schema.ResourceData) error {
	req.DisplayName = d.Get("display_name").(string)
	req.DisableIdPInitiatedLogin = d.Get("disable_idp_initiated_login").(bool)
	return nil
}

type GenericSAMLDraftResponse struct {
	Draft GenericSAMLDraft `json:"draft"`
}

func (resp *GenericSAMLDraftResponse) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(resp.Draft.ID)
	if err := d.Set("display_name", resp.Draft.DisplayName); err != nil {
		return err
	}
	if err := d.Set("disable_idp_initiated_login", resp.Draft.DisableIdPInitiatedLogin); err != nil {
		return err
	}
	if err := d.Set("sp_metadata", resp.Draft.SPMetadata.XMLDocument); err != nil {
		return err
	}
	return nil
}

func (resp *GenericSAMLDraftResponse) ReadFromSchema(d *schema.ResourceData) error {
	resp.Draft.ID = d.Id()
	resp.Draft.DisplayName = d.Get("display_name").(string)
	resp.Draft.DisableIdPInitiatedLogin = d.Get("disable_idp_initiated_login").(bool)
	resp.Draft.SPMetadata = SPMetadata{XMLDocument: d.Get("sp_metadata").(string)}
	return nil
}

type GenericSAMLDraft struct {
	ID                       string `json:"id"`
	DisplayName              string `json:"displayName"`
	IdPType                  string `json:"idpType"`
	DisableIdPInitiatedLogin bool   `json:"disableIdpInitiatedLogin"`
	SPMetadata               `json:"spMetadata"`
	RequiredUserAttributes   `json:"requiredUserAttributes"`
	Completed                bool `json:"completed"`
}

type SPMetadata struct {
	XMLDocument string `json:"xmlDocument"`
}

type RequiredUserAttributes struct {
	FirstName UserAttribute `json:"firstName"`
	LastName  UserAttribute `json:"lastName"`
	Email     UserAttribute `json:"email"`
	Groups    UserAttribute `json:"groups"`
}

type UserAttribute struct {
	Name string `json:"name"`
}

func CreateSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts", c.ControlPlane)
		},
		NewResourceData: func() ResourceData { return &CreateGenericSAMLDraftRequest{} },
		NewResponseData: func() ResponseData { return &GenericSAMLDraftResponse{} },
	}
}

func ReadSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
		NewResponseData: func() ResponseData { return &GenericSAMLDraftResponse{} },
	}
}

func DeleteSAMLDraftConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SAMLDraftResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
	}
}

func resourceIntegrationIdPSAMLDraft() *schema.Resource {
	return &schema.Resource{
		Description: "Manages generic SAML integration drafts.",
		CreateContext: CreateResource(
			CreateSAMLDraftConfig(),
			ReadSAMLDraftConfig(),
		),
		ReadContext:   ReadResource(ReadSAMLDraftConfig()),
		DeleteContext: DeleteResource(DeleteSAMLDraftConfig()),
		Schema: map[string]*schema.Schema{
			// All of the input arguments must force recreation of
			// the resource, because the API does not support
			// updates.
			"display_name": {
				Description:  "Display name used in the Cyral control plane.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
			},
			"disable_idp_initiated_login": {
				Description: "Whether or not IdP-Initiated login should be disabled for this generic SAML integration draft. Defaults to false.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
			},
			"sp_metadata": {
				Description: "The SP Metadata document describing the Cyral service provider for this integration.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"id": {
				Description: "ID of this resource in the Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
