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
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
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
				Description: "Certificate bundle secrets details.",
				Type:        schema.TypeSet,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sidecar": {
							Description: "eneral-purpose certificate bundle for the sidecar.",
							Type:        schema.TypeSet,
							MaxItems:    1,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"secret_id": {
										Description: "Secret identification for the given `type`.",
										Type:        schema.TypeString,
										Required:    true,
									},
									"type": {
										Description: "Secret type. Valid values are `aws` and `k8s`.",
										Type:        schema.TypeString,
										Required:    true,
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

	cbs := getCertificateBundleSecret(d)

	log.Printf("[DEBUG] end getSidecarDataFromResource")
	return &SidecarData{
		ID:                       d.Id(),
		Name:                     d.Get("name").(string),
		Labels:                   sidecarDataLabels,
		SidecarProperty:          sp,
		UserEndpoint:             d.Get("user_endpoint").(string),
		CertificateBundleSecrets: *cbs,
	}, nil
}

func flattenCertificateBundleSecrets(cbs *CertificateBundleSecrets) []interface{} {
	log.Printf("[DEBUG] Init flattenCertificateBundleSecrets")
	if cbs != nil {
		flatCBS := make([]interface{}, 1)
		cb := make(map[string]interface{})

		contentCB := make([]interface{}, 0, len(*cbs))
		for key, val := range *cbs {
			log.Printf("[DEBUG] key: %v", key)
			log.Printf("[DEBUG] val: %v", val)

			contentCBMap := make(map[string]interface{})
			contentCBMap["secret_id"] = val.SecretId
			contentCBMap["engine"] = val.Engine
			contentCBMap["type"] = val.Type

			contentCB = append(contentCB, contentCBMap)
		}

		cb["sidecar"] = contentCB
		flatCBS[0] = cb

		log.Printf("[DEBUG] end flattenCertificateBundleSecrets %v", flatCBS)
		return flatCBS
	}

	log.Printf("[DEBUG] end flattenCertificateBundleSecrets")
	return nil
}

func getCertificateBundleSecret(d *schema.ResourceData) *CertificateBundleSecrets {
	log.Printf("[DEBUG] Init getCertificateBundleSecret")
	rdCBS := d.Get("certificate_bundle_secrets").(*schema.Set).List()
	ret := make(CertificateBundleSecrets)

	if len(rdCBS) > 0 {
		cbsMap := rdCBS[0].(map[string]interface{})
		for k, v := range cbsMap {
			vList := v.(*schema.Set).List()
			// k = "sidecar" or other direct internal elements of certificate_bundle_secrets
			// Also one element on this list due to MaxItems...
			if len(vList) > 0 {
				vMap := vList[0].(map[string]interface{})
				engine := ""
				if val, ok := vMap["engine"]; val != nil && ok {
					engine = val.(string)
				}
				cbsType := vMap["type"].(string)
				secretId := vMap["secret_id"].(string)
				cbs := CertificateBundleSecret{
					SecretId: secretId,
					Engine:   engine,
					Type:     cbsType,
				}
				ret[k] = cbs
			}
		}
	}

	log.Printf("[DEBUG] end getCertificateBundleSecret")
	return &ret
}
