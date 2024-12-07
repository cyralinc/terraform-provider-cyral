package teams

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceContextHandler = core.HTTPContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &MsTeamsIntegration{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &MsTeamsIntegration{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/integrations/notifications/teams", c.ControlPlane)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [integration with Microsoft Teams](https://cyral.com/docs/integrations/alerting/microsoft-teams).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Integration name that will be used internally in the control plane.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "Microsoft Teams webhook URL.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}
