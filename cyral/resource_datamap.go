package cyral

import (
	"context"
	"time"

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
	}
}

func resourceDatamapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	labels := d.Get("labels").([]interface{})
	datamapLabels := []client.DatamapLabel{}

	for _, label := range labels {
		labelMap := label.(map[string]interface{})

		dl := client.DatamapLabel{
			Repo:       labelMap["repo"].(string),
			Attributes: labelMap["attributes"].([]string),
		}

		datamapLabels = append(datamapLabels, dl)
	}

	// datamap, err := c.CreateDatamap(datamapLabels)
	_, err := c.CreateDatamap(datamapLabels)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId(strconv.Itoa(datamap.ID))

	resourceDatamapRead(ctx, d, m)

	return diags
}

func flattenDatamapLabels(datamapLabels *[]client.DatamapLabel) []interface{} {
	if datamapLabels != nil {
		dmLabels := make([]interface{}, len(*datamapLabels), len(*datamapLabels))

		for i, datamapLabel := range *datamapLabels {
			dmLabel := make(map[string]interface{})

			dmLabel["repo"] = datamapLabel.Repo
			dmLabel["attributes"] = datamapLabel.Attributes
			dmLabels[i] = dmLabel
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

	datamapLabels := flattenDatamapLabels(&datamap.Labels)
	if err := d.Set("labels", datamapLabels); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatamapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	if d.HasChange("labels") {
		labels := d.Get("labels").([]interface{})
		datamapLabels := []client.DatamapLabel{}

		for _, label := range labels {
			labelMap := label.(map[string]interface{})

			dl := client.DatamapLabel{
				Repo:       labelMap["repo"].(string),
				Attributes: labelMap["attributes"].([]string),
			}

			datamapLabels = append(datamapLabels, dl)
		}

		_, err := c.UpdateDatamap(datamapLabels)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceDatamapRead(ctx, d, m)
}

func resourceDatamapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
