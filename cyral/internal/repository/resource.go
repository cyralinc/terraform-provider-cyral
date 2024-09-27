package repository

import (
	"fmt"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var resourceContextHandler = core.DefaultContextHandler{
	ResourceName:                 resourceName,
	ResourceType:                 resourcetype.Resource,
	SchemaReaderFactory:          func() core.SchemaReader { return &RepoInfo{} },
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &GetRepoByIDResponse{} },
	BaseURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		return fmt.Sprintf(
			"https://%s/v1/repos",
			c.ControlPlane,
		)
	},
}

func resourceSchema() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages [repositories](https://cyral.com/docs/how-to/track-repos/).",
		CreateContext: resourceContextHandler.CreateContext(),
		ReadContext:   resourceContextHandler.ReadContext(),
		UpdateContext: resourceContextHandler.UpdateContext(),
		DeleteContext: resourceContextHandler.DeleteContext(),
		Schema: map[string]*schema.Schema{
			RepoIDKey: {
				Description: "ID of this resource in Cyral environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			RepoTypeKey: {
				Description:  "Repository type. List of supported types:" + utils.SupportedValuesAsMarkdown(RepositoryTypes()),
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(RepositoryTypes(), false),
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
							Description: "Type of the MongoDB server. Allowed values: " + utils.SupportedValuesAsMarkdown(mongoServerTypes()) +
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
						RepoMongoDBFlavorKey: {
							Description: "The flavor of the MongoDB deployment. Allowed values: " + utils.SupportedValuesAsMarkdown(mongoFlavors()) +
								"\n\n  The following conditions apply:\n" +
								"  - The `" + MongoDBFlavorDocumentDB + "` flavor cannot be combined with the MongoDB Server type `" + Sharded + "`.\n",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(mongoFlavors(), false),
						},
					},
				},
			},
			RepoRedshiftSettingsKey: {
				Description: "Parameters related to Redshift repositories.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoRedshiftClusterIdentifier: {
							Description: "Name of the provisioned cluster.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoRedshiftWorkgroupName: {
							Description: "Workgroup name for serverless cluster.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						RepoRedshiftAWSRegion: {
							Description: "Code of the AWS region where the Redshift instance is deployed.",
							Type:        schema.TypeString,
							Optional:    true,
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
