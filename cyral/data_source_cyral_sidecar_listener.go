package cyral

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"golang.org/x/exp/slices"
)

const (
	SidecarListenerListKey = "listener_list"
)

type ReadDataSourceSidecarListenerAPIResponse struct {
	ListenerConfig []SidecarListener `json:"listenerConfigs"`
}

func (data ReadDataSourceSidecarListenerAPIResponse) WriteToSchema(d *schema.ResourceData) error {
	log.Printf("[DEBUG] Init ReadDataSourceSidecarListenerAPIResponse.WriteToSchema")
	var listenersList []interface{}
	log.Printf("[DEBUG] data.ListenerConfig: %+v", data.ListenerConfig)
	log.Printf("[DEBUG] Init for _, l := range data.ListenerConfig")
	repoTypeFilter := d.Get(RepoTypesKey).(string)
	portFilter := d.Get(PortKey).(int)
	for _, l := range data.ListenerConfig {
		// Check if either the repo filter or the port filter is provided and matches the listener
		if (repoTypeFilter == "" || slices.Contains(l.RepoTypes, repoTypeFilter)) &&
			(portFilter == 0 || l.NetworkAddress.Port == portFilter) {
			argumentVals := map[string]interface{}{
				ListenerIDKey:       l.ListenerId,
				SidecarIDKey:        d.Get(SidecarIDKey).(string),
				RepoTypesKey:        l.RepoTypes,
				NetworkAddressKey:   l.NetworkAddressAsInterface(),
				MySQLSettingsKey:    l.MySQLSettingsAsInterface(),
				S3SettingsKey:       l.S3SettingsAsInterface(),
				DynamoDbSettingsKey: l.DynamoDbSettingsAsInterface(),
			}
			log.Printf("[DEBUG] argumentVals: %q", argumentVals)
			listenersList = append(listenersList, argumentVals)
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
	listenerSchema := getSidecarListenerSchema()
	for _, v := range listenerSchema {
		v.Required = false
		v.Computed = true
		v.MaxItems = 0
		v.ConflictsWith = nil
	}

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
			RepoTypesKey: {
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
