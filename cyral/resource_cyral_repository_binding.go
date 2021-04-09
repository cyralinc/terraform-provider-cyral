package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RepoBindingData struct {
	SidecarID    string
	RepositoryID string
	Enabled      bool
	Listener     Listener `json:"listener"`
}

type Listener struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func resourceRepositoryBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryBindingCreate,
		ReadContext:   resourceRepositoryBindingRead,
		UpdateContext: resourceRepositoryBindingUpdate,
		DeleteContext: resourceRepositoryBindingDelete,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_port": {
				Type:      schema.TypeInt,
				Required:  true,
				Sensitive: true,
			},
			"listener_host": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRepositoryBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryBindingCreate")
	c := m.(*client.Client)

	resourceData, err := getRepoBindingDataFromResource(c, d)
	if err != nil {
		return createError("Unable to bind repository to sidecar", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane,
		resourceData.SidecarID, resourceData.RepositoryID)

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to bind repository to sidecar", fmt.Sprintf("%v", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", resourceData.SidecarID, resourceData.RepositoryID))
	d.Set("sidecar_id", resourceData.SidecarID)
	d.Set("repository_id", resourceData.RepositoryID)

	return resourceRepositoryBindingRead(ctx, d, m)
}

func resourceRepositoryBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryBindingRead")
	c := m.(*client.Client)

	sidecarID := d.Get("sidecar_id").(string)
	repositoryID := d.Get("repository_id").(string)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane, sidecarID, repositoryID)

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read repository. SidecarID: %s, RepositoryID: %s",
			sidecarID, repositoryID), fmt.Sprintf("%v", err))
	}

	response := RepoBindingData{
		SidecarID:    sidecarID,
		RepositoryID: repositoryID,
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON. SidecarID: %s, RepositoryID: %s",
			sidecarID, repositoryID), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("enabled", response.Enabled)
	d.Set("sidecar_id", response.SidecarID)
	d.Set("repository_id", response.RepositoryID)
	d.Set("listener_port", response.Listener.Port)
	if host := response.Listener.Host; host != "" {
		d.Set("listener_host", response.Listener.Host)
	}
	log.Printf("[DEBUG] End resourceRepositoryBindingRead")

	return diag.Diagnostics{}
}

func resourceRepositoryBindingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryBindingUpdate")
	c := m.(*client.Client)

	resourceData, err := getRepoBindingDataFromResource(c, d)
	if err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane,
		resourceData.SidecarID, resourceData.RepositoryID)

	if _, err = c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryBindingUpdate")

	return resourceRepositoryBindingRead(ctx, d, m)
}

func resourceRepositoryBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryBindingDelete")
	c := m.(*client.Client)

	sidecarID := d.Get("sidecar_id").(string)
	repositoryID := d.Get("repository_id").(string)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane, sidecarID, repositoryID)

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete repository binding", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryBindingDelete")

	return diag.Diagnostics{}
}

func getRepoBindingDataFromResource(c *client.Client, d *schema.ResourceData) (RepoBindingData, error) {
	return RepoBindingData{
		Enabled:      d.Get("enabled").(bool),
		SidecarID:    d.Get("sidecar_id").(string),
		RepositoryID: d.Get("repository_id").(string),
		Listener: Listener{
			Host: d.Get("listener_host").(string),
			Port: d.Get("listener_port").(int),
		},
	}, nil
}
