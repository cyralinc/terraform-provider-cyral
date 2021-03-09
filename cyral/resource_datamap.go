package cyral

import (
	"context"
	"log"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

	labels := d.Get("labels").([]interface{})
	sensitiveData := client.SensitiveData{}

	for _, label := range labels {
		labelMap := label.(map[string]interface{})

		labelInfoList := labelMap["label_info"].([]interface{})
		repoAttrsList := []*client.RepoAttrs{}

		for _, labelInfo := range labelInfoList {
			labelInfoMap := labelInfo.(map[string]interface{})

			attrs := labelInfoMap["attributes"].([]interface{})
			attributes := []string{}

			for _, attr := range attrs {
				attributes = append(attributes, attr.(string))
			}

			li := client.RepoAttrs{
				Name:       labelInfoMap["repo"].(string),
				Attributes: attributes,
			}

			repoAttrsList = append(repoAttrsList, &li)
		}

		sensitiveData[labelMap["label_id"].(string)] = repoAttrsList
	}

	sd := client.MakeStrForSD(sensitiveData)
	log.Printf("[DEBUG] resourceDatamapCreate --- sensitiveData: %s", sd)

	if err := c.CreateDatamap(sensitiveData); err != nil {
		return diag.FromErr(err)
	}

	// d.SetId(strconv.Itoa(datamap.ID))

	return resourceDatamapRead(ctx, d, m)
}

func flattenDatamapLabels(datamapLabels *client.SensitiveData) []interface{} {
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

	datamap, err := c.GetDatamap()
	if err != nil {
		return diag.FromErr(err)
	}

	sd := client.MakeStrForSD(datamap.SensitiveData)
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
	// c := m.(*client.Client)

	// if d.HasChange("labels") {
	// 	labels := d.Get("labels").([]interface{})
	// 	datamapLabels := []client.DatamapLabel{}

	// 	for _, label := range labels {
	// 		labelMap := label.(map[string]interface{})

	// 		labelInfoList := labelMap["label_info"].([]interface{})
	// 		labelInfos := []client.LabelInfo{}

	// 		for _, labelInfo := range labelInfoList {
	// 			labelInfoMap := labelInfo.(map[string]interface{})

	// 			attrs := labelInfoMap["attributes"].([]interface{})
	// 			attributes := []string{}

	// 			for _, attr := range attrs {
	// 				attributes = append(attributes, attr.(string))
	// 			}

	// 			li := client.LabelInfo{
	// 				Repo:       labelInfoMap["repo"].(string),
	// 				Attributes: attributes,
	// 			}

	// 			labelInfos = append(labelInfos, li)
	// 		}

	// 		dl := client.DatamapLabel{
	// 			Name: labelMap["label_id"].(string),
	// 			Info: labelInfos,
	// 		}

	// 		datamapLabels = append(datamapLabels, dl)
	// 	}
	// 	_, err := c.UpdateDatamap(datamapLabels)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}

	// 	d.Set("last_updated", time.Now().Format(time.RFC850))
	// }

	return resourceDatamapRead(ctx, d, m)
}

func resourceDatamapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	err := c.DeleteDatamap()
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags

}
