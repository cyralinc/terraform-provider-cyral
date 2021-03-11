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

type CreateRepoResponse struct {
	ID string `json:"ID"`
}

type GetRepoByIDResponse struct {
	Repo RepoData `json:"repo"`
}

type RepoData struct {
	ID       string `json:"id"`
	RepoType string `json:"type"`
	Name     string `json:"name"`
	Host     string `json:"repoHost"`
	Port     int    `json:"repoPort"`
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryCreate")
	c := m.(*client.Client)

	resourceData, err := getRepoDataFromResource(c, d)
	if err != nil {
		return createError("Unable to create repository", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos", c.ControlPlane)

	response := CreateRepoResponse{}
	if err = c.CreateResource(url, http.MethodPost, resourceData, &response); err != nil {
		return createError("Unable to create repository", fmt.Sprintf("%v", err))
	}

	d.SetId(response.ID)

	return resourceRepositoryRead(ctx, d, m)
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())

	response := GetRepoByIDResponse{}
	if err := c.ReadResource(url, &response); err != nil {
		return createError(fmt.Sprintf("Unable to read repository. RepositoryID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	d.Set("type", response.Repo.RepoType)
	d.Set("host", response.Repo.Host)
	d.Set("port", response.Repo.Port)
	d.Set("name", response.Repo.Name)

	log.Printf("[DEBUG] End resourceRepositoryRead")

	return diag.Diagnostics{}
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryUpdate")
	c := m.(*client.Client)

	resourceData, err := getRepoDataFromResource(c, d)
	if err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())
	if err = c.UpdateResource(resourceData, url); err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryUpdate")

	return resourceRepositoryRead(ctx, d, m)
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())
	if err := c.DeleteResource(url); err != nil {
		return createError("Unable to delete sidecar", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryDelete")

	return diag.Diagnostics{}
}

func getRepoDataFromResource(c *client.Client, d *schema.ResourceData) (RepoData, error) {
	repoType := d.Get("type").(string)

	if err := client.ValidateRepoType(repoType); err != nil {
		return RepoData{}, err
	}

	return RepoData{
		ID:       d.Id(),
		RepoType: repoType,
		Host:     d.Get("host").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
	}, nil
}
