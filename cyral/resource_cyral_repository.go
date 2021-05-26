package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

var supportedTypes = map[string]struct{}{
	`bigquery`:   {},
	`cassandra`:  {},
	`dremio`:     {},
	`galera`:     {},
	`mariadb`:    {},
	`mongodb`:    {},
	`mysql`:      {},
	`oracle`:     {},
	`postgresql`: {},
	`s3`:         {},
	`snowflake`:  {},
	`sqlserver`:  {},
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Description:   "CRUD operations for Cyral datamaps",
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository type (see the list of supported types below)",
				ValidateDiagFunc: func(i interface{}, p cty.Path) diag.Diagnostics {
					t := i.(string)
					if _, ok := supportedTypes[t]; !ok {
						return diag.Errorf("repository type %v not supported. supportedTypes: %v", t, supportedTypes)
					}
					return nil
				},
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository host name (ex: `somerepo.cyral.com`)",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Repository access port (ex: `3306`)",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Repository name that will be used internally in Control Plane (ex: `your_repo_name`)",
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

	body, err := c.DoRequest(url, http.MethodPost, resourceData)
	if err != nil {
		return createError("Unable to create repository", fmt.Sprintf("%v", err))
	}

	response := IDBasedResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError("Unable to unmarshall JSON", fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

	d.SetId(response.ID)

	log.Printf("[DEBUG] End resourceRepositoryCreate")

	return resourceRepositoryRead(ctx, d, m)
}

func resourceRepositoryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryRead")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())

	body, err := c.DoRequest(url, http.MethodGet, nil)
	if err != nil {
		return createError(fmt.Sprintf("Unable to read repository. RepositoryID: %s",
			d.Id()), fmt.Sprintf("%v", err))
	}

	response := GetRepoByIDResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return createError(fmt.Sprintf("Unable to unmarshall JSON"), fmt.Sprintf("%v", err))
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", response)

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

	if _, err = c.DoRequest(url, http.MethodPut, resourceData); err != nil {
		return createError("Unable to update repository", fmt.Sprintf("%v", err))
	}

	log.Printf("[DEBUG] End resourceRepositoryUpdate")

	return resourceRepositoryRead(ctx, d, m)
}

func resourceRepositoryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Init resourceRepositoryDelete")
	c := m.(*client.Client)

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, d.Id())

	if _, err := c.DoRequest(url, http.MethodDelete, nil); err != nil {
		return createError("Unable to delete repository", fmt.Sprintf("%v", err))
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
