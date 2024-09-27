package datalabel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/datalabel/classificationrule"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var readUpdateDeleteURLFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/datalabels/%s",
		c.ControlPlane,
		d.Get("name").(string))
}

func resourceSchema() *schema.Resource {
	contextHandler := core.DefaultContextHandler{
		ResourceName:                 resourceName,
		ResourceType:                 resourcetype.Resource,
		SchemaReaderFactory:          func() core.SchemaReader { return &DataLabel{} },
		SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &DataLabel{} },
		ReadUpdateDeleteURLFactory:   readUpdateDeleteURLFactory,
	}
	return &schema.Resource{
		Description: "Manages data labels. Data labels are part of the Cyral [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				Type:                operationtype.Create,
				HttpMethod:          http.MethodPut,
				URLFactory:          readUpdateDeleteURLFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &DataLabel{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &DataLabel{} },
			}, readDataLabelConfig,
		),
		ReadContext:   core.ReadResource(readDataLabelConfig),
		UpdateContext: contextHandler.UpdateContext(),
		DeleteContext: contextHandler.DeleteContext(),
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
					"[Automatic Data Map](https://cyral.com/docs/policy/repo-crawler/use-auto-mapping/) feature to automatically map " +
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
	ResourceName:        resourceName,
	Type:                operationtype.Read,
	HttpMethod:          http.MethodGet,
	URLFactory:          readUpdateDeleteURLFactory,
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &DataLabel{} },
	RequestErrorHandler: &core.IgnoreHttpNotFound{ResName: "Data Label"},
}
