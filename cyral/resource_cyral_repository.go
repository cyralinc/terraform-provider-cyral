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
	RepoIDKey     = "id"
	RepoTypeKey   = "type"
	RepoNameKey   = "name"
	RepoLabelsKey = "labels"
	// Connection draining keys.
	RepoConnDrainingKey         = "connection_draining"
	RepoConnDrainingAutoKey     = "auto"
	RepoConnDrainingWaitTimeKey = "wait_time"
	// Repo node keys.
	RepoNodesKey       = "repo_node"
	RepoHostKey        = "host"
	RepoPortKey        = "port"
	RepoNodeDynamicKey = "dynamic"
	// Access gateway keys.
	RepoPreferredAccessGatewayKey = "preferred_access_gateway"
	RepoSidecarIDKey              = "sidecar_id"
	RepoBindingIDKey              = "binding_id"
	// MongoDB settings keys.
	RepoMongoDBSettingsKey       = "mongodb_settings"
	RepoMongoDBReplicaSetNameKey = "replica_set_name"
	RepoMongoDBServerTypeKey     = "server_type"
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

func mongoServerTypes() []string {
	return []string{
		"replicaset",
		"standalone",
	}
}

type RepoInfo struct {
	ID                       string           `json:"id"`
	Name                     string           `json:"name"`
	Type                     string           `json:"type"`
	Host                     string           `json:"repoHost"`
	Port                     uint32           `json:"repoPort"`
	ConnParams               *ConnParams      `json:"connParams"`
	Labels                   []string         `json:"labels"`
	RepoNodes                []*RepoNode      `json:"repoNodes,omitempty"`
	MongoDBSettings          *MongoDBSettings `json:"mongoDbSettings,omitempty"`
	PreferredAccessGwBinding *BindingKey      `json:"preferredAccessGwBinding,omitempty"`
}

type ConnParams struct {
	ConnDraining *ConnDraining `json:"connDraining"`
}

type ConnDraining struct {
	Auto     bool   `json:"auto"`
	WaitTime uint32 `json:"waitTime"`
}

type MongoDBSettings struct {
	ReplicaSetName string `json:"replicaSetName,omitempty"`
	ServerType     string `json:"serverType,omitempty"`
}

type BindingKey struct {
	SidecarID string `json:"sidecarId,omitempty"`
	BindingID string `json:"bindingId,omitempty"`
}

type RepoNode struct {
	Name    string `json:"name"`
	Host    string `json:"host"`
	Port    uint32 `json:"port"`
	Dynamic bool   `json:"dynamic"`
}

type GetRepoByIDResponse struct {
	Repo RepoInfo `json:"repo"`
}

func (res *GetRepoByIDResponse) WriteToSchema(d *schema.ResourceData) error {
	return res.Repo.WriteToSchema(d)
}

func (res *RepoInfo) WriteToSchema(d *schema.ResourceData) error {
	d.Set(RepoTypeKey, res.Type)
	d.Set(RepoNameKey, res.Name)
	d.Set(RepoLabelsKey, res.LabelsAsInterface())
	d.Set(RepoConnDrainingKey, res.ConnDrainingAsInterface())
	d.Set(RepoNodesKey, res.RepoNodesAsInterface())
	d.Set(RepoMongoDBSettingsKey, res.MongoDBSettingsAsInterface())
	d.Set(RepoPreferredAccessGatewayKey, res.AccessGatewayAsInterface())
	return nil
}

func (r *RepoInfo) ReadFromSchema(d *schema.ResourceData) error {
	r.ID = d.Id()
	r.Name = d.Get(RepoNameKey).(string)
	r.Type = d.Get(RepoTypeKey).(string)
	r.LabelsFromInterface(d.Get(RepoLabelsKey).([]interface{}))
	r.RepoNodesFromInterface(d.Get(RepoNodesKey).([]interface{}))
	r.ConnDrainingFromInterface(d.Get(RepoConnDrainingKey).(*schema.Set).List())
	r.AccessGatewayFromInterface(d.Get(RepoPreferredAccessGatewayKey).(*schema.Set).List())
	r.MongoDBSettingsFromInterface(d.Get(RepoMongoDBSettingsKey).(*schema.Set).List())
	return nil
}

func (r *RepoInfo) LabelsAsInterface() []interface{} {
	if r.Labels == nil {
		return nil
	}
	labels := make([]interface{}, len(r.Labels))
	for i, label := range r.Labels {
		labels[i] = label
	}
	return labels
}

func (r *RepoInfo) LabelsFromInterface(i []interface{}) {
	labels := make([]string, len(i))
	for index, v := range i {
		labels[index] = v.(string)
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

func (r *RepoInfo) AccessGatewayAsInterface() []interface{} {
	if r.PreferredAccessGwBinding == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		RepoBindingIDKey: r.PreferredAccessGwBinding.BindingID,
		RepoSidecarIDKey: r.PreferredAccessGwBinding.SidecarID,
	}}
}

func (r *RepoInfo) AccessGatewayFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	r.PreferredAccessGwBinding = &BindingKey{
		BindingID: i[0].(map[string]interface{})[RepoBindingIDKey].(string),
		SidecarID: i[0].(map[string]interface{})[RepoSidecarIDKey].(string),
	}
}

func (r *RepoInfo) RepoNodesAsInterface() []interface{} {
	if r.RepoNodes == nil {
		return nil
	}
	repoNodes := make([]interface{}, len(r.RepoNodes))
	for i, node := range r.RepoNodes {
		repoNodes[i] = map[string]interface{}{
			RepoNameKey:        node.Name,
			RepoHostKey:        node.Host,
			RepoPortKey:        node.Port,
			RepoNodeDynamicKey: node.Dynamic,
		}
	}
	return repoNodes
}

func (r *RepoInfo) RepoNodesFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	repoNodes := make([]*RepoNode, len(i))
	for index, nodeInterface := range i {
		nodeMap := nodeInterface.(map[string]interface{})
		node := &RepoNode{
			Name:    nodeMap[RepoNameKey].(string),
			Host:    nodeMap[RepoHostKey].(string),
			Port:    uint32(nodeMap[RepoPortKey].(int)),
			Dynamic: nodeMap[RepoNodeDynamicKey].(bool),
		}
		repoNodes[index] = node
	}
	r.RepoNodes = repoNodes
}

func (r *RepoInfo) MongoDBSettingsAsInterface() []interface{} {
	if r.MongoDBSettings == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		RepoMongoDBReplicaSetNameKey: r.MongoDBSettings.ReplicaSetName,
		RepoMongoDBServerTypeKey:     r.MongoDBSettings.ServerType,
	}}
}

func (r *RepoInfo) MongoDBSettingsFromInterface(i []interface{}) {
	if len(i) == 0 {
		return
	}
	r.MongoDBSettings = &MongoDBSettings{
		ReplicaSetName: i[0].(map[string]interface{})[RepoMongoDBReplicaSetNameKey].(string),
		ServerType:     i[0].(map[string]interface{})[RepoMongoDBServerTypeKey].(string),
	}
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
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(repositoryTypes(), false),
			},
			RepoNameKey: {
				Description: "Repository name that will be used internally in the control plane (ex: `your_repo_name`).",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				MaxItems:    1,
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
			RepoPreferredAccessGatewayKey: {
				Description: "Preferred access gateway for this repository.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoSidecarIDKey: {
							Description: "Sidecar ID of the preferred access gateway.",
							Type:        schema.TypeString,
							Required:    true,
						},
						RepoBindingIDKey: {
							Description: "Binding ID of the preferred access gateway.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			RepoNodesKey: {
				Description: "List of nodes for this repository.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoNameKey: {
							Description: "Name of the repo node.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoHostKey: {
							Description: "Repo node host (ex: `somerepo.cyral.com`). Can be empty if node is dynamic.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoPortKey: {
							Description: "Repository access port (ex: `3306`). Can be empty if node is dynamic.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						RepoNodeDynamicKey: {
							Description: "Indicates if node is dynamically discovered. If true, `host` and `port` must be empty.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			RepoMongoDBSettingsKey: {
				Description: "Parameters related to MongoDB repositories.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoMongoDBReplicaSetNameKey: {
							Description: "Name of the replica set, if applicable.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoMongoDBServerTypeKey: {
							Description:  "Type of the MongoDB server. Allowed values: " + supportedTypesMarkdown(mongoServerTypes()),
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(mongoServerTypes(), false),
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
