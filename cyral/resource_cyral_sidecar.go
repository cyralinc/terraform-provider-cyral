package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type CreateSidecarResponse struct {
	ID string `json:"ID"`
}

type SidecarData struct {
	ID                       string                   `json:"id"`
	Name                     string                   `json:"name"`
	Labels                   []string                 `json:"labels"`
	SidecarProperties        *SidecarProperties       `json:"properties"`
	ServicesConfig           SidecarServicesConfig    `json:"services"`
	UserEndpoint             string                   `json:"userEndpoint"`
	CertificateBundleSecrets CertificateBundleSecrets `json:"certificateBundleSecrets,omitempty"`
}

func (sd *SidecarData) BypassMode() string {
	if sd.ServicesConfig != nil {
		if dispConfig, ok := sd.ServicesConfig["dispatcher"]; ok {
			if bypass_mode, ok := dispConfig["bypass"]; ok {
				return bypass_mode
			}
		}
	}
	return ""
}

type SidecarProperties struct {
	DeploymentMethod string `json:"deploymentMethod"`
	LogIntegrationID string `json:"logIntegrationID,omitempty"`
}

func NewSidecarProperties(deploymentMethod, logIntegrationID string) *SidecarProperties {
	return &SidecarProperties{
		DeploymentMethod: deploymentMethod,
		LogIntegrationID: logIntegrationID,
	}
}

type SidecarServicesConfig map[string]map[string]string

type CertificateBundleSecrets map[string]*CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
}

func resourceSidecar() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [sidecars](https://cyral.com/docs/sidecars/sidecar-manage).",
		CreateContext: resourceSidecarCreate,
		ReadContext:   resourceSidecarRead,
		UpdateContext: resourceSidecarUpdate,
		DeleteContext: resourceSidecarDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Sidecar name that will be used internally in Control Plane (ex: `your_sidecar_name`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"deployment_method": {
				Description: "Deployment method that will be used by this sidecar (valid values: `docker`, `cloudFormation`, `terraform`, `helm`, `helm3`, `automated`, `custom`, `terraformGKE`, `linux`, and `singleContainer`).",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"docker", "cloudFormation", "terraform", "helm", "helm3",
						"automated", "custom", "terraformGKE", "singleContainer",
						"linux",
					}, false,
				),
			},
			"log_integration_id": {
				Description: "ID of the log integration mapped to this sidecar.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"labels": {
				Description: "Labels that can be attached to the sidecar and shown in the `Tags` field in the UI.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_endpoint": {
				Description: "User-defined endpoint (also referred as `alias`) that can be used to override the sidecar DNS endpoint shown in the UI.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"bypass_mode": {
				Description: "This argument lets you specify how to handle the connection in the event of an error in the sidecar during a user’s session. Valid modes are: `always`, `failover` or `never`. Defaults to `failover`. If `always` is specified, the sidecar will run in [passthrough mode](https://cyral.com/docs/sidecars/sidecar-manage#passthrough-mode). If `failover` is specified, the sidecar will run in [resiliency mode](https://cyral.com/docs/sidecars/sidecar-manage#resilient-mode-of-sidecar-operation). If `never` is specified and there is an error in the sidecar, connections to bound repositories will fail.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "failover",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"always",
						"failover",
						"never",
					}, false,
				),
			},
			"certificate_bundle_secrets": {
				Description: "Certificate Bundle Secret is a configuration that holds data about the" +
					" location of a particular TLS certificate bundle in a secrets manager.",
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sidecar": {
							Description: "Certificate Bundle Secret for sidecar.",
							Type:        schema.TypeSet,
							MaxItems:    1,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine": {
										Description: "Engine is the name of the engine used with the given secrets" +
											" manager type, when applicable.",
										Type:     schema.TypeString,
										Optional: true,
									},
									"secret_id": {
										Description: "Secret ID is the identifier or location for the secret that" +
											" holds the certificate bundle.",
										Type:     schema.TypeString,
										Required: true,
									},
									"type": {
										Description: "Type identifies the secret manager used to store the secret. Valid values are: `aws` and `k8s`.",
										Type:        schema.TypeString,
										Required:    true,
										ValidateFunc: validation.StringInSlice(
											[]string{
												"aws",
												"k8s",
											}, false,
										),
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
		matched, regexpError := regexp.MatchString(
			"Failed to extract info for wrapper",
			err.Error(),
		)
		if regexpError == nil && matched {
			log.Printf(
				"[DEBUG] Sidecar not found. SidecarID: %s. "+
					"Removing it from state. Error: %v", d.Id(), err,
			)
			d.SetId("")
			return nil
		}

		return createError(
			fmt.Sprintf(
				"Unable to read sidecar. SidecarID: %s",
				d.Id(),
			), fmt.Sprintf("%v", err),
		)
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("name", response.Name)
	if properties := response.SidecarProperties; properties != nil {
		d.Set("deployment_method", properties.DeploymentMethod)
		d.Set("log_integration_id", properties.LogIntegrationID)
	}
	d.Set("labels", response.Labels)
	d.Set("user_endpoint", response.UserEndpoint)
	if bypassMode := response.BypassMode(); bypassMode != "" {
		d.Set("bypass_mode", bypassMode)
	}
	d.Set("certificate_bundle_secrets", flattenCertificateBundleSecrets(response.CertificateBundleSecrets))

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
	logIntegrationID := d.Get("log_integration_id").(string)

	properties := NewSidecarProperties(deploymentMethod, logIntegrationID)

	svcconf := SidecarServicesConfig{
		"dispatcher": map[string]string{
			"bypass": d.Get("bypass_mode").(string),
		},
	}

	labels := d.Get("labels").([]interface{})
	sidecarDataLabels := []string{}
	for _, labelInterface := range labels {
		if label, ok := labelInterface.(string); ok {
			sidecarDataLabels = append(sidecarDataLabels, label)
		}
	}

	cbs := getCertificateBundleSecret(d)

	log.Printf("[DEBUG] end getSidecarDataFromResource")
	return &SidecarData{
		ID:                       d.Id(),
		Name:                     d.Get("name").(string),
		Labels:                   sidecarDataLabels,
		SidecarProperties:        properties,
		ServicesConfig:           svcconf,
		UserEndpoint:             d.Get("user_endpoint").(string),
		CertificateBundleSecrets: cbs,
	}, nil
}

func flattenCertificateBundleSecrets(cbs CertificateBundleSecrets) []interface{} {
	log.Printf("[DEBUG] Init flattenCertificateBundleSecrets")
	var flatCBS []interface{}
	if cbs != nil {
		cb := make(map[string]interface{})

		for key, val := range cbs {
			// Ignore self-signed certificates
			if key != "sidecar-generated-selfsigned" {
				contentCB := make([]interface{}, 1)

				log.Printf("[DEBUG] key: %v", key)
				log.Printf("[DEBUG] val: %v", val)

				contentCBMap := make(map[string]interface{})
				contentCBMap["secret_id"] = val.SecretId
				contentCBMap["engine"] = val.Engine
				contentCBMap["type"] = val.Type

				contentCB[0] = contentCBMap
				cb[key] = contentCB
			}
		}

		if len(cb) > 0 {
			flatCBS = make([]interface{}, 1)
			flatCBS[0] = cb
		}
	}

	log.Printf("[DEBUG] end flattenCertificateBundleSecrets %v", flatCBS)
	return flatCBS
}

func getCertificateBundleSecret(d *schema.ResourceData) CertificateBundleSecrets {
	log.Printf("[DEBUG] Init getCertificateBundleSecret")
	rdCBS := d.Get("certificate_bundle_secrets").(*schema.Set).List()
	ret := make(CertificateBundleSecrets)

	if len(rdCBS) > 0 {
		cbsMap := rdCBS[0].(map[string]interface{})
		for k, v := range cbsMap {
			vList := v.(*schema.Set).List()
			// 1. k = "sidecar" or other direct internal elements of certificate_bundle_secrets
			// 2. Also one element on this list due to MaxItems...
			// 3. Ignore self signed certificates
			if len(vList) > 0 && k != "sidecar-generated-selfsigned" {
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
				ret[k] = &cbs
			}
		}
	}

	// If the occurrence of `sidecar` does not exist, set it to an empty certificate bundle
	// so that the API can remove the `sidecar` key from the persisted certificate bundle map.
	if _, ok := ret["sidecar"]; !ok {
		ret["sidecar"] = &CertificateBundleSecret{}
	}

	log.Printf("[DEBUG] end getCertificateBundleSecret")
	return ret
}
