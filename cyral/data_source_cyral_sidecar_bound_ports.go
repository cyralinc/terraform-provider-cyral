package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type BindingConfig struct {
	Enabled                           bool             `json:"enabled,omitempty"`
	Listener                          *WrapperListener `json:"listener,omitempty"`
	AdditionalListeners               []*TCPListener   `json:"additionalListeners,omitempty"`
	TcpListeners                      *TCPListeners    `json:"tcpListeners,omitempty"`
	IsSelectedIdentityProviderSidecar bool             `json:"isSelectedIdentityProviderSidecar,omitempty"`
}

type WrapperListener struct {
	File string `json:"file,omitempty"`
	Host string `json:"host,omitempty"`
	Port uint32 `json:"port,omitempty"`
}

type TCPListener struct {
	Disabled bool   `json:"disabled,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     uint32 `json:"port,omitempty"`
}

type TCPListeners struct {
	Listeners []*TCPListener `json:"listeners,omitempty"`
}

func dataSourceSidecarBoundPorts() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves all the ports of a given sidecar that are currently bound to repositories.",
		ReadContext: dataSourceSidecarBoundPortsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Computed ID for this data source (locally computed to be used in Terraform state).",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"sidecar_id": {
				Description: "The ID of the sidecar.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"bound_ports": {
				Description: "All the sidecar ports that are currently bound to repositories.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func dataSourceSidecarBoundPortsRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceSidecarBoundPortsRead")
	c := m.(*client.Client)

	var boundPorts []uint32

	sidecarID := d.Get("sidecar_id").(string)
	repoIDs, err := getRepoIDsBoundToSidecar(c, sidecarID)
	if err != nil {
		return createError(fmt.Sprintf("Unable to retrieve repo IDs bound to sidecar. SidecarID: %s",
			sidecarID), err.Error())
	}

	for _, repoID := range repoIDs {
		bindingConfig, err := getRepoBinding(c, sidecarID, repoID)
		if err != nil {
			return createError(fmt.Sprintf("Unable to retrieve repository binding. SidecarID: %s, "+
				"RepositoryID: %s", sidecarID, repoID), err.Error())
		}
		repoInfo, err := getRepoInfo(c, repoID)
		if err != nil {
			return createError(fmt.Sprintf("Unable to retrieve repository info. RepositoryID: %s",
				repoID), err.Error())
		}
		bindingPorts := getBindingPorts(bindingConfig, repoInfo)
		boundPorts = append(boundPorts, bindingPorts...)
	}
	// Sorts ports so that we can have a more organized and deterministic output.
	sort.Slice(boundPorts, func(i, j int) bool { return boundPorts[i] < boundPorts[j] })

	d.SetId(uuid.New().String())
	d.Set("bound_ports", boundPorts)

	log.Printf("[DEBUG] Sidecar bound ports: %v", boundPorts)
	log.Printf("[DEBUG] End dataSourceSidecarBoundPortsRead")

	return diag.Diagnostics{}
}

func getRepoIDsBoundToSidecar(c *client.Client, sidecarID string) ([]string, error) {
	log.Printf("[DEBUG] Init getRepoIDsBoundToSidecar")
	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos", c.ControlPlane, sidecarID)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var repoIDs []string
	if err := json.Unmarshal(body, &repoIDs); err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", repoIDs)
	log.Printf("[DEBUG] End getRepoIDsBoundToSidecar")

	return repoIDs, nil
}

func getRepoBinding(
	c *client.Client,
	sidecarID, repoID string,
) (BindingConfig, error) {
	log.Printf("[DEBUG] Init getRepoBinding")
	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane, sidecarID, repoID)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return BindingConfig{}, err
	}

	var bindingConfig BindingConfig
	if err := json.Unmarshal(body, &bindingConfig); err != nil {
		return BindingConfig{}, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", bindingConfig)
	log.Printf("[DEBUG] End getRepoBinding")

	return bindingConfig, nil
}

func getRepoInfo(c *client.Client, repoID string) (RepoData, error) {
	log.Printf("[DEBUG] Init getRepoInfo")
	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, repoID)
	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return RepoData{}, err
	}

	var repoResponse GetRepoByIDResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return RepoData{}, err
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", repoResponse)
	log.Printf("[DEBUG] End getRepoInfo")

	return repoResponse.Repo, nil
}

// getBindingPorts retrieves all the ports of a repo binding, based on the binding
// config and the repo info attributes. This can be useful for repositories that have
// more than one port bound to a sidecar, such as mongodb replicasets that have more
// than one node, or S3 repos that support S3 Browser, etc.
func getBindingPorts(
	bindingConfig BindingConfig,
	repoInfo RepoData,
) []uint32 {
	var bindingPorts []uint32

	// primaryPort is the main port/listener that the sidecar will use for this binding
	var primaryPort uint32
	if bindingConfig.Listener != nil {
		primaryPort = bindingConfig.Listener.Port
		bindingPorts = append(bindingPorts, primaryPort)
	}

	// MaxAllowedListeners is currently used to generate a repo binding port range
	// between primaryPort and (primaryPort+MaxAllowedListeners-1). Thats why we
	// add (MaxAllowedListeners-1) more ports to the repo binding ports.
	for i := uint32(1); i < repoInfo.MaxAllowedListeners; i++ {
		bindingPorts = append(bindingPorts, primaryPort+i)
	}

	// TcpListeners, also referred as seedListeners, are used as extra data
	// repo out routes, and can be used, for instance, for mongoDB replicasets
	// and sharded clusters.
	if bindingConfig.TcpListeners != nil {
		for _, seedListener := range bindingConfig.TcpListeners.Listeners {
			if seedListener != nil {
				bindingPorts = append(bindingPorts, seedListener.Port)
			}
		}
	}

	// AdditionalListeners is currently used to represent the additional ports of
	// a repo binding, like S3 browser ports for S3 repositories.
	for _, additionalListener := range bindingConfig.AdditionalListeners {
		if additionalListener != nil {
			bindingPorts = append(bindingPorts, additionalListener.Port)
		}
	}

	return bindingPorts
}
