package sidecar

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
)

var urlFactory = func(d *schema.ResourceData, c *client.Client) string {
	return fmt.Sprintf("https://%s/v1/sidecars/%s", c.ControlPlane, d.Id())
}

var readConfig = core.ResourceOperationConfig{
	ResourceName:        resourceName,
	ResourceType:        resourcetype.Resource,
	Type:                operationtype.Read,
	HttpMethod:          http.MethodGet,
	URLFactory:          urlFactory,
	SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &SidecarData{} },
	RequestErrorHandler: &core.IgnoreNotFoundByMessage{
		ResName:        resourceName,
		MessageMatches: "NotFound",
		OperationType:  operationtype.Read,
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [sidecars](https://cyral.com/docs/sidecars/manage).",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				ResourceName: resourceName,
				ResourceType: resourcetype.Resource,
				Type:         operationtype.Create,
				HttpMethod:   http.MethodPost,
				URLFactory: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf("https://%s/v1/sidecars", c.ControlPlane)
				},
				SchemaReaderFactory: func() core.SchemaReader { return &SidecarData{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &core.IDBasedResponse{} },
			},
			readConfig,
		),
		ReadContext: core.ReadResource(readConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				ResourceName:        resourceName,
				ResourceType:        resourcetype.Resource,
				Type:                operationtype.Update,
				HttpMethod:          http.MethodPut,
				URLFactory:          urlFactory,
				SchemaReaderFactory: func() core.SchemaReader { return &SidecarData{} },
				SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &SidecarData{} },
			},
			readConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				ResourceName: resourceName,
				ResourceType: resourcetype.Resource,
				Type:         operationtype.Delete,
				HttpMethod:   http.MethodDelete,
				URLFactory:   urlFactory,
				RequestErrorHandler: &core.IgnoreNotFoundByMessage{
					ResName:        resourceName,
					MessageMatches: "NotFound",
					OperationType:  operationtype.Delete,
				},
			},
		),
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
					" to [our public docs](https://cyral.com/docs/sidecars/deployment/certificates)" +
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
