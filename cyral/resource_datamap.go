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

type SensitiveData map[string][]*RepoAttrs

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
			"labels": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"label_info": {
							Type:     schema.TypeList,
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
	}
}

func resourceDatamapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	sensitiveData := getSensitiveDataFromResource(d)

	sd := makeStrForSD(sensitiveData)
	log.Printf("[DEBUG] resourceDatamapCreate --- sensitiveData: %s", sd)

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
	if err := c.UpdateResource(sensitiveData, url); err != nil {
		return createError("Unable to create datamap", fmt.Sprintf("%v", err))
	}

	// if err := c.UpsertDatamap(sensitiveData); err != nil {
	// 	return diag.FromErr(err)
	// }

	d.SetId(time.Now().Format(time.RFC850))

	return resourceDatamapRead(ctx, d, m)
}

func flattenDatamapLabels(datamapLabels *SensitiveData) []interface{} {
	if datamapLabels != nil {
		dmLabels := make([]interface{}, len(*datamapLabels), len(*datamapLabels))

		for label, repoAttrsList := range *datamapLabels {
			dmLabel := make(map[string]interface{})

			dmLabel["label_id"] = label

			repoAttrs := make([]interface{}, len(repoAttrsList), len(repoAttrsList))

			for i, repoAttr := range repoAttrsList {
				li := make(map[string]interface{})

				li["repo"] = repoAttr.Name
				li["attributes"] = repoAttr.Attributes

				repoAttrs[i] = li
			}

			dmLabel["label_info"] = repoAttrs
			dmLabels = append(dmLabels, dmLabel)
		}

		return dmLabels
	}

	return make([]interface{}, 0)
}

func resourceDatamapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// datamap, err := c.GetDatamap()
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

	response := make(SensitiveData)
	if err := c.ReadResource(d.Id(), url, &response); err != nil {
		return createError("Unable to read datamap", fmt.Sprintf("%v", err))
	}

	datamap := DataMap{SensitiveData: response}

	sd := makeStrForSD(datamap.SensitiveData)
	log.Printf("[DEBUG] resourceDatamapRead --- sensitiveData: %s", sd)

	datamapLabels := flattenDatamapLabels(&datamap.SensitiveData)
	log.Printf("[DEBUG] resourceDatamapRead --- datamapLabels: %s", datamapLabels)

	if err := d.Set("labels", datamapLabels); err != nil {
		log.Printf("[DEBUG] resourceDatamapRead --- ERRO d.Set(\"labels\", datamapLabels): %v", err)
		return diag.FromErr(err)
	}

	log.Print("[DEBUG] resourceDatamapRead --- NORMAL")

	return diags
}

func resourceDatamapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if d.HasChange("labels") {
		sensitiveData := getSensitiveDataFromResource(d)

		sd := makeStrForSD(sensitiveData)
		log.Printf("[DEBUG] resourceDatamapCreate --- sensitiveData: %s", sd)

		url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
		if err := c.UpdateResource(sensitiveData, url); err != nil {
			return createError("Unable to update datamap", fmt.Sprintf("%v", err))
		}
		// if err := c.UpsertDatamap(sensitiveData); err != nil {
		// 	return diag.FromErr(err)
		// }

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceDatamapRead(ctx, d, m)
}

func resourceDatamapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)
	if err := c.DeleteResource(url); err != nil {
		return createError("Unable to delete datamap", fmt.Sprintf("%v", err))
	}

	// err := c.DeleteDatamap()
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags

}

func getSensitiveDataFromResource(d *schema.ResourceData) SensitiveData {
	labels := d.Get("labels").([]interface{})
	sensitiveData := make(SensitiveData)

	for _, label := range labels {
		labelMap := label.(map[string]interface{})

		labelInfoList := labelMap["label_info"].([]interface{})
		repoAttrsList := []*RepoAttrs{}

		for _, labelInfo := range labelInfoList {
			labelInfoMap := labelInfo.(map[string]interface{})

			attrs := labelInfoMap["attributes"].([]interface{})
			attributes := []string{}

			for _, attr := range attrs {
				attributes = append(attributes, attr.(string))
			}

			li := RepoAttrs{
				Name:       labelInfoMap["repo"].(string),
				Attributes: attributes,
			}

			repoAttrsList = append(repoAttrsList, &li)
		}

		sensitiveData[labelMap["label_id"].(string)] = repoAttrsList
	}

	return sensitiveData
}

// // upsertDatamap creates or updates a datamap
// func upsertDatamap(sensitiveData SensitiveData) error {
// 	payloadBytes, err := json.Marshal(sensitiveData)

// 	if err != nil {
// 		return fmt.Errorf("failed to encode 'create/update datamap' payload; err: %v", err)
// 	}

// 	url := fmt.Sprintf("https://%s/v1/datamaps", c.ControlPlane)

// 	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
// 	if err != nil {
// 		return fmt.Errorf("unable to create 'create/update datamap' request; err: %v", err)
// 	}

// 	if _, err := c.doRequest(req); err != nil {
// 		return err
// 	}

// 	return nil
// }

func makeStrForSD(sensitiveData SensitiveData) string {
	var sd string

	for key, value := range sensitiveData {
		s1 := fmt.Sprintf("map[%s]", key)
		sd = sd + s1
		for _, r := range value {
			s2 := fmt.Sprintf("repo: %s, attributes: [ ", r.Name)
			sd += s2
			for _, attr := range r.Attributes {
				s3 := fmt.Sprintf("%s ", attr)
				sd += s3
			}
			sd += "]"
		}
	}

	return sd
}
