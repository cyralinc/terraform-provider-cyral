package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Schema keys.
	RepoIDKey                   = "id"
	RepoTypeKey                 = "type"
	RepoNameKey                 = "name"
	RepoLabelsKey               = "labels"
	RepoConnDrainingKey         = "connection_draining"
	RepoConnDrainingAutoKey     = "auto"
	RepoConnDrainingWaitTimeKey = "wait_time"
	RepoHostKey                 = "host"
	RepoPortKey                 = "port"

	// Deprecated schema keys.
	RepoPropertiesKey        = "properties"
	RepoMongoDBReplicaSetKey = "mongodb_replica_set"
	RepoMaxNodesKey          = "max_nodes"
	RepoReplicaSetIDKey      = "replica_set_id"

	// Values related to deprecrated fields.
	mongodbRepoType                  = "mongodb"
	mongodbReplicaSetServerType      = "replicaset"
	deprecatedHostAndPortMessage     = "`%s` is deprecated. In version 4.X, use `repoNodes` instead, which support single as well as multi-node repo types."
	deprecatedRepoProperitiesMessage = "`%s` is deprecated. In version 4.X, use `mongodb_settings` instead to set MongoDB properties."
)

func repositoryTypes() []string {
	return []string{
		"bigquery",
		"cassandra",
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

type RepoInfo struct {
	ID                  string                `json:"id"`
	Name                string                `json:"name"`
	Type                string                `json:"type"`
	Host                string                `json:"repoHost"`
	Port                uint32                `json:"repoPort"`
	ConnParams          *ConnParams           `json:"connParams"`
	Labels              []string              `json:"labels"`
	MaxAllowedListeners uint32                `json:"maxAllowedListeners,omitempty"`
	Properties          *RepositoryProperties `json:"properties,omitempty"`
}

type ConnParams struct {
	ConnDraining *ConnDraining `json:"connDraining"`
}

type ConnDraining struct {
	Auto     bool   `json:"auto"`
	WaitTime uint32 `json:"waitTime"`
}

type RepositoryProperties struct {
	// Replica set
	MongoDBReplicaSetName string `json:"mongodb-replicaset-name,omitempty"`
	MongoDBServerType     string `json:"mongodb-server-type,omitempty"`
}

type GetRepoByIDResponse struct {
	Repo RepoInfo `json:"repo"`
}

func (res *GetRepoByIDResponse) WriteToSchema(d *schema.ResourceData) error {
	return res.Repo.WriteToSchema(d)
}

func (res *RepoInfo) WriteToSchema(d *schema.ResourceData) error {
	d.Set(RepoTypeKey, res.Type)
	d.Set(RepoHostKey, res.Host)
	d.Set(RepoPortKey, res.Port)
	d.Set(RepoNameKey, res.Name)
	d.Set(RepoLabelsKey, res.LabelsAsInterface())
	d.Set(RepoConnDrainingKey, res.ConnDrainingAsInterface())
	if properties := res.PropertiesAsInterface(); properties != nil {
		d.Set("properties", properties)
	}
	return nil
}

func (r *RepoInfo) LabelsAsInterface() []interface{} {
	if r.Labels == nil {
		return nil
	}
	result := make([]interface{}, len(r.Labels))
	for i, v := range r.Labels {
		result[i] = v
	}
	return result
}

func (r *RepoInfo) LabelsFromInterface(i []interface{}) {
	labels := make([]string, len(i))
	for i, v := range i {
		labels[i] = v.(string)
	}
	r.Labels = labels
}

func (r *RepoInfo) ConnDrainingAsInterface() []interface{} {
	if r.ConnParams == nil || r.ConnParams.ConnDraining == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		RepoConnDrainingAutoKey:     r.ConnParams.ConnDraining.Auto,
		RepoConnDrainingWaitTimeKey: r.ConnParams.ConnDraining.WaitTime,
	}}
}

func (r *RepoInfo) ConnDrainingFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	r.ConnParams = &ConnParams{
		ConnDraining: &ConnDraining{
			Auto:     i[0].(map[string]interface{})[RepoConnDrainingAutoKey].(bool),
			WaitTime: uint32(i[0].(map[string]interface{})[RepoConnDrainingWaitTimeKey].(int)),
		},
	}
}

func (r *RepoInfo) PropertiesAsInterface() []interface{} {
	if !r.IsReplicaSet() {
		return nil
	}

	return []interface{}{map[string]interface{}{
		RepoMongoDBReplicaSetKey: []interface{}{map[string]interface{}{
			RepoMaxNodesKey:     r.MaxAllowedListeners,
			RepoReplicaSetIDKey: r.Properties.MongoDBReplicaSetName,
		},
		}}}

}

func (r *RepoInfo) PropertiesFromInterface(i []interface{}) error {
	if len(i) == 0 {
		return nil
	}
	return r.ReplicaSetFromInterface(i[0].(map[string]interface{})[RepoMongoDBReplicaSetKey].(*schema.Set).List())
}

func (r *RepoInfo) ReplicaSetFromInterface(i []interface{}) error {
	if len(i) == 0 {
		return nil
	}

	if r.Type != mongodbRepoType {
		return fmt.Errorf(
			"replica sets are only supported for repository type '%s'",
			mongodbRepoType)
	}
	r.Properties = &RepositoryProperties{
		MongoDBReplicaSetName: i[0].(map[string]interface{})[RepoReplicaSetIDKey].(string),
		MongoDBServerType:     mongodbReplicaSetServerType,
	}
	r.MaxAllowedListeners = uint32(i[0].(map[string]interface{})[RepoMaxNodesKey].(int))
	return nil
}

func (r *RepoInfo) ReadFromSchema(d *schema.ResourceData) error {
	r.ID = d.Id()
	r.Name = d.Get(RepoNameKey).(string)
	r.Type = d.Get(RepoTypeKey).(string)
	r.Host = d.Get(RepoHostKey).(string)
	r.Port = uint32(d.Get(RepoPortKey).(int))
	r.ConnDrainingFromInterface(d.Get(RepoConnDrainingKey).(*schema.Set).List())
	r.LabelsFromInterface(d.Get(RepoLabelsKey).([]interface{}))
	return r.PropertiesFromInterface(d.Get(RepoPropertiesKey).(*schema.Set).List())
}

func (data *RepoInfo) IsReplicaSet() bool {
	return data.Properties != nil && data.Properties.MongoDBServerType == mongodbReplicaSetServerType
}

var ReadRepositoryConfig = ResourceOperationConfig{
	Name:       "RepositoryRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos/%s",
			c.ControlPlane,
			d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) ResponseData {
		return &GetRepoByIDResponse{}
	},
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [repositories](https://cyral.com/docs/manage-repositories/repo-track)." +
			"\n\nSee also [Cyral Repository Configuration Module](https://github.com/cyralinc/terraform-cyral-repository-config)." +
			"\nThis module provides the repository configuration options as shown in Cyral UI.",
		CreateContext: CreateResource(
			ResourceOperationConfig{
				Name:       "RepositoryCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos",
						c.ControlPlane,
					)
				},
				NewResourceData: func() ResourceData {
					return &RepoInfo{}
				},
				NewResponseData: func(_ *schema.ResourceData) ResponseData {
					return &IDBasedResponse{}
				},
			},
			ReadRepositoryConfig,
		),
		ReadContext: ReadResource(ReadRepositoryConfig),
		UpdateContext: UpdateResource(
			ResourceOperationConfig{
				Name:       "RepositoryUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
				NewResourceData: func() ResourceData {
					return &RepoInfo{}
				},
			},
			ReadRepositoryConfig,
		),
		DeleteContext: DeleteResource(
			ResourceOperationConfig{
				Name:       "RepositoryDelete",
				HttpMethod: http.MethodDelete,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
			},
		),

		Schema: map[string]*schema.Schema{
			RepoIDKey: {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			RepoTypeKey: {
				Description:  "Repository type. List of supported types:" + supportedTypesMarkdown(repositoryTypes()),
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(repositoryTypes(), false),
			},
			RepoHostKey: {
				Description: "Repository host name (ex: `somerepo.cyral.com`).",
				Type:        schema.TypeString,
				Required:    true,
				Deprecated:  fmt.Sprintf(deprecatedHostAndPortMessage, "host"),
			},
			RepoPortKey: {
				Description: "Repository access port (ex: `3306`).",
				Type:        schema.TypeInt,
				Required:    true,
				Deprecated:  fmt.Sprintf(deprecatedHostAndPortMessage, "port"),
			},
			RepoNameKey: {
				Description: "Repository name that will be used internally in the control plane (ex: `your_repo_name`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			RepoLabelsKey: {
				Description: "Labels enable you to categorize your repository.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			RepoConnDrainingKey: {
				Description: "Parameters related to connection draining.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoConnDrainingAutoKey: {
							Description: "Whether connections should be drained automatically after a listener dies.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						RepoConnDrainingWaitTimeKey: {
							Description: "Seconds to wait to let connections drain before starting to kill all the connections, " +
								"if auto is set to true.",
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			RepoPropertiesKey: {
				Description: "Contains advanced repository configuration.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Deprecated:  fmt.Sprintf(deprecatedRepoProperitiesMessage, "properties"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoMongoDBReplicaSetKey: {
							Description: "Used to configure a MongoDB cluster.",
							Type:        schema.TypeSet,
							Optional:    true,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									RepoMaxNodesKey: {
										Description:  "Maximum number of nodes of the replica set cluster.",
										Type:         schema.TypeInt,
										Required:     true,
										Deprecated:   fmt.Sprintf(deprecatedRepoProperitiesMessage, "max_nodes"),
										ValidateFunc: validation.IntAtLeast(1),
									},
									RepoReplicaSetIDKey: {
										Description:  "Identifier of the replica set cluster. Used to construct the URI command (available in Cyral's Access Portal page) that your users will need for connecting to the repository via Cyral.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
										Deprecated:   fmt.Sprintf(deprecatedRepoProperitiesMessage, "replica_set_id"),
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
