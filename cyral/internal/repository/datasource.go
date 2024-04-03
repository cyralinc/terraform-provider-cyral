package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/cyral/client"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core"
	"github.com/cyralinc/terraform-provider-cyral/cyral/core/types/resourcetype"
	"github.com/cyralinc/terraform-provider-cyral/cyral/utils"
)

const (
	RepoListKey = "repository_list"
)

// GetReposSubResponse is different from GetRepoByIDResponse. For the by-id
// response, we expect the ids to be embedded in the RepoInfo struct. For
// GetReposSubResponse, the ids come outside of RepoInfo.
//
// Needles to say we need a new API version to fix these issues. For the
// time being, I'm keeping the model here to isolate it from the rest of
// the code - Wilson.
type GetReposSubResponse struct {
	ID   string   `json:"id"`
	Repo RepoInfo `json:"repo"`
}

type GetReposResponse struct {
	Repos []GetReposSubResponse `json:"repos"`
}

func (resp *GetReposResponse) WriteToSchema(d *schema.ResourceData) error {
	var repoList []interface{}
	for _, repo := range resp.Repos {
		argumentVals := map[string]interface{}{
			RepoIDKey:              repo.ID,
			RepoNameKey:            repo.Repo.Name,
			RepoTypeKey:            repo.Repo.Type,
			RepoLabelsKey:          repo.Repo.Labels.AsInterface(),
			RepoConnDrainingKey:    repo.Repo.ConnParams.AsInterface(),
			RepoNodesKey:           repo.Repo.RepoNodes.AsInterface(),
			RepoMongoDBSettingsKey: repo.Repo.MongoDBSettings.AsInterface(),
		}
		repoList = append(repoList, argumentVals)
	}

	if err := d.Set("repository_list", repoList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

var dsContextHandler = core.DefaultContextHandler{
	ResourceName:                 dataSourceName,
	ResourceType:                 resourcetype.DataSource,
	SchemaWriterFactoryGetMethod: func(_ *schema.ResourceData) core.SchemaWriter { return &GetReposResponse{} },
	IdBasedURLFactory: func(d *schema.ResourceData, c *client.Client) string {
		nameFilter := d.Get("name").(string)
		typeFilter := d.Get("type").(string)
		urlParams := utils.UrlQuery(map[string]string{
			"name": nameFilter,
			"type": typeFilter,
		})

		return fmt.Sprintf("https://%s/v1/repos%s", c.ControlPlane, urlParams)
	},
}

func dataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves a list of repositories. See [`repository_list`](#nestedatt--repository_list).",
		ReadContext: dsContextHandler.ReadContext(),
		Schema: map[string]*schema.Schema{
			RepoNameKey: {
				Description: "Filter the results by a regular expression (regex) that matches names of existing repositories.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			RepoTypeKey: {
				Description:  "Filter the results by type of repository. List of supported types:" + utils.SupportedValuesAsMarkdown(RepositoryTypes()),
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(append(RepositoryTypes(), ""), false),
			},
			RepoListKey: {
				Description: "List of existing repositories satisfying the filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						RepoIDKey: {
							Description: "ID of the repository in the Cyral environment.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						RepoNameKey: {
							Description: "Repository name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						RepoTypeKey: {
							Description: "Repository type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						RepoLabelsKey: {
							Description: "Repository labels.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						RepoConnDrainingKey: {
							Description: "Parameters related to connection draining.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									RepoConnDrainingAutoKey: {
										Description: "Whether connections should be drained automatically after a listener dies.",
										Type:        schema.TypeBool,
										Computed:    true,
									},
									RepoConnDrainingWaitTimeKey: {
										Description: "Seconds to wait to let connections drain before starting to kill all the connections, " +
											"if auto is set to true.",
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						RepoNodesKey: {
							Description: "List of nodes for this repository.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									RepoNameKey: {
										Description: "Name of the repo node.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									RepoHostKey: {
										Description: "Repo node host (ex: `somerepo.cyral.com`). Can be empty if node is dynamic.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									RepoPortKey: {
										Description: "Repository access port (ex: `3306`). Can be empty if node is dynamic.",
										Type:        schema.TypeInt,
										Computed:    true,
									},
									RepoNodeDynamicKey: {
										Description: "Indicates if node is dynamically discovered. If true, `host` and `port` must be empty.",
										Type:        schema.TypeBool,
										Computed:    true,
									},
								},
							},
						},
						RepoMongoDBSettingsKey: {
							Description: "Parameters related to MongoDB repositories.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									RepoMongoDBReplicaSetNameKey: {
										Description: "Name of the replica set, if applicable.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									RepoMongoDBServerTypeKey: {
										Description: "Type of the MongoDB server.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									RepoMongoDBSRVRecordName: {
										Description: "Name of a DNS SRV record which contains cluster topology details.",
										Type:        schema.TypeString,
										Optional:    true,
									},
									RepoMongoDBFlavorKey: {
										Description: "The flavor of the MongoDB deployment.",
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
