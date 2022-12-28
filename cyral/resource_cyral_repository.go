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

const (
	mongodbRepoType             = "mongodb"
	mongodbReplicaSetServerType = "replicaset"
)

func repositoryTypes() []string {
	return []string{
		"denodo",
		"dremio",
		"dynamodb",
		"dynamodbstreams",
		"galera",
		"mariadb",
		"mongodb",
		"mysql",
		"oracle",
		"postgresql",
		"redshift",
		"s3",
		"snowflake",
		"sqlserver",
	}
}

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

func (data *RepoData) WriteToSchema(d *schema.ResourceData) {
	d.Set("type", data.RepoType)
	d.Set("host", data.Host)
	d.Set("port", data.Port)
	d.Set("name", data.Name)
	d.Set("labels", data.Labels)

	if properties := data.PropertiesAsInterface(); properties != nil {
		d.Set("properties", properties)
	}
}

func (data *RepoData) PropertiesAsInterface() []interface{} {
	var properties []interface{}
	if data.Properties != nil {
		if data.IsReplicaSet() {
			propertiesMap := make(map[string]interface{})
			var rset []interface{}
			rsetMap := make(map[string]interface{})
			rsetMap["max_nodes"] = data.MaxAllowedListeners
			rsetMap["replica_set_id"] = data.Properties.MongoDBReplicaSetName
			rset = append(rset, rsetMap)

			propertiesMap["mongodb_replica_set"] = rset
			properties = append(properties, propertiesMap)
		}
	}

	return properties
}

func (data *RepoData) IsReplicaSet() bool {
	return data.Properties != nil && data.Properties.MongoDBServerType == mongodbReplicaSetServerType
}

// RepositoryProperties relates to the field "properties" of the v1/repos
// API. All fields of this struct _must_ be of type string, to comply with the
// API.
type RepositoryProperties struct {
	// Replica set
	MongoDBReplicaSetName string `json:"mongodb-replicaset-name,omitempty"`
	MongoDBServerType     string `json:"mongodb-server-type,omitempty"`
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
				Description:  "Repository type. List of supported types:" + supportedTypesMarkdown(repositoryTypes()),
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(repositoryTypes(), false),
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
			"properties": {
				Description: "Contains advanced repository configuration.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mongodb_replica_set": {
							Description: "Used to configure a MongoDB cluster.",
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_nodes": {
										Description:  "Maximum number of nodes of the replica set cluster.",
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
									"replica_set_id": {
										Description:  "Identifier of the replica set cluster. Used to construct the URI command (available in Cyral's Access Portal page) that your users will need for connecting to the repository via Cyral.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

	response.Repo.WriteToSchema(d)

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
	repoData := RepoData{
		ID:       d.Id(),
		RepoType: d.Get("type").(string),
		Host:     d.Get("host").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
	}

	labels := d.Get("labels").([]interface{})
	repositoryDataLabels := make([]string, len(labels))
	for i, label := range labels {
		repositoryDataLabels[i] = (label).(string)
	}
	repoData.Labels = repositoryDataLabels

	var maxAllowedListeners uint32
	var properties *RepositoryProperties
	if propertiesIface, ok := d.Get("properties").(*schema.Set); ok {
		for _, propertiesMap := range propertiesIface.List() {
			properties = new(RepositoryProperties)
			propertiesMap := propertiesMap.(map[string]interface{})

			// Replica set properties
			if rsetIface, ok := propertiesMap["mongodb_replica_set"]; ok {
				if repoData.RepoType != mongodbRepoType {
					return RepoData{}, fmt.Errorf(
						"replica sets are only supported for repository type '%s'",
						mongodbRepoType)
				}

				for _, rsetMap := range rsetIface.(*schema.Set).List() {
					rsetMap := rsetMap.(map[string]interface{})
					maxAllowedListeners = uint32(rsetMap["max_nodes"].(int))
					properties.MongoDBReplicaSetName = rsetMap["replica_set_id"].(string)
				}
				properties.MongoDBServerType = mongodbReplicaSetServerType
			}
		}
	}
	repoData.MaxAllowedListeners = maxAllowedListeners
	repoData.Properties = properties

	return repoData, nil
}
