package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreatePolicyResponse struct {
	ID string `json:"ID"`
}

type Policy struct {
	Meta *PolicyMetadata `json:"meta" yaml:"meta"`
	Data []string        `json:"data,omitempty" yaml:"data,omitempty,flow"`
}

type PolicyMetadata struct {
	ID          string            `json:"id" yaml:"id"`
	Name        string            `json:"name" yaml:"name"`
	Version     string            `json:"version" yaml:"version"`
	Created     time.Time         `json:"created" yaml:"created"`
	LastUpdated time.Time         `json:"lastUpdated" yaml:"lastUpdated"`
	Type        string            `json:"type" yaml:"type"`
	Tags        []string          `json:"tags" yaml:"tags"`
	Enabled     bool              `json:"enabled" yaml:"enabled"`
	Description string            `json:"description" yaml:"description"`
	Properties  map[string]string `json:"properties" yaml:"properties"`
}

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"properties": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyCreate")
	c := m.(*client.Client)

	d.Set("type", "terraform")
	policy := getPolicyInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, policy)
	if err != nil {
		return createError("Unable to create policy", fmt.Sprintf("%v", err))
	}

	response := CreatePolicyResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourcePolicyCreate")

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError("Unable to read policy", fmt.Sprintf("%v", err))
	}

	response := Policy{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	propertiesList := flattenProperties(response.Meta)
	log.Printf("[DEBUG] resourcePolicyRead - policy: %#v", propertiesList)

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

	log.Printf("[DEBUG] End resourcePolicyRead")
	return diag.Diagnostics{}
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyUpdate")
	c := m.(*client.Client)

	d.Set("type", "terraform")
	policy := getPolicyInfoFromResource(d)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	_, err := c.DoRequest(url, http.MethodPut, policy)
	if err != nil {
		return createError("Unable to update policy", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyUpdate")

	return resourcePolicyRead(ctx, d, m)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete policy", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourcePolicyDelete")

	return diag.Diagnostics{}
}

func getStrListFromSchemaField(d *schema.ResourceData, field string) []string {
	interfaceList := d.Get(field).([]interface{})
	strList := []string{}

	for _, v := range interfaceList {
		strList = append(strList, v.(string))
	}

	return strList
}

func getPolicyInfoFromResource(d *schema.ResourceData) Policy {
	data := getStrListFromSchemaField(d, "data")
	tags := getStrListFromSchemaField(d, "tags")

	propertiesList := d.Get("properties").(*schema.Set).List()
	properties := make(map[string]string)
	for _, property := range propertiesList {
		propertyMap := property.(map[string]interface{})

		name := propertyMap["name"].(string)
		properties[name] = propertyMap["description"].(string)
	}

	policy := Policy{
		Data: data,
		Meta: &PolicyMetadata{
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

func flattenProperties(policyMetadata *PolicyMetadata) []interface{} {
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
