package cyral

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/model"
)

func resourceDatalabel() *schema.Resource {
	return &schema.Resource{
		Description: "Manages data labels. Data labels are part of the Cyral [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "DataLabelResourceCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/datalabels/%s",
						c.ControlPlane,
						d.Get("name").(string))
				},
				NewResourceData: func() ResourceData { return &model.DataLabel{} },
				NewResponseData: func(_ *schema.ResourceData) ResponseData { return &CreateListenerAPIResponse{} },
			}, ReadDataLabelConfig,
		),
		ReadContext: ReadResource(ReadDataLabelConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "DataLabelResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/datalabels/%s",
						c.ControlPlane,
						d.Get("name").(string))
				},
				NewResourceData: func() ResourceData { return &SidecarListenerResource{} },
			}, ReadSidecarListenersConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "DataLabelResourceDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/datalabels/%s",
						c.ControlPlane,
						d.Get("name").(string))
				},
			},
		),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name of the data label.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description of the data label.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "Tags that can be used to categorize data labels.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"classification_rule": {
				Description: "Classification rules are used by the [Automatic Data Map](https://cyral.com/docs/policy/automatic-datamap) feature to automatically map data locations to labels. Currently, only `PREDEFINED` labels have classification rules.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Description: "Type of the classification rule.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rule_code": {
							Description: "Actual code of the classification rule. For example, this attribute may contain REGO code for `REGO`-type classification rules.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rule_status": {
							Description: "Status of the classification rule.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m interface{},
			) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

var ReadDataLabelConfig = ResourceOperationConfig{
	Name:       "DataLabelResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/datalabels/%s",
			c.ControlPlane,
			d.Get("name").(string))
	},
	NewResponseData:     func(_ *schema.ResourceData) ResponseData { return &GetDataLabelResponse{} },
	RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "Data Label"},
}
