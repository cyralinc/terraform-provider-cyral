package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AccessGateway struct {
	SidecarId string `json:"sidecarId,omitempty"`
	BindingId string `json:"bindingId,omitempty"`
}

type GetOrUpdateAccessGateway struct {
	AccessGateway *AccessGateway `json:"accessGateway,omitempty"`
}

func (r *GetOrUpdateAccessGateway) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(RepositoryIDKey).(string))
	d.Set(SidecarIDKey, r.AccessGateway.SidecarId)
	d.Set(BindingIDKey, r.AccessGateway.BindingId)
	return nil
}

func (r *GetOrUpdateAccessGateway) ReadFromSchema(d *schema.ResourceData) error {
	r.AccessGateway = &AccessGateway{
		BindingId: d.Get(BindingIDKey).(string),
		SidecarId: d.Get(SidecarIDKey).(string),
	}
	return nil
}

var ReadRepositoryAccessGatewayConfig = ResourceOperationConfig{
	Name:       "RepositoryAccessGatewayRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos/%s/accessGateway",
			c.ControlPlane,
			d.Get(RepositoryIDKey).(string),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &GetOrUpdateAccessGateway{}
	},
}

func resourceRepositoryAccessGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the sidecar and binding set as the access gateway for [cyral_repositories](./repositories.md).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayCreate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(RepositoryIDKey).(string),
					)
				},
				NewResourceData: func() ResourceData {
					return &GetOrUpdateAccessGateway{}
				},
			},
			ReadRepositoryAccessGatewayConfig,
		),
		ReadContext: ReadResource(ReadRepositoryAccessGatewayConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(RepositoryIDKey).(string),
					)
				},
				NewResourceData: func() ResourceData {
					return &GetOrUpdateAccessGateway{}
				},
			},
			ReadRepositoryAccessGatewayConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(RepositoryIDKey).(string),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			RepositoryIDKey: {
				Description: "ID of the repository the access gateway is associated with. This is also the " +
					"import ID for this resource.",
				Type:     schema.TypeString,
				Required: true,
			},
			SidecarIDKey: {
				Description: "ID of the sidecar that will be set as the access gatway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			BindingIDKey: {
				Description: "ID of the binding that will be set as the access gatway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set(RepositoryIDKey, d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
