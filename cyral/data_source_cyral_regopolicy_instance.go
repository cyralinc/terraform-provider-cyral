package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type GetRegopolicyInstancesResponse struct {
	Instances []*PolicyInstance `json:"instances"`
}

func (resp *GetRegopolicyInstancesResponse) WriteToSchema(d *schema.ResourceData) error {
	if err := writeRegopolicyInstancesToDataSourceSchema(resp.Instances, d); err != nil {
		return err
	}
	// creating a new id for data source
	d.SetId(uuid.New().String())
	return nil
}

func writeRegopolicyInstancesToDataSourceSchema(instances []*PolicyInstance, d *schema.ResourceData) error {
	var instancesList []interface{}
	for _, instance := range instances {
		instancesList = append(instancesList, map[string]interface{}{
			"name":        instance.Name,
			"description": instance.Description,
			"template_id": instance.TemplateId,
			"tags":        instance.TagsAsInterface(),
		})
	}
	if err := d.Set("regopolicy_instance_list", instancesList); err != nil {
		return err
	}
	return nil
}

func dataSourceRegopolicyInstanceReadConfig() ResourceOperationConfig {
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

func dataSourceRegopolicyInstance() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter policy instances. See also resource [`cyral_regopolicy`](../resources/regopolicy_instance.md).",
		ReadContext: ReadResource(dataSourceDatalabelReadConfig()),
		Schema: map[string]*schema.Schema{
			"regopolicy_instance_list": {
				Description: "List of existing regopolicy instances satisfying given filter criteria.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Description: "Name of the policy instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": {
							Description: "Description for the policy instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"template_id": {
							Description: "Template Id on which the instance was based",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tags": {
							Description: "Tags used to categorize policy instance.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
		},
	}
}
