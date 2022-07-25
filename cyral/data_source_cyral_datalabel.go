package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const (
	defaultDataLabelType = "UNKNOWN"
)

func dataLabelTypes() []string {
	return []string{
		"UNKNOWN",
		"PREDEFINED",
		"CUSTOM",
	}
}

type GetDataLabelsResponse struct {
	Labels []DataLabel `json:"labels"`
}

func (resp *GetDataLabelsResponse) WriteToSchema(d *schema.ResourceData) error {
	var labels []interface{}
	for _, label := range resp.Labels {
		var tags []interface{}
		for _, tag := range label.Tags {
			tags = append(tags, tag)
		}
		labels = append(labels, map[string]interface{}{
			"description": label.Description,
			"type":        label.Type,
			"tags":        tags,
		})
		labels = append(labels, label)
	}

	if err := d.Set("datalabel_list", labels); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

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

			return fmt.Sprintf("https://%s/v1/repos%s%s", c.ControlPlane, pathParams, queryParams)
		},
		NewResponseData: func() ResponseData { return &GetReposResponse{} },
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
				ValidateFunc: validation.StringInSlice(dataLabelTypes(), false),
			},
			"datalabel_list": {
				Description: "List of existing data labels satisfying given filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
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
							Description: "Tags used to categorize data labels.",
							Type:        schema.TypeList,
							Optional:    true,
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
