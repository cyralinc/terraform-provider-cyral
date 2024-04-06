package serviceaccount

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &ServiceAccount{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &ServiceAccount{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/serviceAccounts", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Cyral Service Account (A.k.a: " +
			"[Cyral API Access Key](https://cyral.com/docs/api-ref/api-intro/#api-access-key)). See also " +
			"data source [`cyral_permission`](../data-sources/permission.md)." +
			"\n\n-> **Note** This resource does not support importing, since the client secret cannot " +
			"be read after the resource creation.",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),

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
