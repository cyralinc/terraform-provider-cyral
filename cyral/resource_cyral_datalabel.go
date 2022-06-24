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

const (
	dataLabelTypeCustom = "CUSTOM"
)

type DataLabel struct {
	Name        string
	Type        string   `json:"type"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (dl *DataLabel) WriteToSchema(d *schema.ResourceData) error {
	if err := d.Set("id", dl.Name); err != nil {
		return err
	}
	if err := d.Set("description", dl.Description); err != nil {
		return err
	}

	var tagIfaces []interface{}
	for _, tag := range dl.Tags {
		tagIfaces = append(tagIfaces, tag)
	}
	if err := d.Set("tags", tagIfaces); err != nil {
		return err
	}

	return nil
}

func resourceDatalabel() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages data labels. Data labels are part of the Cyral [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: resourceDatalabelCreate,
		ReadContext:   resourceDatalabelRead,
		UpdateContext: resourceDatalabelUpdate,
		DeleteContext: resourceDatalabelDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Name of the data label.",
				Type:        schema.TypeString,
				Computed:    true,
			},
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
				Description: "Tags that can be used to group data labels.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDatalabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceDatalabelCreate")
	c := m.(*client.Client)

	labelName := d.Get("name").(string)
	url := fmt.Sprintf("https://%s/v1/datalabels/%s", c.ControlPlane, labelName)

	dataLabel := getDataLabelFromResource(d)

	_, err := c.DoRequest(url, http.MethodPut, dataLabel)
	if err != nil {
		return createError("Unable to create data label", err.Error())
	}

	d.SetId(labelName)

	log.Printf("[DEBUG] End resourceDatalabelCreate")

	return resourceDatalabelRead(ctx, d, m)
}

func resourceDatalabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceDatalabelRead")
	c := m.(*client.Client)

	labelName := d.Get("name").(string)
	url := fmt.Sprintf("https://%s/v1/datalabels/%s", c.ControlPlane, labelName)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to create data label", err.Error())
	}

	dataLabel := DataLabel{}
	if err := json.Unmarshal(body, &dataLabel); err != nil {
		return createError("Unable to unmarshall JSON", err.Error())
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", dataLabel)

	if err := dataLabel.WriteToSchema(d); err != nil {
		return createError("Unable to read data label", err.Error())
	}

	log.Printf("[DEBUG] End resourceDatalabelRead")

	return diag.Diagnostics{}
}

func resourceDatalabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceDatalabelUpdate")
	c := m.(*client.Client)

	dataLabel := getDataLabelFromResource(d)

	labelName := d.Get("name").(string)
	url := fmt.Sprintf("https://%s/v1/datalabels/%s", c.ControlPlane, labelName)

	_, err := c.DoRequest(url, http.MethodPut, dataLabel)
	if err != nil {
		return createError("Unable to create data label", err.Error())
	}

	log.Printf("[DEBUG] End resourceDatalabelUpdate")

	return resourceDatalabelRead(ctx, d, m)
}

func resourceDatalabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceDatalabelDelete")
	c := m.(*client.Client)

	labelName := d.Get("name").(string)
	url := fmt.Sprintf("https://%s/v1/datalabels/%s", c.ControlPlane, labelName)

	_, err := c.DoRequest(url, http.MethodDelete, nil)
	if err != nil {
		return createError("Unable to delete data label", err.Error())
	}

	log.Printf("[DEBUG] End resourceDatalabelDelete")

	return diag.Diagnostics{}
}

func getDataLabelFromResource(d *schema.ResourceData) DataLabel {
	var tags []string
	tagIfaces := d.Get("tags").([]interface{})
	for _, tagIface := range tagIfaces {
		tags = append(tags, tagIface.(string))
	}

	return DataLabel{
		Name:        d.Get("name").(string),
		Type:        dataLabelTypeCustom,
		Description: d.Get("description").(string),
		Tags:        tags,
	}
}
