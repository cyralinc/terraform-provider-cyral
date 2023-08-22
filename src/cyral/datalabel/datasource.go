package datalabel

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/src/client"
	"github.com/cyralinc/terraform-provider-cyral/src/core"
	"github.com/cyralinc/terraform-provider-cyral/src/cyral"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
)

func init() {
	sr := &core.SchemaRegister{
		Name:   "cyral_datalabel",
		Schema: DataSourceSchema,
		Type:   core.DataSourceSchema,
	}
	cyral.RegisterToProvider(sr)
}

type GetDataLabelResponse DataLabel

func (resp *GetDataLabelResponse) WriteToSchema(d *schema.ResourceData) error {
	if err := writeDataLabelsToDataSourceSchema([]*DataLabel{(*DataLabel)(resp)}, d); err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	return nil
}

type GetDataLabelsResponse struct {
	Labels []*DataLabel `json:"labels"`
}

func (dl *GetDataLabelsResponse) WriteToSchema(d *schema.ResourceData) error {
	if err := writeDataLabelsToDataSourceSchema(dl.Labels, d); err != nil {
		return err
	}

	d.SetId(uuid.New().String())
	return nil
}

func writeDataLabelsToDataSourceSchema(labels []*DataLabel, d *schema.ResourceData) error {
	var labelsList []interface{}
	for _, label := range labels {
		labelsList = append(labelsList, map[string]interface{}{
			"name":                label.Name,
			"description":         label.Description,
			"type":                label.Type,
			"tags":                label.Tags.AsInterface(),
			"classification_rule": label.ClassificationRule.AsInterface(),
			"implicit":            label.Implicit,
		})
	}
	if err := d.Set("datalabel_list", labelsList); err != nil {
		return err
	}
	return nil
}

func readConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		Name:       "DatalabelDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			nameFilter := d.Get("name").(string)
			typeFilter := d.Get("type").(string)
			var pathParams string
			if nameFilter != "" {
				pathParams = fmt.Sprintf("/%s", nameFilter)
			}
			queryParams := utils.UrlQuery(map[string]string{
				"type": typeFilter,
			})

			return fmt.Sprintf("https://%s/v1/datalabels%s%s", c.ControlPlane, pathParams, queryParams)
		},
		NewResponseData: func(d *schema.ResourceData) core.ResponseData {
			nameFilter := d.Get("name").(string)
			if nameFilter == "" {
				return &GetDataLabelsResponse{}
			} else {
				return &GetDataLabelResponse{}
			}
		},
	}
}

func DataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter data labels. See also resource [`cyral_datalabel`](../resources/datalabel.md).",
		ReadContext: core.ReadResource(readConfig()),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Retrieve the unique label with this name, if it exists.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  fmt.Sprintf("Filter the results by type of data label. Defaults to `%s`, which will return all label types. The labels you create will always have type `CUSTOM`. Labels that come pre-configured in the control plane have type `PREDEFINED`. List of supported types:", Default) + utils.SupportedTypesMarkdown(TypesAsString()),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      Default,
				ValidateFunc: validation.StringInSlice(append(TypesAsString(), ""), false),
			},
			"datalabel_list": {
				Description: "List of existing data labels satisfying the filter criteria.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the data label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Type of the data label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Description of the data label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"tags": {
							Description: "Tags used to categorize data labels.",
							Type:        schema.TypeList,
							Computed:    true,
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
						"implicit": {
							Description: "If true, the label only exists implicitly in the legacy data map API (i.e. `v1/datamaps`). Implicit labels always have `CUSTOM` type.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
