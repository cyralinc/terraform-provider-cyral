package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSidecarResponse struct {
	ID        string `json:"ID"`
	AccessKey string `json:"accessKey"`
}

type SidecarData struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	SidecarProperty SidecarProperty `json:"properties"`
}

type SidecarProperty struct {
	DeploymentMethod     string `json:"deploymentMethod"`
	AWSRegion            string `json:"awsRegion"`
	KeyName              string `json:"keyName"`
	VPC                  string `json:"vpc"`
	Subnets              string `json:"publicSubnets"`
	PubliclyAccessible   string `json:"publiclyAccessible"`
	MetricsIntegrationID string `json:"metricsIntegrationID"`
	LogIntegrationID     string `json:"logIntegrationID"`
}

func resourceSidecar() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSidecarCreate,
		ReadContext:   resourceSidecarRead,
		UpdateContext: resourceSidecarUpdate,
		DeleteContext: resourceSidecarDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_integration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"metrics_integration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"aws_configuration": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"aws_region": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"vpc": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"subnets": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"publicly_accessible": {
							Type:     schema.TypeBool,
							Required: true,
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

func resourceSidecarCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarCreate")
	c := m.(*client.Client)

	resourceData, err := getSidecarDataFromResource(c, d)
	if err != nil {
		return createError("Unable to create sidecar", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create sidecar", fmt.Sprintf("%v", err))
	}

	response := CreateSidecarResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read sidecar. SidecarID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("deployment_method", response.SidecarProperty.DeploymentMethod)
	d.Set("log_integration_id", response.SidecarProperty.LogIntegrationID)
	d.Set("metrics_integration_id", response.SidecarProperty.MetricsIntegrationID)
	awsConfiguration := make([]map[string]interface{}, 0)
	ac := make(map[string]interface{})
	ac["key_name"] = response.SidecarProperty.KeyName
	ac["aws_region"] = response.SidecarProperty.AWSRegion
	ac["vpc"] = response.SidecarProperty.VPC
	ac["subnets"] = response.SidecarProperty.Subnets
	if p, err := strconv.ParseBool(response.SidecarProperty.PubliclyAccessible); err == nil {
		ac["publicly_accessible"] = p
	}
	awsConfiguration = append(awsConfiguration, ac)
	d.Set("aws_configuration", awsConfiguration)

	log.Printf("[DEBUG] End resourceSidecarRead")

	return diag.Diagnostics{}
}

func resourceSidecarUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarUpdate")
	c := m.(*client.Client)

	resourceData, err := getSidecarDataFromResource(c, d)
	if err != nil {
		return createError("Unable to update sidecar", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	if _, err = c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update sidecar", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceSidecarUpdate")

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete sidecar", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceSidecarDelete")

	return diag.Diagnostics{}
}

func getSidecarDataFromResource(c *client.Client, d *schema.ResourceData) (SidecarData, error) {
	deploymentMethod := d.Get("deployment_method").(string)
	if err := client.ValidateDeploymentMethod(deploymentMethod); err != nil {
		return SidecarData{}, err
	}

	sp := SidecarProperty{
		DeploymentMethod:     deploymentMethod,
		LogIntegrationID:     d.Get("log_integration_id").(string),
		MetricsIntegrationID: d.Get("metrics_integration_id").(string),
	}

	if v, ok := d.GetOk("aws_configuration"); ok {
		vL := v.(*schema.Set).List()
		for _, v := range vL {
			configMap := v.(map[string]interface{})
			if v, ok := configMap["key_name"].(string); ok && v != "" {
				sp.KeyName = v
			}
			if v, ok := configMap["vpc"].(string); ok && v != "" {
				sp.VPC = v
			}
			if v, ok := configMap["subnets"].(string); ok && v != "" {
				sp.Subnets = v
			}
			if v, ok := configMap["publicly_accessible"].(bool); ok {
				sp.PubliclyAccessible = strconv.FormatBool(v)
			}
			if v, ok := configMap["aws_region"].(string); ok && v != "" {
				if err := client.ValidateAWSRegion(v); err != nil {
					return SidecarData{}, err
				}
				sp.AWSRegion = v
			}
		}

	}

	return SidecarData{
		ID:              d.Id(),
		Name:            d.Get("name").(string),
		SidecarProperty: sp,
	}, nil
}
