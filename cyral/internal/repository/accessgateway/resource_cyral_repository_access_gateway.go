package accessgateway

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AGData struct {
	SidecarId string `json:"sidecarId,omitempty"`
	BindingId string `json:"bindingId,omitempty"`
}

type AccessGateway struct {
	AGData *AGData `json:"accessGateway,omitempty"`
}

func (r *AccessGateway) WriteToSchema(d *schema.ResourceData) error {
	d.SetId(d.Get(utils.RepositoryIDKey).(string))
	d.Set(utils.SidecarIDKey, r.AGData.SidecarId)
	d.Set(utils.BindingIDKey, r.AGData.BindingId)
	return nil
}

func (r *AccessGateway) ReadFromSchema(d *schema.ResourceData) error {
	r.AGData = &AGData{
		BindingId: d.Get(utils.BindingIDKey).(string),
		SidecarId: d.Get(utils.SidecarIDKey).(string),
	}
	return nil
}

var ReadRepositoryAccessGatewayConfig = core.ResourceOperationConfig{
	Name:       "RepositoryAccessGatewayRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos/%s/accessGateway",
			c.ControlPlane,
			d.Get(utils.RepositoryIDKey).(string),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
		return &AccessGateway{}
	},
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Repository access gateway"},
}

func ResourceRepositoryAccessGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Manages the sidecar and binding set as the access gateway for [cyral_repositories](./repositories.md).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayCreate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(utils.RepositoryIDKey).(string),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &AccessGateway{}
				},
			},
			ReadRepositoryAccessGatewayConfig,
		),
		ReadContext: core.ReadResource(ReadRepositoryAccessGatewayConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(utils.RepositoryIDKey).(string),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &AccessGateway{}
				},
			},
			ReadRepositoryAccessGatewayConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "RepositoryAccessGatewayDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s/accessGateway",
						c.ControlPlane,
						d.Get(utils.RepositoryIDKey).(string),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			utils.RepositoryIDKey: {
				Description: "ID of the repository the access gateway is associated with. This is also the " +
					"import ID for this resource.",
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			utils.SidecarIDKey: {
				Description: "ID of the sidecar that will be set as the access gateway for the given repository.",
				Type:        schema.TypeString,
				Required:    true,
			},
			utils.BindingIDKey: {
				Description: "ID of the binding that will be set as the access gateway for the given repository.  " +
					"Note that modifications to this field will result in terraform replacing the given " +
					"access gateway resource, since the access gateway must be deleted before binding. ",
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set(utils.RepositoryIDKey, d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
