package sidecar

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
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
	DeploymentMethod           string `json:"deploymentMethod"`
	LogIntegrationID           string `json:"logIntegrationID,omitempty"`
	DiagnosticLogIntegrationID string `json:"diagnosticLogIntegrationID,omitempty"`
}

func NewSidecarProperties(deploymentMethod, activityLogIntegrationID, diagnosticLogIntegrationID string) *SidecarProperties {
	return &SidecarProperties{
		DeploymentMethod:           deploymentMethod,
		LogIntegrationID:           activityLogIntegrationID,
		DiagnosticLogIntegrationID: diagnosticLogIntegrationID,
	}
}

type SidecarServicesConfig map[string]map[string]string

type CertificateBundleSecrets map[string]*CertificateBundleSecret

type CertificateBundleSecret struct {
	Engine   string `json:"engine,omitempty"`
	SecretId string `json:"secretId,omitempty"`
	Type     string `json:"type,omitempty"`
}

func ResourceSidecar() *schema.Resource {
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
				Description: "Deployment method that will be used by this sidecar (valid values: `docker`, `cft-ec2`, `terraform`, `helm3`, `automated`, `custom`, `terraformGKE`, `linux`, and `singleContainer`).",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"docker", "cft-ec2", "terraform", "helm3",
						"automated", "custom", "terraformGKE", "singleContainer",
						"linux",
					}, false,
				),
			},
			"log_integration_id": {
				Description:   "ID of the log integration mapped to this sidecar, used for Cyral activity logs.",
				Type:          schema.TypeString,
				Optional:      true,
				Deprecated:    "Since sidecar v4.8. Use `activity_log_integration_id` instead.",
				ConflictsWith: []string{"activity_log_integration_id"},
			},
			"activity_log_integration_id": {
				Description: "ID of the log integration mapped to this sidecar, used for Cyral activity logs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"diagnostic_log_integration_id": {
				Description: "ID of the log integration mapped to this sidecar, used for sidecar diagnostic logs.",
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
				Deprecated: "Since sidecar v4.7 the certificate is managed at deployment level. Refer" +
					" to [our public docs](https://cyral.com/docs/v4.7/sidecars/sidecar-certificates)" +
					" for more information.",
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
	tflog.Debug(ctx, "Init resourceSidecarCreate")
	c := m.(*client.Client)

	resourceData, err := getSidecarDataFromResource(c, d)
	if err != nil {
		return utils.CreateError("Unable to create sidecar", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)

	body, err := c.DoRequest(ctx, url, http.MethodPost, resourceData)
	if err != nil {
		return utils.CreateError("Unable to create sidecar", fmt.Sprintf("%v", err))
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.SetId(response.ID)

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceSidecarRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(ctx, url, http.MethodGet, nil)
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
			tflog.Debug(ctx, fmt.Sprintf("Sidecar not found. SidecarID: %s. "+
				"Removing it from state. Error: %v", d.Id(), err))
			d.SetId("")
			return nil
		}

		return utils.CreateError(
			fmt.Sprintf(
				"Unable to read sidecar. SidecarID: %s",
				d.Id(),
			), fmt.Sprintf("%v", err),
		)
	}

	response := SidecarData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return utils.CreateError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	tflog.Debug(ctx, fmt.Sprintf("Response body (unmarshalled): %#v", response))

	d.Set("name", response.Name)
	if properties := response.SidecarProperties; properties != nil {
		d.Set("deployment_method", properties.DeploymentMethod)
		d.Set("activity_log_integration_id", properties.LogIntegrationID)
		d.Set("diagnostic_log_integration_id", properties.DiagnosticLogIntegrationID)
	}
	d.Set("labels", response.Labels)
	d.Set("user_endpoint", response.UserEndpoint)
	if bypassMode := response.BypassMode(); bypassMode != "" {
		d.Set("bypass_mode", bypassMode)
	}
	d.Set("certificate_bundle_secrets", flattenCertificateBundleSecrets(response.CertificateBundleSecrets))

	tflog.Debug(ctx, "End resourceSidecarRead")

	return diag.Diagnostics{}
}

func resourceSidecarUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceSidecarUpdate")
	c := m.(*client.Client)

	resourceData, err := getSidecarDataFromResource(c, d)
	if err != nil {
		return utils.CreateError("Unable to update sidecar", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	if _, err = c.DoRequest(ctx, url, http.MethodPut, resourceData); err != nil {
		return utils.CreateError("Unable to update sidecar", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceSidecarUpdate")

	return resourceSidecarRead(ctx, d, m)
}

func resourceSidecarDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "Init resourceSidecarDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(ctx, url, http.MethodDelete, nil); err != nil {
		return utils.CreateError("Unable to delete sidecar", fmt.Sprintf("%v", err))
	}

	tflog.Debug(ctx, "End resourceSidecarDelete")

	return diag.Diagnostics{}
}

func getSidecarDataFromResource(c *client.Client, d *schema.ResourceData) (*SidecarData, error) {
	ctx := context.Background()
	tflog.Debug(ctx, "Init getSidecarDataFromResource")

	deploymentMethod := d.Get("deployment_method").(string)

	activityLogIntegrationID := d.Get("activity_log_integration_id").(string)
	if activityLogIntegrationID == "" {
		activityLogIntegrationID = d.Get("log_integration_id").(string)
	}
	diagnosticLogIntegrationID := d.Get("diagnostic_log_integration_id").(string)

	properties := NewSidecarProperties(deploymentMethod, activityLogIntegrationID, diagnosticLogIntegrationID)

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

	tflog.Debug(ctx, "end getSidecarDataFromResource")
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
	ctx := context.Background()
	tflog.Debug(ctx, "Init flattenCertificateBundleSecrets")
	var flatCBS []interface{}
	if cbs != nil {
		cb := make(map[string]interface{})

		for key, val := range cbs {
			// Ignore self-signed certificates
			if key != "sidecar-generated-selfsigned" {
				contentCB := make([]interface{}, 1)

				tflog.Debug(ctx, fmt.Sprintf("key: %v", key))
				tflog.Debug(ctx, fmt.Sprintf("val: %v", val))

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

	tflog.Debug(ctx, fmt.Sprintf("end flattenCertificateBundleSecrets %v", flatCBS))
	return flatCBS
}

func getCertificateBundleSecret(d *schema.ResourceData) CertificateBundleSecrets {
	ctx := context.Background()
	tflog.Debug(ctx, "Init getCertificateBundleSecret")
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

	tflog.Debug(ctx, "end getCertificateBundleSecret")
	return ret
}
