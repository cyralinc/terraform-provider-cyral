package cyral

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

// GetReposSubResponse is different from GetRepoByIDResponse. For the by-id
// reponse, we expect the ids to be embedded in the RepoData struct. For
// GetReposSubResponse, the ids come outside of RepoData.
type GetReposSubResponse struct {
	ID   string   `json:"id"`
	Repo RepoData `json:"repo"`
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
			"id":         repoID,
			"name":       repoData.Name,
			"type":       repoData.RepoType,
			"host":       repoData.Host,
			"port":       repoData.Port,
			"labels":     repoData.Labels,
			"properties": repoData.PropertiesAsInterface(),
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
		NewResponseData: func() ResponseData { return &GetReposResponse{} },
	}
}

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter repositories.",
		ReadContext: ReadResource(dataSourceRepositoryReadConfig()),
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Filter the results by a regular expression (regex) that matches names of existing repositories.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  "Filter the results by type of repository. List of supported types:" + supportedTypesMarkdown(repositoryTypes()),
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(append(repositoryTypes(), ""), false),
			},
			"repository_list": {
				Description: "List of existing repositories satisfying given filter criteria.",
				Computed:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the repository in the Cyral environment.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Repository name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "Repository type.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"host": {
							Description: "Repository host name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"port": {
							Description: "Repository access port.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"labels": {
							Description: "Repository labels.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"properties": {
							Description: "Advanced repository configuration.",
							Type:        schema.TypeSet,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mongodb_replica_set": {
										Description: "MongoDB replica set configuration.",
										Type:        schema.TypeSet,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_nodes": {
													Description: "Maximum number of nodes of the replica set cluster.",
													Type:        schema.TypeInt,
													Computed:    true,
												},
												"replica_set_id": {
													Description: "Identifier of the replica set cluster.",
													Type:        schema.TypeString,
													Computed:    true,
												},
											},
										},
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
