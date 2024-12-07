package credentials

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                  resourceName,
	ResourceType:                  resourcetype.Resource,
	SchemaReaderFactory:           func() core.SchemaReader { return &CreateSidecarCredentialsRequest{} },
	SchemaWriterFactoryGetMethod:  func(_ *schema.ResourceData) core.SchemaWriter { return &SidecarCredentialsData{} },
	SchemaWriterFactoryPostMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &SidecarCredentialsData{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/users/sidecarAccounts", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Create new [credentials for Cyral sidecar](https://cyral.com/docs/sidecars/manage/#rotate-the-client-secret-for-a-sidecar).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Same as `client_id`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sidecar_id": {
				Description: "ID of the sidecar to create new credentials.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Description: "Sidecar client ID.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"client_secret": {
				Description: "Sidecar client secret.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
