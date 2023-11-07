package serviceaccount

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Schema keys
	ServiceAccountResourceDisplayNameKey   = "display_name"
	ServiceAccountResourcePermissionIDsKey = "permission_ids"
	ServiceAccountResourceClientIDKey      = "client_id"
	ServiceAccountResourceClientSecretKey  = "client_secret"
)

var (
	ReadServiceAccountConfig = core.ResourceOperationConfig{
		Name:       "ServiceAccountRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			return fmt.Sprintf(
				"https://%s/v1/users/serviceAccounts/%s",
				c.ControlPlane,
				d.Id(),
			)
		},
		NewResponseData: func(_ *schema.ResourceData) core.SchemaWriter {
			return &ServiceAccount{}
		},
		RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Service account"},
	}
)

func ResourceServiceAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Cyral Service Account (A.k.a: " +
			"[Cyral API Access Key](https://cyral.com/docs/api-ref/api-intro/#api-access-key)). See also " +
			"data source [`cyral_permission`](../data-sources/permission.md)." +
			"\n\n-> **Note** This resource does not support importing, since the client secret cannot " +
			"be read after the resource creation.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "ServiceAccountCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts",
						c.ControlPlane,
					)
				},
				NewResourceData: func() core.SchemaReader { return &ServiceAccount{} },
				NewResponseData: func(_ *schema.ResourceData) core.SchemaWriter { return &ServiceAccount{} },
			},
			ReadServiceAccountConfig,
		),
		ReadContext: core.ReadResource(ReadServiceAccountConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "ServiceAccountUpdate",
				HttpMethod: http.MethodPatch,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
				NewResourceData: func() core.SchemaReader { return &ServiceAccount{} },
			},
			ReadServiceAccountConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
				Name:       "ServiceAccountDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/users/serviceAccounts/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			ServiceAccountResourceDisplayNameKey: {
				Description: "The service account display name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			ServiceAccountResourcePermissionIDsKey: {
				Description: "A list of permission IDs that will be assigned to this service account. See " +
					"also data source [`cyral_permission`](../data-sources/permission.md).",
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			utils.IDKey: {
				Description: fmt.Sprintf(
					"The resource identifier. It's equal to `%s`.",
					ServiceAccountResourceClientIDKey,
				),
				Type:     schema.TypeString,
				Computed: true,
			},
			ServiceAccountResourceClientIDKey: {
				Description: "The service account client ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			ServiceAccountResourceClientSecretKey: {
				Description: "The service account client secret. **Note**: This resource is not able to recognize " +
					"changes to the client secret after its creation, so keep in mind that if the client secret is " +
					"rotated, the value present in this attribute will be outdated. If you need to rotate the client " +
					"secret it's recommended that you recreate this terraform resource.",
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
