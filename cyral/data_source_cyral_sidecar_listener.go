package cyral

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"golang.org/x/exp/slices"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const (
	SidecarListenerListKey = "listener_list"
	DSRepoTypeKey          = "repo_type"
)

type ReadDataSourceSidecarListenerAPIResponse struct {
	ListenerConfigs []SidecarListener `json:"listenerConfigs"`
}

func (data ReadDataSourceSidecarListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	log.Printf("[DEBUG] Init ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")
	var listenersList []any
	log.Printf("[DEBUG] data.ListenerConfig: %+v", data.ListenerConfigs)
	log.Printf("[DEBUG] Init for _, l := range data.ListenerConfig")
	repoTypeFilter := d.Get(DSRepoTypeKey).(string)
	portFilter := d.Get(PortKey).(int)
	for _, listenerConfig := range data.ListenerConfigs {
		// Check if either the repo filter or the port filter is provided and matches the listener
		if (repoTypeFilter == "" || slices.Contains(listenerConfig.RepoTypes, repoTypeFilter)) &&
			(portFilter == 0 || listenerConfig.NetworkAddress.Port == portFilter) {
			listener := map[string]any{
				ListenerIDKey:        listenerConfig.ListenerId,
				SidecarIDKey:         d.Get(SidecarIDKey).(string),
				RepoTypesKey:         listenerConfig.RepoTypes,
				NetworkAddressKey:    listenerConfig.NetworkAddressAsInterface(),
				MySQLSettingsKey:     listenerConfig.MySQLSettingsAsInterface(),
				S3SettingsKey:        listenerConfig.S3SettingsAsInterface(),
				DynamoDbSettingsKey:  listenerConfig.DynamoDbSettingsAsInterface(),
				SQLServerSettingsKey: listenerConfig.SQLServerSettingsAsInterface(),
			}
			log.Printf("[DEBUG] listener: %q", listener)
			listenersList = append(listenersList, listener)
		}
	}

	log.Printf("[DEBUG] listenersList: %q", listenersList)
	log.Printf("[DEBUG] End for _, l := range data.ListenerConfig")

	if err := d.Set(SidecarListenerListKey, listenersList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	log.Printf("[DEBUG] End ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")

	return nil
}

func dataSourceSidecarListenerReadConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "SidecarListenerDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			sidecarID := d.Get(SidecarIDKey).(string)

			return fmt.Sprintf("https://%s/v1/sidecars/%s/listeners", c.ControlPlane, sidecarID)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &ReadDataSourceSidecarListenerAPIResponse{} },
	}
}

func dataSourceSidecarListener() *schema.Resource {
	log.Printf("[DEBUG] Init dataSourceSidecarListener")
	listenerSchema := convertSchemaFieldsToComputed(getSidecarListenerSchema())

	log.Printf("[DEBUG] End dataSourceSidecarListener")
	return &schema.Resource{
		Description: "Retrieve and filter sidecar listeners.",
		ReadContext: ReadResource(dataSourceSidecarListenerReadConfig()),
		Schema: map[string]*schema.Schema{
			SidecarIDKey: {
				Description: "Filter the results by sidecar ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			DSRepoTypeKey: {
				Description: "Filter the results per repository type. Supported repo types:" + supportedTypesMarkdown(repositoryTypes()),
				Type:        schema.TypeString,
				Optional:    true,
			},
			PortKey: {
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
