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
	DisableIdPInitiatedLogin bool   `json:"disableIdPInitiatedLogin"`
	IdpType                  string `json:"idpType,omitempty"`
	*RequiredUserAttributes  `json:"attributes,omitempty"`
}

func (req *CreateGenericSAMLDraftRequest) ReadFromSchema(d *schema.ResourceData) error {
	req.DisplayName = d.Get("display_name").(string)
	req.DisableIdPInitiatedLogin = d.Get("disable_idp_initiated_login").(bool)
	req.IdpType = d.Get("idp_type").(string)

	attributes, err := RequiredUserAttributesFromSchema(d)
	if err != nil {
		return err
	}
	req.RequiredUserAttributes = attributes

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
	if err := d.Set("idp_type", resp.Draft.IdpType); err != nil {
		return err
	}
	if resp.Draft.RequiredUserAttributes != nil {
		if err := resp.Draft.RequiredUserAttributes.WriteToSchema(d); err != nil {
			return err
		}
	}
	return nil
}

func (resp *GenericSAMLDraftResponse) ReadFromSchema(d *schema.ResourceData) error {
	resp.Draft.ID = d.Id()
	resp.Draft.DisplayName = d.Get("display_name").(string)
	resp.Draft.IdpType = d.Get("idp_type").(string)
	resp.Draft.DisableIdPInitiatedLogin = d.Get("disable_idp_initiated_login").(bool)
	resp.Draft.SPMetadata = &SPMetadata{XMLDocument: d.Get("sp_metadata").(string)}
	attributes, err := RequiredUserAttributesFromSchema(d)
	if err != nil {
		return err
	}
	resp.Draft.RequiredUserAttributes = attributes

	return nil
}

type GenericSAMLDraft struct {
	ID                       string `json:"id"`
	DisplayName              string `json:"displayName"`
	IdpType                  string `json:"idpType"`
	DisableIdPInitiatedLogin bool   `json:"disableIdpInitiatedLogin"`
	*SPMetadata              `json:"spMetadata,omitempty"`
	*RequiredUserAttributes  `json:"requiredUserAttributes,omitempty"`
	Completed                bool `json:"completed"`
}

type SPMetadata struct {
	XMLDocument string `json:"xmlDocument"`
}

type RequiredUserAttributes struct {
	FirstName UserAttribute `json:"firstName,omitempty"`
	LastName  UserAttribute `json:"lastName,omitempty"`
	Email     UserAttribute `json:"email,omitempty"`
	Groups    UserAttribute `json:"groups,omitempty"`
}

func NewRequiredUserAttributes(firstName, lastName, email, groups string) *RequiredUserAttributes {
	return &RequiredUserAttributes{
		FirstName: UserAttribute{
			Name: firstName,
		},
		LastName: UserAttribute{
			Name: lastName,
		},
		Email: UserAttribute{
			Name: email,
		},
		Groups: UserAttribute{
			Name: groups,
		},
	}
}

func (uatt *RequiredUserAttributes) WriteToSchema(d *schema.ResourceData) error {
	var attributes []interface{}
	attributes = append(attributes, map[string]interface{}{
		"first_name": uatt.FirstName.Name,
		"last_name":  uatt.LastName.Name,
		"email":      uatt.Email.Name,
		"groups":     uatt.Groups.Name,
	})
	return d.Set("attributes", attributes)
}

func RequiredUserAttributesFromSchema(d *schema.ResourceData) (*RequiredUserAttributes, error) {
	attributesSet := d.Get("attributes").(*schema.Set).List()
	if len(attributesSet) > 1 {
		return nil, fmt.Errorf("Expected 'attributes' to be a set with at "+
			"most one element, got %d elements", len(attributesSet))
	} else if len(attributesSet) > 0 {
		attributesMap := attributesSet[0].(map[string]interface{})
		return NewRequiredUserAttributes(
			attributesMap["first_name"].(string),
			attributesMap["last_name"].(string),
			attributesMap["email"].(string),
			attributesMap["groups"].(string),
		), nil
	}
	return nil, nil
}

type UserAttribute struct {
	Name string `json:"name"`
}

func CreateGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceCreate",
		HttpMethod: http.MethodPost,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts", c.ControlPlane)
		},
		NewResourceData: func() ResourceData { return &CreateGenericSAMLDraftRequest{} },
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &GenericSAMLDraftResponse{} },
	}
}

func ReadGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &GenericSAMLDraftResponse{} },
	}
}

func DeleteGenericSAMLConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "GenericSAMLResourceDelete",
		HttpMethod: http.MethodDelete,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf("https://%s/v1/integrations/generic-saml/drafts/%s", c.ControlPlane, d.Id())
		},
	}
}

func resourceIntegrationIdPGenericSAML() *schema.Resource {
	return &schema.Resource{
		Description: "Manages generic SAML integration drafts.",
		CreateContext: CreateResource(
			CreateGenericSAMLConfig(),
			ReadGenericSAMLConfig(),
		),
		ReadContext:   ReadResource(ReadGenericSAMLConfig()),
		DeleteContext: DeleteResource(DeleteGenericSAMLConfig()),
		Schema: map[string]*schema.Schema{
			// All of the input arguments must force recreation of
			// the resource, because the API does not support
			// updates. If you try to use the Create API, it will do
			// the equivalent of deleting + recreating, generating a
			// new ID every time. In this case, the terraform
			// resource ID would become inconsistent with the Cyral
			// ID.
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
			"idp_type": {
				Description: "Identity provider type. The value provided can be used as a filter when retrieving SAML drafts. See data source `cyral_integration_idp_generic_saml`.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"attributes": {
				Description: "SAML Attribute names for the identity attributes required by the Cyral SP. Each attribute name MUST be at least 3 characters long.",
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				// Fix min items to 1 because it doesn't make
				// sense to have an empty set.
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users first name (given name).",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "firstName",
							ValidateFunc: validationStringLenAtLeast(3),
						},
						"last_name": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users last name (family name).",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "lastName",
							ValidateFunc: validationStringLenAtLeast(3),
						},
						"email": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users email address.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "email",
							ValidateFunc: validationStringLenAtLeast(3),
						},
						"groups": {
							Description:  "The name of the attribute in the incoming SAML assertion containing the users group membership in the IdP.",
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "memberOf",
							ValidateFunc: validationStringLenAtLeast(3),
						},
					},
				},
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
