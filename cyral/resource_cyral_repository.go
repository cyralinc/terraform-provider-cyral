package cyral

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cyralinc/terraform-provider-cyral/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type CreateRepoResponse struct {
	ID string `json:"ID"`
}

type GetRepoByNameResponse struct {
	Repos []RepoCompoundData `json:"repos"`
}

type RepoCompoundData struct {
	ID   string   `json:"id"`
	Repo RepoData `json:"repo"`
}

type RepoData struct {
	RepoType string `json:"type"`
	Name     string `json:"name"`
	Host     string `json:"repoHost"`
	Port     int    `json:"repoPort"`
}

func resourceCyralRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceCyralRepositoryCreate,
		Read:   resourceCyralRepositoryRead,
		Update: resourceCyralRepositoryUpdate,
		Delete: resourceCyralRepositoryDelete,

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceCyralRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)
	repoData, err := getRepoDataFromResource(c, d)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/v1/repos", c.ControlPlane)
	payloadBytes, err := json.Marshal(repoData)
	if err != nil {
		return fmt.Errorf("failed to encode 'create repo' payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create 'create repo' request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute 'create repo' request. Check the control plane address; err: %v", err)
	}

	if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("repository name already exists in control plane")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("unable to read data from request body; err: %v", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected response from 'create repo' request; status code: %d; body: %q",
			res.StatusCode, body)
	}

	bodyRep := CreateRepoResponse{}
	err = json.Unmarshal(body, &bodyRep)
	if err != nil {
		return fmt.Errorf("unable to unmarshall json; err: %v", err)
	}

	c.Repository.Name = repoData.Name
	d.SetId(repoData.Name)

	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	repoData, err := resourceCyralRepositoryFindByName(c, d.Id())

	if err != nil {
		return err
	}

	if repoData == nil {
		return fmt.Errorf("repo not found; name: %s", d.Id())
	}

	d.Set("type", repoData.Repo.RepoType)
	d.Set("host", repoData.Repo.Host)
	d.Set("port", repoData.Repo.Port)
	d.Set("name", repoData.Repo.Name)

	return nil
}

func resourceCyralRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	repoDataFromResource, err := getRepoDataFromResource(c, d)
	if err != nil {
		return err
	}

	repoDataFromCP, err := resourceCyralRepositoryFindByName(c, d.Id())
	if err != nil {
		return fmt.Errorf("unable to find repo by name during update operation: %v", err)
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, repoDataFromCP.ID)
	payloadBytes, err := json.Marshal(repoDataFromResource)
	if err != nil {
		return fmt.Errorf("failed to encode 'update repo' payload: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return fmt.Errorf("unable to create 'update repo' request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to execute 'update repo' request. Check the control plane address; err: %v", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository not found; id: %s", d.Id())
	} else if res.StatusCode == http.StatusConflict {
		return fmt.Errorf("repository name already exists in control plane")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from 'update repo'; status code: %d; body: %q",
			res.StatusCode, res.Body)
	}

	// Name may have been changed by user, so we must update the tf state
	d.SetId(repoDataFromResource.Name)

	return resourceCyralRepositoryRead(d, m)
}

func resourceCyralRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*client.Client)

	repoData, err := resourceCyralRepositoryFindByName(c, d.Id())
	if err != nil {
		return fmt.Errorf("unable to find repo by name during delete operation: %v", err)
	}

	url := fmt.Sprintf("https://%s/v1/repos/%s", c.ControlPlane, repoData.ID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create 'delete repo' request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable execute 'delete repo' request. Check the control plane address; err: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response from 'delete repo' request; status code: %d; body: %q", res.StatusCode, res.Body)
	}

	return nil
}

func resourceCyralRepositoryFindByName(c *client.Client, name string) (*RepoCompoundData, error) {
	url := fmt.Sprintf("https://%s/v1/repos?name=%s", c.ControlPlane, name)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new request; err: %v", err)
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", c.TokenType, c.Token))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to execute request at resourceCyralRepositoryFindByName."+
			" Check the control plane address; err: %v", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"unable to read data from request body at resourceCyralRepositoryFindByName; err: %v",
			err)
	}

	// Not an error, nor any data was found
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected response from resourceCyralRepositoryFindByName; status code: %d; body: %q",
			res.StatusCode, res.Body)
	}
	getRepoByNameResp := GetRepoByNameResponse{
		Repos: []RepoCompoundData{},
	}
	err = json.Unmarshal(body, &getRepoByNameResp)
	if err != nil {
		return nil, fmt.Errorf("unable to get repo json by name, err: %v", err)
	}

	r := findRepoByName(getRepoByNameResp, name)
	if r == nil {
		return nil, fmt.Errorf("unable to find repo named '%s' in response %v",
			name, getRepoByNameResp)
	}

	return r, nil
}

func findRepoByName(resp GetRepoByNameResponse, name string) *RepoCompoundData {
	for _, r := range resp.Repos {
		if r.Repo.Name == name {
			return &r
		}
	}
	return nil
}

func getRepoDataFromResource(c *client.Client, d *schema.ResourceData) (RepoData, error) {
	repoType := d.Get("type").(string)
	err := containsRepoType(repoType)
	if err != nil {
		return RepoData{}, err
	}

	return RepoData{
		RepoType: repoType,
		Host:     d.Get("host").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
	}, nil
}

func containsRepoType(repoType string) error {
	// This code was copied here to remove dependency of CRUD,
	// but we should move the CRUD code to CRUD-API (or somewhere
	// else) in the future.
	repoTypes := map[string]bool{
		"bigquery":   true,
		"cassandra":  true,
		"dremio":     true,
		"galera":     true,
		"mariadb":    true,
		"mongodb":    true,
		"mysql":      true,
		"postgresql": true,
		"snowflake":  true,
		"sqlserver":  true,
	}
	if repoTypes[repoType] == false {
		return fmt.Errorf("repo type must be one of %v", repoTypes)
	}
	return nil
}
