package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateSidecarResponse struct {
	ID string `json:"ID"`
}

type SidecarData struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	Labels                   []string                 `json:"labels"`
	SidecarProperty          SidecarProperty          `json:"properties"`
	UserEndpoint             string                   `json:"userEndpoint"`
	CertificateBundleSecrets CertificateBundleSecrets `json:"certificateBundleSecrets"`
}

type SidecarProperty struct {
	DeploymentMethod string `json:"deploymentMethod"`
}

type CertificateBundleSecrets map[string]CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine"`
	SecretId string `json:"secretId"`
	Type     string `json:"type"`
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
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_bundle_secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sidecar": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"secret_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
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

	response := SidecarData{}
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
		// Currently, the sidecar API always returns a status code of 500 for every error,
		// so its not possible to distinguish if the error returned is related to
		// a 404 Not Found or not by its status code. This way, a workaround for that is to
		// check if the error message matches a 'Failed to extract info for wrapper' message,
		// since thats the current message returned when the sidecar is not found. Once this
		// issue is fixed in the sidecar API, we should handle the error here by its status
		// code, and only remove the resource from the state (d.SetId("")) if it returns a 404
		// Not Found.
		matched, regexpError := regexp.MatchString("Failed to extract info for wrapper",
			err.Error())
		if regexpError == nil && matched {
			log.Printf("[DEBUG] Sidecar not found. SidecarID: %s. "+
				"Removing it from state. Error: %v", d.Id(), err)
			d.SetId("")
			return nil
		}

		return createError(fmt.Sprintf("Unable to read sidecar. SidecarID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	d.Set("deployment_method", response.SidecarProperty.DeploymentMethod)
	d.Set("labels", response.Labels)
	d.Set("user_endpoint", response.UserEndpoint)
	d.Set("certificate_bundle_secrets", flattenCertificateBundleSecrets(&response.CertificateBundleSecrets))

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

func getSidecarDataFromResource(c *client.Client, d *schema.ResourceData) (*SidecarData, error) {
	log.Printf("[DEBUG] Init getSidecarDataFromResource")
	if err := validateCertificateBundleSecretsBlock(d); err != nil {
		return nil, err
	}

	deploymentMethod := d.Get("deployment_method").(string)
	if err := client.ValidateDeploymentMethod(deploymentMethod); err != nil {
		return &SidecarData{}, err
	}

	sp := SidecarProperty{
		DeploymentMethod: deploymentMethod,
	}
	labels := d.Get("labels").([]interface{})
	sidecarDataLabels := make([]string, len(labels))
	for i, label := range labels {
		sidecarDataLabels[i] = (label).(string)
	}
	rdCBS := d.Get("certificate_bundle_secrets").(*schema.Set).List()
	cbsMap := make(map[string]CertificateBundleSecret)
	for _, c := range rdCBS {
		typeMap := c.(map[string]interface{})
		cbsType := typeMap["type"].(string)
		secretId := typeMap["secretId"].(string)
		engine := typeMap["engine"].(string)
		cbs := CertificateBundleSecret{
			SecretId: secretId,
			Engine:   engine,
			Type:     cbsType,
		}
		cbsMap[cbsType] = cbs
	}
	log.Printf("[DEBUG] end getSidecarDataFromResource")
	return &SidecarData{
		ID:                       d.Id(),
		Name:                     d.Get("name").(string),
		Labels:                   sidecarDataLabels,
		SidecarProperty:          sp,
		UserEndpoint:             d.Get("user_endpoint").(string),
		CertificateBundleSecrets: cbsMap,
	}, nil
}

func validateCertificateBundleSecretsBlock(d *schema.ResourceData) error {
	log.Printf("[DEBUG] Init validateCertificateBundleSecretsBlock")
	set := make(map[string]bool)
	var repeated []string
	rdCBS := d.Get("certificate_bundle_secrets").(*schema.Set).List()

	for _, c := range rdCBS {
		typeMap := c.(map[string]interface{})

		cbsType := typeMap["type"].(string)
		if set[cbsType] {
			repeated = append(repeated, cbsType)
		} else {
			set[cbsType] = true
		}
	}

	log.Printf("[DEBUG] end validateCertificateBundleSecretsBlock")

	if len(repeated) > 0 {
		return fmt.Errorf("there is more than one `certificate_bundle_secret`"+
			" block with the same type. Types must be unique. Repeated types: %v", repeated)
	}

	return nil
}

func flattenCertificateBundleSecrets(cbs *CertificateBundleSecrets) []interface{} {
	log.Printf("[DEBUG] Init flattenCertificateBundleSecrets")
	if cbs != nil {
		flatCBS := make([]interface{}, 0, len(*cbs))

		for key, val := range *cbs {
			cbsMap := make(map[string]interface{})

			fooCB := make(map[string]string)
			if val.SecretId != "" {
				fooCB["secret_id"] = val.SecretId
			}
			if val.Engine != "" {
				fooCB["engine"] = val.Engine
			}
			if val.Type != "" {
				fooCB["type"] = val.Type
			}
			cbsMap[key] = fooCB
			flatCBS = append(flatCBS, cbsMap)
		}

		log.Printf("[DEBUG] end flattenCertificateBundleSecrets")
		return flatCBS
	}

	log.Printf("[DEBUG] end flattenCertificateBundleSecrets")
	return make([]interface{}, 0)
}
