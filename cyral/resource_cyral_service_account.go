package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	serviceAccountResourceIDKey           = "id"
	serviceAccountResourceDisplayNameKey  = "display_name"
	serviceAccountResourceClientIDKey     = "client_id"
	serviceAccountResourceClientSecretKey = "client_secret"
	serviceAccountResourcePermissionsKey  = "permissions"
)

var (
	ReadServiceAccountConfig = ResourceOperationConfig{
		Name:       "ServiceAccountRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/users/serviceAccounts/%s",
				c.ControlPlane,
				d.Get(serviceAccountResourceClientIDKey),
			)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData {
			return &ServiceAccount{}
		},
		RequestErrorHandler: &ReadIgnoreHttpNotFound{resName: "Service account"},
	}
)

func resourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Cyral Service Account (A.k.a: " +
			"[Cyral API Access Key](https://cyral.com/docs/api-ref/api-intro/#api-access-key)).",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "ServiceAccountCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts",
						c.ControlPlane,
					)
				},
				NewResourceData: func() ResourceData {
					return &ServiceAccount{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &ServiceAccount{}
				},
			},
			ReadServiceAccountConfig,
		),
		ReadContext: ReadResource(ReadServiceAccountConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "ServiceAccountUpdate",
				HttpMethod: http.MethodPatch,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts/%s",
						c.ControlPlane,
						d.Get(serviceAccountResourceClientIDKey),
					)
				},
				NewResourceData: func() ResourceData {
					return &ServiceAccount{}
				},
			},
			ReadServiceAccountConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "ServiceAccountDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts/%s",
						c.ControlPlane,
						d.Get(serviceAccountResourceClientIDKey),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			serviceAccountResourceDisplayNameKey: {
				Description: "The service account display name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			serviceAccountResourcePermissionsKey: {
				Description: "A block responsible for configuring the service account permissions.",
				Type:        schema.TypeSet,
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: permissionsSchema,
				},
			},
			IDKey: {
				Description: fmt.Sprintf(
					"The resource identifier. It's equal to `%s`.",
					serviceAccountResourceClientIDKey,
				),
				Type:     schema.TypeString,
				Computed: true,
			},
			serviceAccountResourceClientIDKey: {
				Description: "The service account client ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			serviceAccountResourceClientSecretKey: {
				Description: "The service account client secret. **Note**: This resource is not able to recognize " +
					"changes to the client secret after its creation, so keep in mind that if the client secret is " +
					"rotated, the value present in this attribute will be outdated. If you need to rotate the client " +
					"secret it's recommended that you recreate this terraform resource.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
