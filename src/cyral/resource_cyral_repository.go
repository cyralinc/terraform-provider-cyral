package cyral

import (
	"fmt"
	"net/http"

	"github.com/cyralinc/terraform-provider-cyral/src/client"
	"github.com/cyralinc/terraform-provider-cyral/src/core"
	"github.com/cyralinc/terraform-provider-cyral/src/utils"
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
	// MongoDB settings keys.
	RepoMongoDBSettingsKey       = "mongodb_settings"
	RepoMongoDBReplicaSetNameKey = "replica_set_name"
	RepoMongoDBServerTypeKey     = "server_type"
	RepoMongoDBSRVRecordName     = "srv_record_name"
)

const (
	Denodo          = "denodo"
	Dremio          = "dremio"
	DynamoDB        = "dynamodb"
	DynamoDBStreams = "dynamodbstreams"
	Galera          = "galera"
	MariaDB         = "mariadb"
	MongoDB         = "mongodb"
	MySQL           = "mysql"
	Oracle          = "oracle"
	PostgreSQL      = "postgresql"
	Redshift        = "redshift"
	S3              = "s3"
	Snowflake       = "snowflake"
	SQLServer       = "sqlserver"
)

func repositoryTypes() []string {
	return []string{
		Denodo,
		Dremio,
		DynamoDB,
		DynamoDBStreams,
		Galera,
		MariaDB,
		MongoDB,
		MySQL,
		Oracle,
		PostgreSQL,
		Redshift,
		S3,
		Snowflake,
		SQLServer,
	}
}

const (
	ReplicaSet = "replicaset"
	Standalone = "standalone"
	Sharded    = "sharded"
)

func mongoServerTypes() []string {
	return []string{
		ReplicaSet,
		Standalone,
		Sharded,
	}
}

type RepoInfo struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Type            string           `json:"type"`
	Host            string           `json:"repoHost"`
	Port            uint32           `json:"repoPort"`
	ConnParams      *ConnParams      `json:"connParams"`
	Labels          []string         `json:"labels"`
	RepoNodes       []*RepoNode      `json:"repoNodes,omitempty"`
	MongoDBSettings *MongoDBSettings `json:"mongoDbSettings,omitempty"`
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
	SRVRecordName  string `json:"srvRecordName,omitempty"`
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
	return nil
}

func (r *RepoInfo) ReadFromSchema(d *schema.ResourceData) error {
	r.ID = d.Id()
	r.Name = d.Get(RepoNameKey).(string)
	r.Type = d.Get(RepoTypeKey).(string)
	r.LabelsFromInterface(d.Get(RepoLabelsKey).([]interface{}))
	r.RepoNodesFromInterface(d.Get(RepoNodesKey).([]interface{}))
	r.ConnDrainingFromInterface(d.Get(RepoConnDrainingKey).(*schema.Set).List())
	var mongoDBSettings = d.Get(RepoMongoDBSettingsKey).(*schema.Set).List()
	if r.Type == MongoDB && (mongoDBSettings == nil || len(mongoDBSettings) == 0) {
		return fmt.Errorf("'%s' block must be provided when '%s=%s'", RepoMongoDBSettingsKey, TypeKey, MongoDB)
	} else if r.Type != MongoDB && len(mongoDBSettings) > 0 {
		return fmt.Errorf("'%s' block is only allowed when '%s=%s'", RepoMongoDBSettingsKey, TypeKey, MongoDB)
	}
	return r.MongoDBSettingsFromInterface(mongoDBSettings)
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
		RepoMongoDBSRVRecordName:     r.MongoDBSettings.SRVRecordName,
	}}
}

func (r *RepoInfo) MongoDBSettingsFromInterface(i []interface{}) error {
	if len(i) == 0 {
		return nil
	}
	var replicaSetName = i[0].(map[string]interface{})[RepoMongoDBReplicaSetNameKey].(string)
	var serverType = i[0].(map[string]interface{})[RepoMongoDBServerTypeKey].(string)
	var srvRecordName = i[0].(map[string]interface{})[RepoMongoDBSRVRecordName].(string)
	if serverType == ReplicaSet && replicaSetName == "" {
		return fmt.Errorf("'%s' must be provided when '%s=\"%s\"'", RepoMongoDBReplicaSetNameKey,
			RepoMongoDBServerTypeKey, ReplicaSet)
	}
	if serverType != ReplicaSet && replicaSetName != "" {
		return fmt.Errorf("'%s' cannot be provided when '%s=\"%s\"'", RepoMongoDBReplicaSetNameKey,
			RepoMongoDBServerTypeKey, serverType)
	}
	if serverType == Standalone && srvRecordName != "" {
		return fmt.Errorf(
			"'%s' cannot be provided when '%s=\"%s\"'",
			RepoMongoDBSRVRecordName,
			RepoMongoDBServerTypeKey,
			Standalone,
		)
	}
	r.MongoDBSettings = &MongoDBSettings{
		ReplicaSetName: i[0].(map[string]interface{})[RepoMongoDBReplicaSetNameKey].(string),
		ServerType:     i[0].(map[string]interface{})[RepoMongoDBServerTypeKey].(string),
		SRVRecordName:  i[0].(map[string]interface{})[RepoMongoDBSRVRecordName].(string),
	}
	return nil
}

var ReadRepositoryConfig = core.ResourceOperationConfig{
	Name:       "RepositoryRead",
	HttpMethod: http.MethodGet,
	CreateURL: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos/%s",
			c.ControlPlane,
			d.Id(),
		)
	},
	NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
		return &GetRepoByIDResponse{}
	},
	RequestErrorHandler: &core.ReadIgnoreHttpNotFound{ResName: "Repository"},
}

func resourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Manages [repositories](https://cyral.com/docs/manage-repositories/repo-track)." +
			"\n\nSee also [Cyral Repository Configuration Module](https://github.com/cyralinc/terraform-cyral-repository-config)." +
			"\nThis module provides the repository configuration options as shown in Cyral UI.",
		CreateContext: core.CreateResource(
			core.ResourceOperationConfig{
				Name:       "RepositoryCreate",
				HttpMethod: http.MethodPost,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos",
						c.ControlPlane,
					)
				},
				NewResourceData: func() core.ResourceData {
					return &RepoInfo{}
				},
				NewResponseData: func(_ *schema.ResourceData) core.ResponseData {
					return &core.IDBasedResponse{}
				},
			},
			ReadRepositoryConfig,
		),
		ReadContext: core.ReadResource(ReadRepositoryConfig),
		UpdateContext: core.UpdateResource(
			core.ResourceOperationConfig{
				Name:       "RepositoryUpdate",
				HttpMethod: http.MethodPut,
				CreateURL: func(d *schema.ResourceData, c *client.Client) string {
					return fmt.Sprintf(
						"https://%s/v1/repos/%s",
						c.ControlPlane,
						d.Id(),
					)
				},
				NewResourceData: func() core.ResourceData {
					return &RepoInfo{}
				},
			},
			ReadRepositoryConfig,
		),
		DeleteContext: core.DeleteResource(
			core.ResourceOperationConfig{
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
				Description:  "Repository type. List of supported types:" + utils.SupportedTypesMarkdown(repositoryTypes()),
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
							Description: "*Only supported for MongoDB in cluster configurations.*\n" +
								"Indicates if the node is dynamically discovered, meaning that the sidecar " +
								"will query the cluster to get the topology information and discover the " +
								"addresses of the dynamic nodes. If set to `true`, `host` and `port` must " +
								"be empty. A node with value of this field as false considered `static`.\n" +
								"The following conditions apply: \n" +
								"  - The total number of declared `" + RepoNodesKey + "` blocks must match " +
								"the actual number of nodes in the cluster.\n" +
								"  - If there are static nodes in the configuration, they must be declared " +
								"before all dynamic nodes.\n" +
								"  - See the MongoDB-specific configuration in the [" + RepoMongoDBSettingsKey +
								"](#nested-schema-for-" + RepoMongoDBSettingsKey + ").",
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			RepoMongoDBSettingsKey: {
				Description: "Parameters related to MongoDB repositories.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoMongoDBReplicaSetNameKey: {
							Description: "Name of the replica set, if applicable.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoMongoDBServerTypeKey: {
							Description: "Type of the MongoDB server. Allowed values: " + utils.SupportedTypesMarkdown(mongoServerTypes()) +
								"\n\n  The following conditions apply:\n" +
								"  - If `" + Sharded + "` and `" + RepoMongoDBSRVRecordName + "` *not* provided, then all `" +
								RepoNodesKey + "` blocks must be static (see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")).\n" +
								"  - If `" + Sharded + "` and `" + RepoMongoDBSRVRecordName + "` provided, then all `" +
								RepoNodesKey + "` blocks must be dynamic (see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")).\n" +
								"  - If `" + Standalone + "`, then only one `" + RepoNodesKey +
								"` block can be declared and it must be static (see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")). The `" +
								RepoMongoDBSRVRecordName + "` is not supported in this configuration.\n" +
								"  - If `" + ReplicaSet + "` and `" + RepoMongoDBSRVRecordName + "` *not* provided, then `" +
								RepoNodesKey + "` blocks may mix dynamic and static nodes (see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")).\n" +
								"  - If `" + ReplicaSet + "` and `" + RepoMongoDBSRVRecordName + "` provided, then `" +
								RepoNodesKey + "` blocks must be dynamic (see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")).\n",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(mongoServerTypes(), false),
						},
						RepoMongoDBSRVRecordName: {
							Description: "Name of a DNS SRV record which contains cluster topology details. " +
								"If specified, then all `" + RepoNodesKey + "` blocks must be declared dynamic " +
								"(see [`" + RepoNodeDynamicKey + "`](#" + RepoNodeDynamicKey + ")). " +
								"Only supported for `" + RepoMongoDBServerTypeKey + "=\"" + Sharded + "\"` or `" +
								RepoMongoDBServerTypeKey + "=\"" + ReplicaSet + "\".",
							Type:     schema.TypeString,
							Optional: true,
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
