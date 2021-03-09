package cyral

import (
	"context"
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
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"aws_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"vpc": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"subnets": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"publicly_accessible": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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

	response := CreateSidecarResponse{}
	if err := c.CreateResource(url, http.MethodPost, resourceData, &response); err != nil {
		return createError("Unable to create sidecar", fmt.Sprintf("%v", err))
	}

	d.SetId(response.ID)

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	response := SidecarData{}
	if err := c.ReadResource(url, &response); err != nil {
		return createError(fmt.Sprintf("Unable to read sidecar. SidecarID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	d.Set("name", response.Name)
	d.Set("deployment_method", response.SidecarProperty.DeploymentMethod)
	d.Set("key_name", response.SidecarProperty.KeyName)
	d.Set("aws_region", response.SidecarProperty.AWSRegion)
	d.Set("vpc", response.SidecarProperty.VPC)
	d.Set("subnets", response.SidecarProperty.Subnets)
	d.Set("publicly_accessible", response.SidecarProperty.PubliclyAccessible)
	d.Set("log_integration_id", response.SidecarProperty.LogIntegrationID)
	d.Set("metrics_integration_id", response.SidecarProperty.MetricsIntegrationID)

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
	if err = c.UpdateResource(resourceData, url); err != nil {
		return createError("Unable to update sidecar", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] End resourceSidecarUpdate")

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceSidecarDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())
	if err := c.DeleteResource(url); err != nil {
		return createError("Unable to delete sidecar", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceSidecarDelete")

	return diag.Diagnostics{}
}

func getSidecarDataFromResource(c *client.Client, d *schema.ResourceData) (SidecarData, error) {
	deploymentMethod := d.Get("deployment_method").(string)
	if err := validateDeploymentMethod(deploymentMethod); err != nil {
		return SidecarData{}, err
	}
	awsRegion := d.Get("aws_region").(string)
	if awsRegion != "" {
		if err := validateAWSRegion(awsRegion); err != nil {
			return SidecarData{}, err
		}
	}
	logIntegrationID := d.Get("log_integration_id").(string)
	metricsIntegrationID := d.Get("metrics_integration_id").(string)

	return SidecarData{
		ID:   d.Id(),
		Name: d.Get("name").(string),
		SidecarProperty: SidecarProperty{
			DeploymentMethod:     deploymentMethod,
			KeyName:              d.Get("key_name").(string),
			AWSRegion:            awsRegion,
			VPC:                  d.Get("vpc").(string),
			Subnets:              d.Get("subnets").(string),
			PubliclyAccessible:   strconv.FormatBool(d.Get("publicly_accessible").(bool)),
			LogIntegrationID:     logIntegrationID,
			MetricsIntegrationID: metricsIntegrationID,
		},
	}, nil
}

func validateDeploymentMethod(param string) error {
	validValues := map[string]bool{
		"docker":         true,
		"cloudformation": true,
		"terraform":      true,
		"helm":           true,
		"helm3":          true,
		"automated":      true,
		"custom":         true,
		"terraformGKE":   true,
	}
	if validValues[param] == false {
		return fmt.Errorf("deployment method must be one of %v", validValues)
	}
	return nil
}

func validateAWSRegion(param string) error {
	validValues := map[string]bool{
		"us-east-2":      true,
		"us-east-1":      true,
		"us-west-1":      true,
		"us-west-2":      true,
		"af-south-1":     true,
		"ap-east-1":      true,
		"ap-south-1":     true,
		"ap-northeast-3": true,
		"ap-northeast-2": true,
		"ap-southeast-1": true,
		"ap-southeast-2": true,
		"ap-northeast-1": true,
		"ca-central-1":   true,
		"eu-central-1":   true,
		"eu-west-1":      true,
		"eu-west-2":      true,
		"eu-south-1":     true,
		"eu-west-3":      true,
		"eu-north-1":     true,
		"me-south-1":     true,
		"sa-east-1":      true,
	}
	if validValues[param] == false {
		return fmt.Errorf("AWS region must be one of %v", validValues)
	}
	return nil
}
