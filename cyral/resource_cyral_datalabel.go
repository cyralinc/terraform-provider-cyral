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

func resourceDatalabel() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages data labels. Data labels are part of the Cyral [Data Map](https://cyral.com/docs/policy/datamap).",
		CreateContext: resourceDatalabelCreate,
		ReadContext:   resourceDatalabelRead,
		UpdateContext: resourceDatalabelUpdate,
		DeleteContext: resourceDatalabelDelete,
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
				Description: "Tags that can be used to categorize data labels.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"classification_rule": {
				Description: "Classification rules are used by the " +
					"[Automatic Data Map](https://cyral.com/docs/policy/automatic-datamap) feature to automatically map " +
					"data locations to labels. Currently, only `PREDEFINED` labels have classification rules.",
				Optional: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_type": {
							Description: "Type of the classification rule. Valid values are: `UNKNOWN` and `REGO`.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"rule_code": {
							Description: "Actual code of the classification rule. For example, this attribute may contain " +
								"REGO code for `REGO`-type classification rules.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_status": {
							Description: "Status of the classification rule. Valid values are: `ENABLED` and  `DISABLED`.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(
				ctx context.Context,
				d *schema.ResourceData,
				m any,
			) ([]*schema.ResourceData, error) {
				d.Set("name", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceDatalabelCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceDatalabelRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

	if err := writeDataLabelToResourceSchema(dataLabel, d); err != nil {
		return createError("Unable to read data label", err.Error())
	}

	log.Printf("[DEBUG] End resourceDatalabelRead")

	return diag.Diagnostics{}
}

func resourceDatalabelUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceDatalabelDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
	tagIfaces := d.Get("tags").([]any)
	for _, tagIface := range tagIfaces {
		tags = append(tags, tagIface.(string))
	}

	var classificationRule *DataLabelClassificationRule
	classificationRuleList := d.Get("classification_rule").(*schema.Set).List()
	if len(classificationRuleList) > 0 {
		classificationRuleMap := classificationRuleList[0].(map[string]any)
		classificationRule = &DataLabelClassificationRule{
			RuleType:   classificationRuleMap["rule_type"].(string),
			RuleCode:   classificationRuleMap["rule_code"].(string),
			RuleStatus: classificationRuleMap["rule_status"].(string),
		}
	}

	return DataLabel{
		Name:               d.Get("name").(string),
		Type:               dataLabelTypeCustom,
		Description:        d.Get("description").(string),
		Tags:               tags,
		ClassificationRule: classificationRule,
	}
}

func writeDataLabelToResourceSchema(label DataLabel, d *schema.ResourceData) error {
	if err := d.Set("description", label.Description); err != nil {
		return fmt.Errorf("error setting 'description' field: %w", err)
	}

	tagIfaces := label.TagsAsInterface()
	if err := d.Set("tags", tagIfaces); err != nil {
		return fmt.Errorf("error setting 'tags' field: %w", err)
	}

	classificationRule := label.ClassificationRuleAsInterface()
	if err := d.Set("classification_rule", classificationRule); err != nil {
		return fmt.Errorf("error setting 'classification_rule' field: %w", err)
	}

	return nil
}
