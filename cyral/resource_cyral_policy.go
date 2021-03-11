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
	Meta             *PolicyMetadata `json:"meta" yaml:"meta"`
	Data             []string        `json:"data,omitempty" yaml:"data,omitempty,flow"`
	WhitelistedUsers []string        `json:"whitelistedUsers,omitempty" yaml:"whitelistedUsers,omitempty"`
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
			"data": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"meta": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true, // TODO: is this correct?
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lastUpdated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"properties": {
							Type: schema.TypeSet,
							// TODO: is this optional?
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
							Required: true, // TODO: is this correct?
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"whitelistedUsers": {
				Type: schema.TypeList,
				// TODO: is this optional?
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

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyCreate")
	c := m.(*client.Client)

	policy := getPolicyInfoFromResource(d)

	// sd := sensitiveData.String()
	// log.Printf("[DEBUG] resourcePolicyCreate - sensitiveData: %s", sd)

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

	// sd := datamap.SensitiveData.String()
	// log.Printf("[DEBUG] resourceDatamapRead - sensitiveData: %s", sd)

	//TODO
	policy := flattenPolicy(&response)
	log.Printf("[DEBUG] resourcePolicyRead - policy: %s", policy)

	// if err := d.Set("mapping", datamapLabels); err != nil {
	// 	return createError("Unable to read policy", fmt.Sprintf("%v", err))
	// }

	return diag.Diagnostics{}
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourcePolicyUpdate")
	c := m.(*client.Client)

	if d.HasChange("mapping") {
		policy := getPolicyInfoFromResource(d)

		// sd := sensitiveData.String()
		// log.Printf("[DEBUG] resourcePolicyUpdate - sensitiveData: %s", sd)

		url := fmt.Sprintf("https://%s/v1/policies/%s", c.ControlPlane, d.Id())

		_, err := c.DoRequest(url, http.MethodPut, policy)
		if err != nil {
			return createError("Unable to update policy", fmt.Sprintf("%v", err))
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
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

func getPolicyInfoFromResource(d *schema.ResourceData) Policy {
	var policy Policy
	return policy
}

func flattenPolicy(policy *Policy) []interface{} {
	return make([]interface{}, 0)
}
