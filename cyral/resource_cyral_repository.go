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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type GetRepoByIDResponse struct {
	Repo RepoData `json:"repo"`
}

type RepoData struct {
	ID                  string                `json:"id"`
	RepoType            string                `json:"type"`
	Name                string                `json:"name"`
	Host                string                `json:"repoHost"`
	Port                int                   `json:"repoPort"`
	Labels              []string              `json:"labels"`
	MaxAllowedListeners uint32                `json:"maxAllowedListeners,omitempty"`
	Properties          *RepositoryProperties `json:"properties,omitempty"`
}

type RepositoryProperties struct {
	// Replica set
	MaxNodes     string `json:"max-nodes,omitempty"`
	ReplicaSetID string `json:"mongodb-replicaset-name,omitempty"`
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [repositories](https://cyral.com/docs/manage-repositories/repo-track)." +
			"\n\nSee also [Cyral Repository Configuration Module](https://github.com/cyralinc/terraform-cyral-repository-config)." +
			"\nThis module provides the repository configuration options as shown in Cyral UI.",
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "Repository type. List of supported types:" +
					"\n  - `bigquery`" +
					"\n  - `cassandra`" +
					"\n  - `denodo`" +
					"\n  - `dremio`" +
					"\n  - `galera`" +
					"\n  - `mariadb`" +
					"\n  - `mongodb`" +
					"\n  - `mysql`" +
					"\n  - `oracle`" +
					"\n  - `postgresql`" +
					"\n  - `redshift`" +
					"\n  - `s3`" +
					"\n  - `snowflake`" +
					"\n  - `sqlserver`",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bigquery",
					"cassandra",
					"denodo",
					"dremio",
					"galera",
					"mariadb",
					"mongodb",
					"mysql",
					"oracle",
					"postgresql",
					"redshift",
					"snowflake",
					"s3",
					"sqlserver",
				}, false),
			},
			"host": {
				Description: "Repository host name (ex: `somerepo.cyral.com`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"port": {
				Description: "Repository access port (ex: `3306`).",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": {
				Description: "Repository name that will be used internally in the control plane (ex: `your_repo_name`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"labels": {
				Description: "Labels enable you to categorize your repository.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// 'advanced' is equivalent to 'properties' in the Cyral
			// v1/repos API, but with a user-friendly name.
			"advanced": {
				Description: "Contains advanced repository configuration.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replica_set": {
							Description: "Used to configure a distributed database, such as a MongoDB cluster.",
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_nodes": {
										Description: "Maximum number of nodes of the replica set cluster.",
										Type:        schema.TypeInt,
										Required:    true,
									},
									"replica_set_id": {
										Description: "Identifier of the replica set cluster. Used to construct the URI command (available in Cyral's Access Token page) that your users will need for connecting to the repository via Cyral.",
										Type:        schema.TypeString,
										Required:    true,
									},
								},
							},
						},
					},
				},
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
	d.Set("labels", response.Repo.Labels)
	d.Set("properties", response.Repo.Properties)

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
	labels := d.Get("labels").([]interface{})
	repositoryDataLabels := make([]string, len(labels))
	for i, label := range labels {
		repositoryDataLabels[i] = (label).(string)
	}

	properties := new(RepositoryProperties)
	for _, rsetIface := range d.Get("replica_set").(*schema.Set).List() {
		rsetMap := rsetIface.(map[string]interface{})
		properties.MaxNodes = rsetMap["max_nodes"].(string)
		properties.ReplicaSetID = rsetMap["replica_set_id"].(string)
	}

	return RepoData{
		ID:         d.Id(),
		RepoType:   d.Get("type").(string),
		Host:       d.Get("host").(string),
		Name:       d.Get("name").(string),
		Port:       d.Get("port").(int),
		Labels:     repositoryDataLabels,
		Properties: properties,
	}, nil
}
