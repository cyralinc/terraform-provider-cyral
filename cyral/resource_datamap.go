package cyral

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DataMap struct {
	SensitiveData SensitiveData `json:"sensitiveData" yaml:"sensitiveData"`
}

type SensitiveData map[string][]RepoAttrs

type RepoAttrs struct {
	Name       string   `json:"repo" yaml:"repo"`
	Attributes []string `json:"attributes" yaml:"attributes"`
}

func resourceDatamap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatamapCreate,
		ReadContext:   resourceDatamapRead,
		UpdateContext: resourceDatamapUpdate,
		DeleteContext: resourceDatamapDelete,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mapping": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"data_location": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"repo": {
										Type:     schema.TypeString,
										Required: true,
									},
									"attributes": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
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

func resourceDatamapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	sensitiveData, err := getSensitiveDataFromResource(d)
	if err != nil {
		return createError("Unable to create datamap", fmt.Sprintf("%v", err))
	}

	sd := sensitiveData.String()
	log.Printf("[DEBUG] resourceDatamapCreate - sensitiveData: %s", sd)

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
	if err := c.UpdateResource(sensitiveData, url); err != nil {
		return createError("Unable to create datamap", fmt.Sprintf("%v", err))
	}

	d.SetId("datamap")

	return resourceDatamapRead(ctx, d, m)
}

func resourceDatamapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	response := make(SensitiveData)
	if err := c.ReadResource(url, &response); err != nil {
		return createError("Unable to read datamap", fmt.Sprintf("%v", err))
	}

	datamap := DataMap{SensitiveData: response}

	sd := datamap.SensitiveData.String()
	log.Printf("[DEBUG] resourceDatamapRead - sensitiveData: %s", sd)

	datamapLabels := flattenSensitiveData(&datamap.SensitiveData)
	log.Printf("[DEBUG] resourceDatamapRead - datamapLabels: %s", datamapLabels)

	if err := d.Set("mapping", datamapLabels); err != nil {
		return createError("Unable to read datamap", fmt.Sprintf("%v", err))
	}

	return diag.Diagnostics{}
}

func resourceDatamapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if d.HasChange("mapping") {
		sensitiveData, err := getSensitiveDataFromResource(d)
		if err != nil {
			return createError("Unable to update datamap", fmt.Sprintf("%v", err))
		}

		sd := sensitiveData.String()
		log.Printf("[DEBUG] resourceDatamapCreate - sensitiveData: %s", sd)

		url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
		if err := c.UpdateResource(sensitiveData, url); err != nil {
			return createError("Unable to update datamap", fmt.Sprintf("%v", err))
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceDatamapRead(ctx, d, m)
}

func resourceDatamapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
	if err := c.DeleteResource(url); err != nil {
		return createError("Unable to delete datamap", fmt.Sprintf("%v", err))
	}

	return diag.Diagnostics{}
}

func getSensitiveDataFromResource(d *schema.ResourceData) (SensitiveData, error) {
	if err := validateMappingBlock(d); err != nil {
		return nil, err
	}

	mappings := d.Get("mapping").(*schema.Set).List()
	sensitiveData := make(SensitiveData)

	for _, m := range mappings {
		labelMap := m.(map[string]interface{})

		labelInfoList := labelMap["data_location"].(*schema.Set).List()
		label := labelMap["label"].(string)

		for _, labelInfo := range labelInfoList {
			labelInfoMap := labelInfo.(map[string]interface{})

			attrs := labelInfoMap["attributes"].([]interface{})
			attributes := []string{}

			for _, attr := range attrs {
				attributes = append(attributes, attr.(string))
			}

			repoAttr := RepoAttrs{
				Name:       labelInfoMap["repo"].(string),
				Attributes: attributes,
			}

			sensitiveData[label] = append(sensitiveData[label], repoAttr)
		}
	}

	return sensitiveData, nil
}

func validateMappingBlock(d *schema.ResourceData) error {
	labelsSet := make(map[string]bool)
	var labels []string
	mappings := d.Get("mapping").(*schema.Set).List()

	for _, m := range mappings {
		labelMap := m.(map[string]interface{})

		label := labelMap["label"].(string)

		if labelsSet[label] {
			labels = append(labels, label)
		} else {
			labelsSet[label] = true
		}
	}

	if len(labels) > 0 {
		return fmt.Errorf("there is more than one mapping block with the same label, please join them into one, labels: %v", labels)
	}

	return nil
}

func flattenSensitiveData(sensitiveData *SensitiveData) []interface{} {
	if sensitiveData != nil {
		labels := make([]interface{}, 0, len(*sensitiveData))

		for label, repoAttrsList := range *sensitiveData {
			labelMap := make(map[string]interface{})

			labelMap["label"] = label

			labelInfoList := make([]interface{}, len(repoAttrsList), len(repoAttrsList))

			for i, repoAttr := range repoAttrsList {
				labelInfoMap := make(map[string]interface{})

				labelInfoMap["repo"] = repoAttr.Name
				labelInfoMap["attributes"] = repoAttr.Attributes

				labelInfoList[i] = labelInfoMap
			}

			labelMap["data_location"] = labelInfoList
			labels = append(labels, labelMap)
		}

		return labels
	}

	return make([]interface{}, 0)
}

func (sensitiveData SensitiveData) String() string {
	var sd string

	for label, repoAttrsList := range sensitiveData {
		s1 := fmt.Sprintf("map[%s]", label)
		sd = sd + s1
		for _, repoAttr := range repoAttrsList {
			s2 := fmt.Sprintf("repo: %s, attributes: [ ", repoAttr.Name)
			sd += s2
			for _, attr := range repoAttr.Attributes {
				s3 := fmt.Sprintf("%s ", attr)
				sd += s3
			}
			sd += "]"
		}
		sd += "\n"
	}

	return sd
}
