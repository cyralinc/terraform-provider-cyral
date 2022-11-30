package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

const (
	RepoListKey = "repository_list"
)

// GetReposSubResponse is different from GetRepoByIDResponse. For the by-id
// reponse, we expect the ids to be embedded in the RepoInfo struct. For
// GetReposSubResponse, the ids come outside of RepoInfo.
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
		repoID := repo.ID
		repoData := repo.Repo
		argumentVals := map[string]interface{}{
			RepoIDKey:     repoID,
			RepoNameKey:   repoData.Name,
			RepoTypeKey:   repoData.Type,
			RepoLabelsKey: repoData.LabelsAsInterface(),
			RepoNodesKey:  repoData.RepoNodesAsInterface(),
		}
		repoList = append(repoList, argumentVals)
	}

	if err := d.Set("repository_list", repoList); err != nil {
		return err
	}

	d.SetId(uuid.New().String())

	return nil
}

func dataSourceRepositoryReadConfig() ResourceOperationConfig {
	return ResourceOperationConfig{
		Name:       "RepositoryDataSourceRead",
		HttpMethod: http.MethodGet,
		CreateURL: func(d *schema.ResourceData, c *client.Client) string {
			nameFilter := d.Get("name").(string)
			typeFilter := d.Get("type").(string)
			urlParams := urlQuery(map[string]string{
				"name": nameFilter,
				"type": typeFilter,
			})

			return fmt.Sprintf("https://%s/v1/repos%s", c.ControlPlane, urlParams)
		},
		NewResponseData: func(_ *schema.ResourceData) ResponseData { return &GetReposResponse{} },
	}
}

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter repositories.",
		ReadContext: ReadResource(dataSourceRepositoryReadConfig()),
		Schema: map[string]*schema.Schema{
			RepoNameKey: {
				Description: "Filter the results by a regular expression (regex) that matches names of existing repositories.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			RepoTypeKey: {
				Description:  "Filter the results by type of repository. List of supported types:" + supportedTypesMarkdown(repositoryTypes()),
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(append(repositoryTypes(), ""), false),
			},
			RepoListKey: {
				Description: "List of existing repositories satisfying given filter criteria.",
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
						RepoNodesKey: {
							Description: "List of nodes for this repository.",
							Type:        schema.TypeList,
							Optional:    true,
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
					},
				},
			},
		},
	}
}
