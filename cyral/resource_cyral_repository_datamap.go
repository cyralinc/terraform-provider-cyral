package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type DataMapRequest struct {
	DataMap `json:"dataMap,omitempty"`
}

// This is called 'DataMap' and not 'Datamap', because although we consider
// 'datamap' to be a single word in the resource name 'cyral_repository_datamap'
// for ease of writing, 'data map' is actually two words in English.
type DataMap struct {
	Labels map[string]*DataMapMapping `json:"labels,omitempty"`
}

func (dm *DataMap) WriteToSchema(d *schema.ResourceData) error {
	var mappings []interface{}
	for label, mapping := range dm.Labels {
		mappingContents := make(map[string]interface{})

		var attributes []string
		if mapping != nil {
			attributes = mapping.Attributes
		}

		mappingContents["label"] = label
		mappingContents["attributes"] = attributes

		mappings = append(mappings, mappingContents)
	}

	return d.Set("mapping", mappings)
}

func (dm *DataMap) equal(other DataMap) bool {
	for label, thisMapping := range dm.Labels {
		if otherMapping, ok := other.Labels[label]; ok {
			if !elementsMatch(thisMapping.Attributes, otherMapping.Attributes) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

type DataMapMapping struct {
	Attributes []string `json:"attributes,omitempty"`
}

func resourceRepositoryDatamap() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: resourceRepositoryDatamapCreate,
		ReadContext:   resourceRepositoryDatamapRead,
		UpdateContext: resourceRepositoryDatamapUpdate,
		DeleteContext: resourceRepositoryDatamapDelete,
		Schema: map[string]*schema.Schema{
			"repository_id": {
				Description: "ID of the repository for which to configure a data map.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"mapping": {
				Description: "Mapping of a label to a list of data locations (attributes).",
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Description: "Label given to the attributes in this mapping.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"attributes": {
							Description: "List containing the specific locations of the data within the repo, " +
								"following the pattern `{SCHEMA}.{TABLE}.{ATTRIBUTE}` (ex: " +
								"`[your_schema_name.your_table_name.your_attr_name]`).\n\n" +
								"> Note: When referencing data in Dremio repository, please include the complete " +
								"location in `attributes`, separating spaces by dots. For example, an attribute " +
								"`my_attr` from table `my_tbl` within space `inner_space` within space `outer_space` " +
								"would be referenced as `outer_space.inner_space.my_tbl.my_attr`. For more information, " +
								"please see the [Policy Guide](https://cyral.com/docs/reference/policy/).",
							Type:     schema.TypeList,
							Required: true,
							// TODO: this ForceNew propagates to the parent attribute `mapping`. Therefore, any
							// new mapping will force recreation. In the future, it would be good to use the
							// `v1/repos/{repoID}/datamap/labels/{label}/attributes/{attribute}` endpoint to
							// avoid unnecessary resource recreation. -aholmquist 2022-08-04
							ForceNew: true,
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

func resourceRepositoryDatamapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDatamapCreate")
	c := m.(*client.Client)

	repoID := d.Get("repository_id").(string)
	url := fmt.Sprintf("https://%s/v1/repos/%s/datamap", c.ControlPlane, repoID)

	dataMap := getDatamapFromResource(d)
	dataMapRequest := DataMapRequest{DataMap: dataMap}

	_, err := c.DoRequest(url, http.MethodPut, dataMapRequest)
	if err != nil {
		return createError("Unable to create repository datamap", err.Error())
	}

	d.SetId(repoID)
	// Write data map here to avoid issues with the order of the attributes
	if err := dataMap.WriteToSchema(d); err != nil {
		return createError("Unable to create repository datamap", err.Error())
	}

	log.Printf("[DEBUG] End resourceRepositoryDatamapCreate")

	return resourceRepositoryDatamapRead(ctx, d, m)
}

func resourceRepositoryDatamapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDatamapRead")
	c := m.(*client.Client)

	repoID := d.Id()
	url := fmt.Sprintf("https://%s/v1/repos/%s/datamap", c.ControlPlane, repoID)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to create repository datamap", err.Error())
	}

	dataMap := DataMap{}
	if err := json.Unmarshal(body, &dataMap); err != nil {
		return createError("Unable to unmarshall JSON", err.Error())
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", dataMap)

	currentDataMap := getDatamapFromResource(d)
	if !currentDataMap.equal(dataMap) {
		if err := dataMap.WriteToSchema(d); err != nil {
			return createError("Unable to read repository datamap", err.Error())
		}
	}

	log.Printf("[DEBUG] End resourceRepositoryDatamapRead")

	return diag.Diagnostics{}
}

func resourceRepositoryDatamapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDatamapUpdate")
	c := m.(*client.Client)

	repoID := d.Id()
	url := fmt.Sprintf("https://%s/v1/repos/%s/datamap", c.ControlPlane, repoID)

	dataMap := getDatamapFromResource(d)
	dataMapRequest := DataMapRequest{DataMap: dataMap}

	_, err := c.DoRequest(url, http.MethodPut, dataMapRequest)
	if err != nil {
		return createError("Unable to create repository datamap", err.Error())
	}

	log.Printf("[DEBUG] End resourceRepositoryDatamapUpdate")

	return resourceRepositoryDatamapRead(ctx, d, m)
}

func resourceRepositoryDatamapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDatamapDelete")
	c := m.(*client.Client)

	repoID := d.Id()
	url := fmt.Sprintf("https://%s/v1/repos/%s/datamap", c.ControlPlane, repoID)

	_, err := c.DoRequest(url, http.MethodDelete, nil)
	if err != nil {
		return createError("Unable to delete repository datamap", err.Error())
	}

	log.Printf("[DEBUG] End resourceRepositoryDatamapDelete")

	return diag.Diagnostics{}
}

func getDatamapFromResource(d *schema.ResourceData) DataMap {
	mappings := d.Get("mapping").(*schema.Set).List()

	dataMap := DataMap{
		Labels: make(map[string]*DataMapMapping),
	}
	for _, mappingIface := range mappings {
		mapping := mappingIface.(map[string]interface{})

		label := mapping["label"].(string)
		var attributes []string
		if mappingAtts, ok := mapping["attributes"]; ok {
			for _, attributeIface := range mappingAtts.([]interface{}) {
				attributes = append(attributes, attributeIface.(string))
			}
		}
		dataMap.Labels[label] = &DataMapMapping{
			Attributes: attributes,
		}
	}

	return dataMap
}
