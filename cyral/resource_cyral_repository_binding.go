package cyral

import (
	"context"
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
	TCPListeners TCPListener `json:"tcpListeners"`
}

type TCPListener struct {
	Listeners []Listener `json:"listeners"`
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
			"sidecar_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_port": {
				Type:     schema.TypeInt,
				Required: true,
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

	err = c.CreateResource(url, http.MethodPut, resourceData, nil)
	if err != nil {
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

	response := RepoBindingData{
		SidecarID:    sidecarID,
		RepositoryID: repositoryID,
	}
	err := c.ReadResource(url, &response)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read repository. SidecarID: %s, RepositoryID: %s",
			sidecarID, repositoryID), fmt.Sprintf("%v", err))
	}

	d.Set("sidecar_id", response.SidecarID)
	d.Set("repository_id", response.RepositoryID)
	if len(response.TCPListeners.Listeners) > 0 {
		d.Set("listener_port", response.TCPListeners.Listeners[0].Port)
		if host := response.TCPListeners.Listeners[0].Host; host != "" {
			d.Set("listener_host", response.TCPListeners.Listeners[0].Host)
		}
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

	err = c.UpdateResource(resourceData, url)

	if err != nil {
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

	err := c.DeleteResource(url)
	if err != nil {
		return createError("Unable to delete sidecar", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryBindingDelete")

	return diag.Diagnostics{}
}

func getRepoBindingDataFromResource(c *client.Client, d *schema.ResourceData) (RepoBindingData, error) {
	return RepoBindingData{
		SidecarID:    d.Get("sidecar_id").(string),
		RepositoryID: d.Get("repository_id").(string),
		TCPListeners: TCPListener{
			Listeners: []Listener{
				{
					Host: d.Get("listener_host").(string),
					Port: d.Get("listener_port").(int),
				},
			},
		},
	}, nil
}
