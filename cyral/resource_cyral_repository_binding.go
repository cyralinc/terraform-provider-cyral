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
	SidecarID                 string
	RepositoryID              string
	Enabled                   bool
	SidecarAsIdPAccessGateway bool     `json:"isSelectedIdentityProviderSidecar,omitempty"`
	Listener                  Listener `json:"listener"`
}

type Listener struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func resourceRepositoryBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [repositories to sidecars binding](https://cyral.com/docs/sidecars/sidecar-assign-repo).",
		CreateContext: resourceRepositoryBindingCreate,
		ReadContext:   resourceRepositoryBindingRead,
		UpdateContext: resourceRepositoryBindingUpdate,
		DeleteContext: resourceRepositoryBindingDelete,

		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Enable|Disable the repository in the target sidecar. It is important to notice that the resource will always be created, but will remain inactive if set to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"sidecar_id": {
				Description: "ID of the sidecar that the repository(ies) will be bound to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"repository_id": {
				Description: "ID of the repository that will be bound to the sidecar.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"listener_port": {
				Description: "Port in which the sidecar will listen for the given repository.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"listener_host": {
				Description: "Address in which the sidecar will listen for the given repository. By default, the sidecar will listen in all interfaces.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.0.0",
			},
			"sidecar_as_idp_access_gateway": {
				Description: "Indicates whether or not the sidecar in the binding configuration is selected as the Access Gateway for Identity Provider users connecting to the underlying data repository. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"id": {
				Description: "Computed ID for this resource (locally computed to be used in Terraform state).",
				Computed:    true,
				Type:        schema.TypeString,
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

	resourceData := getRepoBindingDataFromResource(d)

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane,
		resourceData.SidecarID, resourceData.RepositoryID)

	if _, err := c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to bind repository to sidecar", fmt.Sprintf("%v", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", resourceData.SidecarID, resourceData.RepositoryID))

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

	response := RepoBindingData{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON. SidecarID: %s, RepositoryID: %s",
			sidecarID, repositoryID), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.Set("enabled", response.Enabled)
	d.Set("sidecar_as_idp_access_gateway", response.SidecarAsIdPAccessGateway)
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

	resourceData := getRepoBindingDataFromResource(d)

	if err := updateRepositoryBinding(c, resourceData); err != nil {
		return createError("Unable to update repository binding", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryBindingUpdate")

	return resourceRepositoryBindingRead(ctx, d, m)
}

func resourceRepositoryBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryBindingDelete")
	c := m.(*client.Client)

	// SidecarAsIdPAccessGateway is set to false to stop
	// using the bound sidecar as the Access Gateway for Identity
	// Provider users. This is needed so that the binding can
	// be deleted, otherwise it will throw a validation error.
	resourceData := getRepoBindingDataFromResource(d)
	resourceData.SidecarAsIdPAccessGateway = false
	if err := updateRepositoryBinding(c, resourceData); err != nil {
		return createError("Unable to delete repository binding",
			fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane,
		resourceData.SidecarID, resourceData.RepositoryID)

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete repository binding", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryBindingDelete")

	return diag.Diagnostics{}
}

func getRepoBindingDataFromResource(d *schema.ResourceData) RepoBindingData {
	return RepoBindingData{
		Enabled:                   d.Get("enabled").(bool),
		SidecarID:                 d.Get("sidecar_id").(string),
		RepositoryID:              d.Get("repository_id").(string),
		SidecarAsIdPAccessGateway: d.Get("sidecar_as_idp_access_gateway").(bool),
		Listener: Listener{
			Host: d.Get("listener_host").(string),
			Port: d.Get("listener_port").(int),
		},
	}
}

func updateRepositoryBinding(c *client.Client, resourceData RepoBindingData) error {
	url := fmt.Sprintf("https://%s/v1/sidecars/%s/repos/%s", c.ControlPlane,
		resourceData.SidecarID, resourceData.RepositoryID)
	_, err := c.DoRequest(url, http.MethodPut, resourceData)
	return err
}
