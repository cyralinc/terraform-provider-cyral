package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type GetDataLabelResponse DataLabel

func (resp *GetDataLabelResponse) WriteToSchema(d *schema.ResourceData) error {
	if err := WriteDataLabelsToDataSourceSchema([]*DataLabel{(*DataLabel)(resp)}, d); err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	return nil
}

type GetDataLabelsResponse struct {
	Labels []*DataLabel `json:"labels"`
}

func (resp *GetDataLabelsResponse) WriteToSchema(d *schema.ResourceData) error {
	if err := WriteDataLabelsToDataSourceSchema(resp.Labels, d); err != nil {
		return err
	}
	d.SetId(uuid.New().String())
	return nil
}

func WriteDataLabelsToDataSourceSchema(labels []*DataLabel, d *schema.ResourceData) error {
	var labelsList []interface{}
	for _, label := range labels {
		labelsList = append(labelsList, map[string]interface{}{
			"name":        label.Name,
			"description": label.Description,
			"type":        label.Type,
			"tags":        label.TagsAsInterface(),
		})
	}
	if err := d.Set("datalabel_list", labelsList); err != nil {
		return err
	}
	return nil
}

func dataSourceDatalabelReadConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "DatalabelDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			nameFilter := d.Get("name").(string)
			typeFilter := d.Get("type").(string)
			var pathParams string
			if nameFilter != "" {
				pathParams = fmt.Sprintf("/%s", nameFilter)
			}
			queryParams := urlQuery(map[string]string{
				"type": typeFilter,
			})

			return fmt.Sprintf("https://%s/v1/datalabels%s%s", c.ControlPlane, pathParams, queryParams)
		},
		NewResponseData: func(d *schema.ResourceData) ResponseData {
			nameFilter := d.Get("name").(string)
			if nameFilter == "" {
				return &GetDataLabelsResponse{}
			} else {
				return &GetDataLabelResponse{}
			}
		},
	}
}

func dataSourceDatalabel() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter data labels. See also resource [`cyral_datalabel`](../resources/datalabel.md).",
		ReadContext: ReadResource(dataSourceDatalabelReadConfig()),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Retrieve the unique label with this name, if it exists.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  fmt.Sprintf("Filter the results by type of data label. Defaults to `%s`, which will return all label types. List of supported types:", defaultDataLabelType) + supportedTypesMarkdown(dataLabelTypes()),
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultDataLabelType,
				ValidateFunc: validation.StringInSlice(append(dataLabelTypes(), ""), false),
			},
			"datalabel_list": {
				Description: "List of existing data labels satisfying given filter criteria.",
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
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}
