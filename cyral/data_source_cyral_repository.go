package cyral

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/cyralinc/terraform-provider-cyral/client"
)

type GetReposRequest struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type GetReposSubResponse struct {
	ID   string   `json:"id"`
	Repo RepoData `json:"repo"`
}

type GetReposResponse struct {
	Repos []GetReposSubResponse `json:"repos"`
}

func dataSourceRepository() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve and filter repositories.",
		ReadContext: dataSourceRepositoryRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Filter the results by a _regular expression_ (regex) that matches names of existing repositories.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  "Filter the results by type of repository. List of supported types:" + repositoryTypesMarkdown,
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(repositoryTypes(), false),
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
							Description: "Contains advanced repository configuration.",
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

func dataSourceRepositoryRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	log.Printf("[DEBUG] Init dataSourceRepositoryRead")
	c := m.(*client.Client)

	nameFilter := d.Get("name").(string)
	typeFilter := d.Get("type").(string)
	getReposRequest := &GetReposRequest{
		Name: nameFilter,
		Type: typeFilter,
	}
	url := fmt.Sprintf("https://%s/v1/repos", c.ControlPlane)
	body, err := c.DoRequest(url, http.MethodGet, getReposRequest)
	if err != nil {
		return createError("Unable to execute request to read repositories",
			err.Error())
	}

	resp := GetReposResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return createError("Unable to unmarshal repository list response",
			err.Error())
	}
	log.Printf("[DEBUG] Response body (unmarshalled): %#v", resp)

	var repoList []interface{}
	for _, repo := range resp.Repos {
		repoID := repo.ID
		repoData := repo.Repo
		repoList = append(repoList, map[string]interface{}{
			"id":         repoID,
			"name":       repoData.Name,
			"type":       repoData.RepoType,
			"host":       repoData.Host,
			"port":       repoData.Port,
			"labels":     repoData.Labels,
			"properties": repoData.PropertiesAsInterface(),
		})
	}

	d.Set("repository_list", repoList)

	d.SetId(uuid.New().String())

	log.Printf("[DEBUG] End dataSourceRepositoryRead")

	return nil
}
