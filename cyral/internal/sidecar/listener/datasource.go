package listener

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

var dsContextHandler = core.DefaultContextHandler{
	ResourceName:                 dataSourceName,
	ResourceType:                 resourcetype.DataSource,
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &ReadDataSourceSidecarListenerAPIResponse{} },
	GetPutDeleteURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners", c.ControlPlane, d.Get(utils.SidecarIDKey).(string))
	},
}

func dataSourceSchema() *schema.Resource {
	listenerSchema := utils.ConvertSchemaFieldsToComputed(getSidecarListenerSchema())
	return &schema.Resource{
		Description: "Retrieve and filter sidecar listeners.",
		ReadContext: dsContextHandler.ReadContext(),
		Schema: map[string]*schema.Schema{
			utils.SidecarIDKey: {
				Description: "Filter the results by sidecar ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			DSRepoTypeKey: {
				Description: "Filter the results per repository type. Supported repo types:" + utils.SupportedValuesAsMarkdown(repository.RepositoryTypes()),
				Type:        schema.TypeString,
				Optional:    true,
			},
			utils.PortKey: {
				Description: "Filter the results per port.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			SidecarListenerListKey: {
				Description: "List of existing listeners satisfying the filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: listenerSchema,
				},
			},
		},
	}
}
