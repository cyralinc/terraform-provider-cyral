package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreatePolicyRuleResponse struct {
	ID string `json:"ID"`
}

type PolicyRule struct {
	Deletes    []Rule     `json:"deletes,omitempty"`
	Hosts      []string   `json:"hosts,omitempty"`
	Identities []Identity `json:"identities,omitempty"`
	Reads      []Rule     `json:"reads,omitempty"`
	RuleID     string     `json:"ruleId"`
	Updates    []Rule     `json:"updates,omitempty"`
}

type Rule struct {
	AdditionalChecks string           `json:"additionalChecks"`
	Data             []string         `json:"data,omitempty"`
	DatasetRewrites  []DatasetRewrite `json:"datasetRewrite,omitempty"`
	Rows             int              `json:"rows"`
	Severity         string           `json:"severity"`
}

type DatasetRewrite struct {
	Dataset      string   `json:"dataset"`
	Parameters   []string `json:"parameters,omitempty"`
	Repo         string   `json:"repo"`
	Substitution string   `json:"substitution"`
}

type Identity struct {
	DBRoles  string `json:"dbRoles,omitempty"`
	Groups   string `json:"groups,omitempty"`
	Services string `json:"services,omitempty"`
	Users    string `json:"users,omitempty"`
}

func resourcePolicyRule() *schema.Resource {
	ruleSchema := &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"additional_checks": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"data": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"dataset_rewrites": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"repo": {
								Type:     schema.TypeString,
								Required: true,
							},
							"substitution": {
								Type:     schema.TypeString,
								Required: true,
							},
							"dataset": {
								Type:     schema.TypeString,
								Required: true,
							},
							"parameters": {
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"rows": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"severity": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}

	return &schema.Resource{
		CreateContext: resourcePolicyRuleCreate,
		ReadContext:   resourcePolicyRuleRead,
		UpdateContext: resourcePolicyRuleUpdate,
		DeleteContext: resourcePolicyRuleDelete,
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"delete": ruleSchema,
			"read":   ruleSchema,
			"update": ruleSchema,
			"identities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_roles": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"groups": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"users": {
							Type:     schema.TypeList,
							Optional: true,
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

func resourcePolicyRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleCreate")
	c := m.(*client.Client)

	policyID := d.Get("policy_id").(string)
	resourceData := getPolicyRuleInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies/%s/rules", c.ControlPlane, policyID)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create policy", fmt.Sprintf("%v", err))
	}

	response := CreatePolicyRuleResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)
	//d.Set("created", time.Now().Format(time.RFC850))

	log.Printf("[DEBUG] End resourcePolicyRuleCreate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to read policy", fmt.Sprintf("%v", err))
	}

	response := PolicyRule{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	propertiesList := flattenPropertiesRule(response.Meta)
	log.Printf("[DEBUG] resourcePolicyRuleRead - policy: %#v", propertiesList)

	d.Set("created", response.Meta.Created.String())
	d.Set("data", response.Data)
	d.Set("description", response.Meta.Description)
	d.Set("enabled", response.Meta.Enabled)
	d.Set("last_updated", response.Meta.LastUpdated.String())
	d.Set("name", response.Meta.Name)
	d.Set("properties", propertiesList)
	d.Set("tags", response.Meta.Tags)
	d.Set("type", response.Meta.Type)
	d.Set("version", response.Meta.Version)

	log.Printf("[DEBUG] End resourcePolicyRuleRead")
	return diag.Diagnostics{}
}

func resourcePolicyRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleUpdate")
	c := m.(*client.Client)

	d.Set("type", "terraform")
	policy := getPolicyRuleInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	_, err := c.DoRequest(url, http.MethodPut, policy)
	if err != nil {
		return createError("Unable to update policy", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyRuleUpdate")

	return resourcePolicyRuleRead(ctx, d, m)
}

func resourcePolicyRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRuleDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete policy", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyRuleDelete")

	return diag.Diagnostics{}
}

func getPolicyRuleInfoFromResource(d *schema.ResourceData) PolicyRule {
	data := getStrListFromSchemaField(d, "data")
	tags := getStrListFromSchemaField(d, "tags")

	propertiesList := d.Get("properties").(*schema.Set).List()
	properties := make(map[string]string)
	for _, property := range propertiesList {
		propertyMap := property.(map[string]interface{})

		name := propertyMap["name"].(string)
		properties[name] = propertyMap["description"].(string)
	}

	policy := PolicyRule{
		Data: data,
		Meta: &PolicyRuleMetadata{
			Tags:       tags,
			Properties: properties,
		},
	}

	if v, ok := d.Get("name").(string); ok {
		policy.Meta.Name = v
	}

	if v, ok := d.Get("version").(string); ok {
		policy.Meta.Version = v
	}

	if v, ok := d.Get("type").(string); ok {
		policy.Meta.Type = v
	}

	if v, ok := d.Get("enabled").(bool); ok {
		policy.Meta.Enabled = v
	}

	if v, ok := d.Get("description").(string); ok {
		policy.Meta.Description = v
	}

	return policy
}

func flattenPropertiesRule(policyMetadata *PolicyRuleMetadata) []interface{} {
	if policyMetadata != nil {

		propertiesList := make([]map[string]interface{},
			len(policyMetadata.Properties), len(policyMetadata.Properties))

		for name, description := range policyMetadata.Properties {
			propertyMap := make(map[string]interface{})

			propertyMap["name"] = name
			propertyMap["description"] = description
			propertiesList = append(propertiesList, propertyMap)
		}

	}

	return make([]interface{}, 0)
}
