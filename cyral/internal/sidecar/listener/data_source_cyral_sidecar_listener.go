package listener

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/exp/slices"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/operationtype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/internal/repository"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

const (
	SidecarListenerListKey = "listener_list"
	DSRepoTypeKey          = "repo_type"
)

type ReadDataSourceSidecarListenerAPIResponse struct {
	ListenerConfigs []SidecarListener `json:"listenerConfigs"`
}

func (data ReadDataSourceSidecarListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	ctx := context.Background()
	tflog.Debug(ctx, "Init ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")
	var listenersList []any
	tflog.Debug(ctx, fmt.Sprintf("data.ListenerConfig: %+v", data.ListenerConfigs))
	tflog.Debug(ctx, "Init for _, l := range data.ListenerConfig")
	repoTypeFilter := d.Get(DSRepoTypeKey).(string)
	portFilter := d.Get(utils.PortKey).(int)
	for _, listenerConfig := range data.ListenerConfigs {
		// Check if either the repo filter or the port filter is provided and matches the listener
		if (repoTypeFilter == "" || slices.Contains(listenerConfig.RepoTypes, repoTypeFilter)) &&
			(portFilter == 0 || listenerConfig.NetworkAddress.Port == portFilter) {
			listener := map[string]any{
				utils.ListenerIDKey:  listenerConfig.ListenerId,
				utils.SidecarIDKey:   d.Get(utils.SidecarIDKey).(string),
				RepoTypesKey:         listenerConfig.RepoTypes,
				NetworkAddressKey:    listenerConfig.NetworkAddressAsInterface(),
				MySQLSettingsKey:     listenerConfig.MySQLSettingsAsInterface(),
				S3SettingsKey:        listenerConfig.S3SettingsAsInterface(),
				DynamoDbSettingsKey:  listenerConfig.DynamoDbSettingsAsInterface(),
				SQLServerSettingsKey: listenerConfig.SQLServerSettingsAsInterface(),
			}
			tflog.Debug(ctx, fmt.Sprintf("listener: %q", listener))
			listenersList = append(listenersList, listener)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("listenersList: %q", listenersList))
	tflog.Debug(ctx, "End for _, l := range data.ListenerConfig")

	if err := d.Set(SidecarListenerListKey, listenersList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	tflog.Debug(ctx, "End ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")

	return nil
}

func dataSourceSidecarListenerReadConfig() core.ResourceOperationConfig {
	return core.ResourceOperationConfig{
		ResourceName: "SidecarListenerDataSourceRead",
		Type:         operationtype.Read,
		HttpMethod:   http.MethodGet,
		URLFactory: func(d *schema.ResourceData, c *client.Client) string {
			sidecarID := d.Get(utils.SidecarIDKey).(string)

			return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners", c.ControlPlane, sidecarID)
		},
		SchemaWriterFactory: func(_ *schema.ResourceData) core.SchemaWriter { return &ReadDataSourceSidecarListenerAPIResponse{} },
	}
}

func DataSourceSidecarListener() *schema.Resource {
	listenerSchema := utils.ConvertSchemaFieldsToComputed(getSidecarListenerSchema())
	return &schema.Resource{
		Description: "Retrieve and filter sidecar listeners.",
		ReadContext: core.ReadResource(dataSourceSidecarListenerReadConfig()),
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
