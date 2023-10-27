package datalabel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel/classificationrule"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages data labels. Data labels are part of the Cyral [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "DataLabelResourceCreate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/datalabels/%s",
						c.ControlPlane,
						d.Get("name").(string))
				},
				NewResourceData: func() core.ResourceData { return &DataLabel{} },
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData { return &DataLabel{} },
			}, readDataLabelConfig,
		),
		ReadContext: core.ReadResource(readDataLabelConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "DataLabelResourceUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/datalabels/%s",
						c.ControlPlane,
						d.Get("name").(string))
				},
				NewResourceData: func() core.ResourceData { return &DataLabel{} },
			}, readDataLabelConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
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
				Description: "Classification rules are used by the " +
					"[Automatic Data Map](https://cyral.com/docs/policy/automatic-datamap) feature to automatically map " +
					"data locations to labels.",
				Optional: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Description: "Type of the classification rule. List of supported values: " +
								utils.SupportedValuesAsMarkdown(classificationrule.TypesAsString()),
							Type:         schema.TypeString,
							Optional:     true,
							Default:      classificationrule.Unknown,
							ValidateFunc: validation.StringInSlice(classificationrule.TypesAsString(), false),
						},
						"rule_code": {
							Description: "Actual code of the classification rule. For example, this attribute may contain " +
								"REGO code for `REGO`-type classification rules.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_status": {
							Description: "Status of the classification rule. List of supported values: " +
								utils.SupportedValuesAsMarkdown(classificationrule.StatusesAsString()),
							Type:         schema.TypeString,
							Optional:     true,
							Default:      classificationrule.Enabled,
							ValidateFunc: validation.StringInSlice(classificationrule.StatusesAsString(), false),
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

var readDataLabelConfig = core.ResourceOperationConfig{
	Name:       "DataLabelResourceRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/datalabels/%s",
			c.ControlPlane,
			d.Get("name").(string))
	},
	NewResponseData:     func(_ *schema.ResourceData) core.ResponseData { return &DataLabel{} },
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Data Label"},
}
